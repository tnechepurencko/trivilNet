#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

#include "rt_defs.h"
#include "rt_api.h"

//==== crash

void panic() {
    exit(1);
}

void runtime_crash(char* s) {
	printf("!runtime_crash: %s\n", s);
    panic();
}

//==== memory

void* mm_allocate(size_t size) {
	void* a = malloc(size);
	if (a == NULL) {
		runtime_crash("memory not allocated");
	}
	return a;
}	

//==== utf-8

// Сохраняет code point в буфер, должен быть не менее 4 байтов
size_t encode_symbol(TSymbol cp, TByte *buf) {
  if (cp < 0x00) {
    return 0;
  } else if (cp < 0x80) {
    buf[0] = (TByte) cp;
    return 1;
  } else if (cp < 0x800) {
    buf[0] = (TByte)(0xC0 + (cp >> 6));
    buf[1] = (TByte)(0x80 + (cp & 0x3F));
    return 2;
  // Не учитываю диапазон 0xd800-0xdfff, хотя полученный UTF-8 будет не корректным
  } else if (cp < 0x10000) {
    buf[0] = (TByte)(0xE0 + (cp >> 12));
    buf[1] = (TByte)(0x80 + ((cp >> 6) & 0x3F));
    buf[2] = (TByte)(0x80 + (cp & 0x3F));
    return 3;
  } else if (cp < 0x110000) {
    buf[0] = (TByte)(0xF0 + (cp >> 18));
    buf[1] = (TByte)(0x80 + ((cp >> 12) & 0x3F));
    buf[2] = (TByte)(0x80 + ((cp >> 6) & 0x3F));
    buf[3] = (TByte)(0x80 + (cp & 0x3F));
    return 4;
  } else return 0;
}

// Возвращает число байтов кодировки code point в UTF-8
size_t encode_bytes(TSymbol cp) {
  if (cp < 0x00) {
    return 0;
  } else if (cp < 0x80) {
    return 1;
  } else if (cp < 0x800) {
    return 2;
  } else if (cp < 0x10000) {
    return 3;
  } else if (cp < 0x110000) {
    return 4;
  } else return 0;
}

#define utf_cont(ch)  (((ch) & 0xc0) == 0x80)

// Извлекает code point из UTF-8 буфера.
// Если успешно: code point записан в cp_ref, возвращает число прочитанных байтов
// Если ошибк, возвращает -1
size_t decode_symbol(TByte* buf, size_t buflen, TSymbol* cp_ref) {
  int32_t cp;
  const TByte *end;

  if (!buflen) return -1;

  *cp_ref = 0;
  end = buf + ((buflen < 0) ? 4 : buflen);
  
  cp = *buf++;
  if (cp < 0x80) {
    *cp_ref = cp;
    return 1;
  }
  // Первый байт должен быть в диапазоне [0xc2..0xf4]
  if ((TSymbol)(cp - 0xc2) > (0xf4-0xc2)) return -1;

  if (cp < 0xe0) {         // 2-byte sequence
     // Must have valid continuation character
     if (buf >= end || !utf_cont(*buf)) return -1;
     *cp_ref = ((cp & 0x1f)<<6) | (*buf & 0x3f);
     return 2;
  }
  if (cp < 0xf0) {        // 3-byte sequence
     if ((buf + 1 >= end) || !utf_cont(*buf) || !utf_cont(buf[1]))
        return -1;
     // Check for surrogate chars
     if (cp == 0xed && *buf > 0x9f)
         return -1;
     cp = ((cp & 0xf)<<12) | ((*buf & 0x3f)<<6) | (buf[1] & 0x3f);
     if (cp < 0x800)
         return -1;
     *cp_ref = cp;
     return 3;
  }
  // 4-byte sequence
  // Must have 3 valid continuation characters
  if ((buf + 2 >= end) || !utf_cont(*buf) || !utf_cont(buf[1]) || !utf_cont(buf[2]))
     return -1;
  // Make sure in correct range (0x10000 - 0x10ffff)
  if (cp == 0xf0) {
    if (*buf < 0x90) return -1;
  } else if (cp == 0xf4) {
    if (*buf > 0x8f) return -1;
  }
  *cp_ref = ((cp & 7)<<18) | ((*buf & 0x3f)<<12) | ((buf[1] & 0x3f)<<6) | (buf[2] & 0x3f);
  return 4;
}


//==== strings

EXPORTED TString tri_newLiteralString(TString* sptr, TInt64 bytes, TInt64 symbols, char* body) {

	if (*sptr != NULL) return *sptr;

	size_t sz = sizeof(StringDesc) + bytes + 1; // +1 для 0x0
//	printf("sz=%lld\n", sz);
	void* mem = mm_allocate(sz);
//	printf("mem=%p\n", mem);

	TString s = mem;
	s->bytes = bytes;
	s->symbols = symbols;
	s->body = mem + sizeof(StringDesc);
	memcpy(s->body, body, bytes);
	s->body[bytes] = 0x0;

	*sptr = s;
	
	return s;
}

EXPORTED TString tri_newString(TInt64 bytes, TInt64 symbols, char* body) {

	size_t sz = sizeof(StringDesc) + bytes + 1; // +1 для 0x0
	void* mem = mm_allocate(sz);

	TString s = mem;
	s->bytes = bytes;
	s->symbols = symbols;
	s->body = mem + sizeof(StringDesc);
	memcpy(s->body, body, bytes);
	s->body[bytes] = 0x0;

	return s;
}

// Initialized in tri_init
StringDesc emptyStringDesc;

EXPORTED TString tri_emptyString() {
    return &emptyStringDesc;
}

// Делает дескриптор, но не копирует содержимое
EXPORTED TString tri_newStringDesc(TInt64 bytes, TInt64 symbols) {

	size_t sz = sizeof(StringDesc) + bytes + 1; // +1 для 0x0
	void* mem = mm_allocate(sz);

	TString s = mem;
	s->bytes = bytes;
	s->symbols = symbols;
	s->body = mem + sizeof(StringDesc);

	return s;
}


EXPORTED TInt64 tri_lenString(TString s) {
	if (s->symbols >= 0) return s->symbols;
	
	TInt64 count = 0;
	TSymbol cp;

	size_t i = 0;
	size_t symlen;
	TByte* buf = s->body;
	while (i < s->bytes) {
		symlen = decode_symbol(buf, s->bytes - i, &cp);
		if (symlen < 0) {
			break;
		}
		count++;
		i += symlen;
		buf += symlen;
	}	

	return count;
}

EXPORTED TBool tri_equalStrings(TString s1, TString s2) {
    if (s1 == s2) {
        return true;
    }
    if (s1->bytes != s2->bytes) {
        return false;
    }
    
    return memcmp(s1->body, s2->body, s1->bytes) == 0;
}

//==== vector


typedef struct VectorDesc { 
    //TODO: Tag
    TInt64 len;
    //TODO: capacity
    void* body; 
} VectorDesc;


EXPORTED void* tri_newVector(size_t element_size, TInt64 len) {
	VectorDesc* v = mm_allocate(sizeof(VectorDesc));
	v->len = len;
	v->body = mm_allocate(element_size * len);
	
	memset(v->body, 0x0, element_size * len);
	return v;
}

EXPORTED void* tri_newVectorDesc() {
	VectorDesc* v = mm_allocate(sizeof(VectorDesc));
	return v;
}

EXPORTED TInt64 tri_lenVector(void* vd) {
	VectorDesc* v = vd;
	return v->len;	
}

EXPORTED TInt64 tri_vcheck(void* vd, TInt64 inx) {
	VectorDesc* v = vd;
	if (inx < 0 || inx >= v->len) {
		runtime_crash("vector index out of bounds");
	}
	
	return inx;
}

//==== class

typedef struct VTMini { size_t self_size; } VTMini;
typedef struct MetaMini { size_t object_size; void* base_desc; } MetaMini;
typedef struct ClassMini { void* vtable; } ClassMini;

EXPORTED void* tri_newObject(void* class_desc) {
	
	VTMini* vt = class_desc;
	size_t vt_sz = vt->self_size;

	MetaMini* m = class_desc + vt_sz;
	size_t o_sz = m->object_size;
	
	ClassMini* c = mm_allocate(o_sz);
	memset(c, 0x0, o_sz);
	c->vtable = vt;
	
	return c;
}

EXPORTED void* tri_checkClassType(void* object, void* target_desc) {
	
	VTMini* current_vt = ((ClassMini*)object)->vtable;

//printf("object_vt = %p\n", current_vt);
//printf("target_desc = %p\n", target_desc);
	
	if (current_vt == target_desc) {
//printf("found self\n");
		return object;
	}
	
	MetaMini* m = (void *)current_vt + current_vt->self_size;
	
	while (m->base_desc != NULL) {
		printf("base_desc = %p\n", m->base_desc);
		
		if (m->base_desc == target_desc) return object;
		
		current_vt = m->base_desc;
		m = (void *)current_vt + current_vt->self_size;
	}
	
	runtime_crash("failed class type check");
	
	return NULL;
}

//==== conversions

EXPORTED TByte tri_TInt64_to_TByte(TInt64 x) {
	if (x < 0 || x > 255) {
		runtime_crash("conversion to byte out of range");
	}
	return (TByte)x;
}

EXPORTED TByte tri_TSymbol_to_TByte(TSymbol x) {
	if (x > 255) {
		runtime_crash("conversion to byte out of range");
	}
	return (TByte)x;
}

EXPORTED TInt64 tri_TFloat64_to_TInt64(TFloat64 x) {
	return llround(x);
}

#define MaxSymbol 0x10FFFF

EXPORTED TSymbol tri_TInt64_to_TSymbol(TInt64 x) {
	if (x < 0 || x > MaxSymbol) {
		runtime_crash("conversion to symbol out of range");
	}
	return (TSymbol)x;	
}

EXPORTED TString tri_TSymbol_to_TString(TSymbol x) {
	TByte buf[8];
	size_t bytes = encode_symbol(x, buf);
	
	return tri_newString(bytes, 1, (char *)buf);
}

EXPORTED TString tri_Bytes_to_TString(void *vd) {
	VectorDesc* v = vd;
	//TODO: check meta and crash

	//TODO: calculate symbols? lazy?
	return tri_newString(v->len, -1, (char *)v->body);	
}	

EXPORTED TString tri_Symbols_to_TString(void *vd) {
	VectorDesc* v = vd;
	//TODO: check meta and crash

	TInt64 bytes = 0;
	TSymbol *symbuf = v->body;
	for (int i = 0; i < v->len; i++) {
		bytes += encode_bytes(symbuf[i]);
	}	

	TString s = tri_newStringDesc(bytes, v->len);
	
	TByte *bytebuf = s->body;
	int len;
	for (int i = 0; i < v->len; i++) {
		len = encode_symbol(symbuf[i], bytebuf);
		bytebuf += len;
	}	

	return s;	
}	

EXPORTED void* tri_TString_to_Bytes(TString s) {
	VectorDesc* v = tri_newVectorDesc();
	
	v->len = s->bytes;
	v->body = mm_allocate(sizeof(TByte) * v->len);
	memcpy(v->body, s->body, s->bytes);
	
	return v;
}

EXPORTED void* tri_TString_to_Symbols(TString s) {
	TInt64 count = 0;
	TSymbol cp;

	size_t i = 0;
	size_t symlen;
	TByte* buf = s->body;
	while (i < s->bytes) {
		symlen = decode_symbol(buf, s->bytes - i, &cp);
		if (symlen < 0) {
			runtime_crash("invalid utf-8 bytes");
			return NULL;
		}
		count++;
		i += symlen;
		buf += symlen;
	}

	VectorDesc* v = tri_newVectorDesc();
	v->len = count;
	v->body = mm_allocate(sizeof(TSymbol) * count);	
	
	TSymbol* symbuf = v->body;
	i = 0;
	buf = s->body;
	while (i < count) {
		symlen = decode_symbol(buf, s->bytes - i, &cp);
		symbuf[i] = cp;
		buf += symlen;
		i++;
	}	
	
	return v;
}

//==== console

/*
void print_int(int i) {
  printf("%d", i);
}
*/

EXPORTED void print_int64(TInt64 i) {
  printf("%lld", i);
}

EXPORTED void print_float64(TFloat64 f) {
  printf("%g", f);
}

EXPORTED void print_symbol(TSymbol s) {
  printf("0x%x", s);
}	

EXPORTED void print_string(TString s) {
  printf("%s", s->body);
}	

EXPORTED void print_bool(TBool b) {
	if (b) printf("истина"); else printf("ложь");
}

EXPORTED void println() {
  printf("\n");
}

//==== crash

EXPORTED void tri_crash(char* msg, char* pos) {
	printf("авария '%s' в позиции %s\n", msg, pos);
    panic();
}

//==== init

EXPORTED void tri_init() {
    emptyStringDesc.bytes = 0;
    emptyStringDesc.symbols = 0;
    emptyStringDesc.body = (TByte*)"";
}
