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

void* mm_reallocate(void* ptr, size_t size) {
	void* a = realloc(ptr, size);
	if (a == NULL) {
		runtime_crash("memory not reallocated");
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
    
    if (bytes < 0) {
        bytes = strlen(body);
    	//printf("bytes=%lld symbols=%lld\n", bytes, symbols);
    }

	size_t sz = sizeof(StringDesc) + bytes + 1; // +1 для 0x0
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
    TInt64 capacity;
    void* body; 
} VectorDesc;


EXPORTED void* tri_newVector(size_t element_size, TInt64 len, TInt64 cap) {
	VectorDesc* v = mm_allocate(sizeof(VectorDesc));
	v->len = len;
    
    if (cap < len) { cap = len; }
    v->capacity = cap;
    
    if (cap == 0) {
        v->body = NULL;
        return v;
    }

    v->body = mm_allocate(element_size * cap);
    //memset(v->body, 0x0, element_size * cap); //TODO: не надо, см. fill
	
	return v;
}

void vectorFill(VectorDesc* v, size_t element_size, TWord64 filler) {
    
   switch (element_size) {
    case 1: 
        memset(v->body, (int)filler, v->len);
        break;
    case 8:
        TWord64 *a = v->body;
        for (int i = 0; i < v->len; i++) a[i] = filler;
        break;
    default: 
        char buf[128];
        sprintf_s(buf, 128, "vectorFill not implemented for element size=%lld", element_size);
		runtime_crash(buf);    
    }        
}

EXPORTED void* tri_newVectorFill(size_t element_size, TInt64 len, TInt64 cap, TWord64 filler) {
	VectorDesc* v = mm_allocate(sizeof(VectorDesc));
	v->len = len;
    
    if (cap < len) { cap = len; }
    v->capacity = cap;
    
    if (cap == 0) {
        v->body = NULL;
        return v;
    }

    v->body = mm_allocate(element_size * cap);
    vectorFill(v, element_size, filler);
    
	return v;    
}    

EXPORTED void* tri_newVectorDesc() {
	VectorDesc* v = mm_allocate(sizeof(VectorDesc));
	return v;
}

/* //unused 
EXPORTED TInt64 tri_lenVector(void* vd) {
	VectorDesc* v = vd;
	return v->len;	
}
*/

EXPORTED TInt64 tri_indexcheck(TInt64 inx, TInt64 len) {
	if (inx < 0 || inx >= len) {
        char buf[128];
        sprintf_s(buf, 128, "index %lld out of bounds [0..%lld[", inx, len);
		runtime_crash(buf);
	}
	
	return inx;
}    

//=== vector methods

void vectorExtend(VectorDesc* v, size_t element_size, TInt64 new_cap) {
    if (new_cap < v->capacity * 2) new_cap = v->capacity * 2;

    //TODO: нужно копировать по длине (не по capacity)
    v->body = mm_reallocate(v->body, new_cap * element_size);
    v->capacity = new_cap;
}    


EXPORTED void tri_vectorAppend(void* vd, size_t element_size, TInt64 len, void* add_body) {

    if (len <= 0) return;

	VectorDesc* v = vd;
    TInt64 new_len = v->len + len;

    if (new_len > v->capacity) {
        vectorExtend(v, element_size, new_len);
    }    

    //TODO: убрать
    memcpy(v->body + v->len * element_size, add_body, len * element_size);
    
    v->len = new_len;
}

/*
EXPORTED void* tri_vectorFill(void* vd, size_t element_size, TInt64 len, TWord64 filler) {
    	runtime_crash("tri_vectorFill deprecated");
        return NULL;


    if (len <= 0) return vd;

	VectorDesc* v = vd;
    TInt64 new_len = v->len + len;

    if (new_len > v->capacity) {
        vectorExtend(v, element_size, new_len);
    }    

    switch (element_size) {
    case 1: 
        memset(v->body + v->len * element_size, (int)filler, len);
        break;
    case 8:
        TWord64 *a = v->body + v->len * element_size;
        for (int i = 0; i < len; i++) a[i] = filler;
        break;
    default: 
        char buf[128];
        sprintf_s(buf, 128, "vectorFill not implemented for element size=%lld", element_size);
		runtime_crash(buf);    
    }    

    v->len = new_len;

    return vd;
    
}
*/

EXPORTED void tri_vectorAppend_TSymbol_to_Bytes(void *vd, TSymbol x) {

    TByte buf[4];
    size_t len = encode_symbol(x, buf);

	VectorDesc* v = vd;
    TInt64 new_len = v->len + len;

    if (new_len > v->capacity) {
        vectorExtend(v, sizeof(TByte), new_len);
    }    

    memcpy(v->body + v->len * sizeof(TByte), buf, len * sizeof(TByte));
    
    v->len = new_len;
}

//==== nil check

EXPORTED void* tri_nilcheck(void* r) {
    if (r == NULL) {
        runtime_crash("nil check");    
    }
    return r;
}

//==== class

EXPORTED void* tri_newObject(void* class_desc) {
	
	_BaseVT* vt = class_desc;
	size_t vt_sz = vt->self_size;

	_BaseMeta* m = class_desc + vt_sz;
	size_t o_sz = m->object_size;
	
	_BaseObject* o = mm_allocate(o_sz);
	memset(o, 0x0, o_sz);
	o->vtable = vt;
    
    vt->__init__(o);
	
	return o;
}

EXPORTED void* tri_checkClassType(void* object, void* target_desc) {
	
	_BaseVT* current_vt = ((_BaseObject*)object)->vtable;
	
	if (current_vt == target_desc) {
//printf("found self\n");
		return object;
	}
	
	_BaseMeta* m = (void *)current_vt + current_vt->self_size;
	
	while (m->base_desc != NULL) {
		//printf("base_desc = %p\n", m->base_desc);
		
		if (m->base_desc == target_desc) return object;
		
		current_vt = m->base_desc;
		m = (void *)current_vt + current_vt->self_size;
	}
	
	runtime_crash("failed class type check");
	
	return NULL;
}

EXPORTED TBool tri_isClassType(void* object, void* target_desc) {
    	_BaseVT* current_vt = ((_BaseObject*)object)->vtable;
	
	if (current_vt == target_desc) {
//printf("found self\n");
		return true;
	}
	
	_BaseMeta* m = (void *)current_vt + current_vt->self_size;
	
	while (m->base_desc != NULL) {
		//printf("base_desc = %p\n", m->base_desc);
		
		if (m->base_desc == target_desc) return true;
		
		current_vt = m->base_desc;
		m = (void *)current_vt + current_vt->self_size;
	}
	
	return false;
}


//==== conversions

EXPORTED TByte tri_TInt64_to_TByte(TInt64 x) {
	if (x < 0 || x > 255) {
		runtime_crash("conversion to byte out of range");
	}
	return (TByte)x;
}

EXPORTED TByte tri_TWord64_to_TByte(TWord64 x) {
	if (x > 255) {
		runtime_crash("conversion to byte out of range");
	}
	return (TByte)x;
}

EXPORTED TInt64 tri_TWord64_to_TInt64(TWord64 x) {
	if (x > 0x7FFFFFFFFFFFFFFF) {
		runtime_crash("conversion to int64 out of range");
	}
	return (TInt64)x;
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

EXPORTED TSymbol tri_TWord64_to_TSymbol(TWord64 x) {
	if (x > MaxSymbol) {
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
    v->capacity = s->bytes;
	v->body = mm_allocate(sizeof(TByte) * v->len);
	memcpy(v->body, s->body, s->bytes);
	
	return v;
}

EXPORTED void* tri_TSymbol_to_Bytes(TSymbol x) {

    TByte buf[4];
    size_t len = encode_symbol(x, buf);

	VectorDesc* v = tri_newVectorDesc();

 	v->len = len;
    v->capacity = len;
	v->body = mm_allocate(sizeof(TByte) * len);
	memcpy(v->body, buf, len);
	
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
    v->capacity = count;
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

//==== tags

enum Tags {
    tag_unsigned = 1,
    tag_signed,
    tag_float,
    tag_bool,
    tag_symbol,
    tag_string,
    tag_tag,
    
    tag_class
};

#define size_shift 8
#define tag_id_shift 3

#define flag_lang 1

EXPORTED TWord64 tri_tagTByte() {
    return 1 << size_shift | tag_unsigned << tag_id_shift | flag_lang;
}

EXPORTED TWord64 tri_tagTInt64() {
    return (8 << size_shift) | (tag_signed << tag_id_shift) | flag_lang;
}

EXPORTED TWord64 tri_tagTFloat64() {
    return 8 << size_shift | tag_float << tag_id_shift | flag_lang;
}

EXPORTED TWord64 tri_tagTWord64() {
    return 8 << size_shift | tag_unsigned << tag_id_shift | flag_lang;
}

EXPORTED TWord64 tri_tagTBool() {
    return 1 << size_shift | tag_bool << tag_id_shift | flag_lang;
}

EXPORTED TWord64 tri_tagTSymbol() {
    return 4 << size_shift | tag_symbol << tag_id_shift | flag_lang;
}

EXPORTED TWord64 tri_tagTString() {
    return 8 << size_shift | tag_string << tag_id_shift | flag_lang;
}

//==== console

/*
void print_int(int i) {
  printf("%d", i);
}
*/

EXPORTED void print_byte(TByte i) {
  printf("%02x", i);
}

EXPORTED void print_int64(TInt64 i) {
  printf("%lld", i);
}

EXPORTED void print_float64(TFloat64 f) {
  printf("%g", f);
}

EXPORTED void print_word64(TWord64 x) {
  printf("0x%llx", x);
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

//==== аргументы

static int _argc  = 0;
static char **_argv;

EXPORTED TInt64 tri_argc() {
    return _argc;
}

EXPORTED TString tri_arg(TInt64 n) {
    if (n < 0 || n >= _argc) {
        return &emptyStringDesc;
    }
    
    TInt64 bytes = strlen(_argv[n]);
    
    return tri_newString(bytes, -1, _argv[n]);
}    

//==== init/exit

EXPORTED void tri_init(int argc, char *argv[]) {
    
    _argc = argc;
    _argv = argv;
    
    emptyStringDesc.bytes = 0;
    emptyStringDesc.symbols = 0;
    emptyStringDesc.body = (TByte*)"";
}

EXPORTED void tri_exit(TInt64 x) {
    exit(x);
}    
