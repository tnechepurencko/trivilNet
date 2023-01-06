#include "trirun.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

//==== crash

void crash(char* s) {
	printf("! crash: %s\n", s);
	exit(1);
}

//==== memory

void* mm_allocate(size_t size) {
	void* a = malloc(size);
	if (a == NULL) {
		crash("memory not allocated");
	}
	return a;
}	

//==== utf-8

// Сохраняет code point в буфер, не менее 4 байтов
size_t encode_symbol(TSymbol uc, TByte *buf) {
  if (uc < 0x00) {
    return 0;
  } else if (uc < 0x80) {
    buf[0] = (TByte) uc;
    return 1;
  } else if (uc < 0x800) {
    buf[0] = (TByte)(0xC0 + (uc >> 6));
    buf[1] = (TByte)(0x80 + (uc & 0x3F));
    return 2;
  // Не учитываю диапазон 0xd800-0xdfff, хотя полученный UTF-8 будет не корректным
  } else if (uc < 0x10000) {
    buf[0] = (TByte)(0xE0 + (uc >> 12));
    buf[1] = (TByte)(0x80 + ((uc >> 6) & 0x3F));
    buf[2] = (TByte)(0x80 + (uc & 0x3F));
    return 3;
  } else if (uc < 0x110000) {
    buf[0] = (TByte)(0xF0 + (uc >> 18));
    buf[1] = (TByte)(0x80 + ((uc >> 12) & 0x3F));
    buf[2] = (TByte)(0x80 + ((uc >> 6) & 0x3F));
    buf[3] = (TByte)(0x80 + (uc & 0x3F));
    return 4;
  } else return 0;
}


//==== strings

TString tri_newLiteralString(TString* sptr, TInt64 bytes, TInt64 symbols, char* body) {

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

TString tri_newString(TInt64 bytes, TInt64 symbols, char* body) {

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


TInt64 tri_lenString(TString s) {
	if (s->symbols >= 0) return s->symbols;
	crash("string len ni");
	return 0;
}

//==== vector


typedef struct VectorDesc { TInt64 len; TInt64* body; } VectorDesc;


void* tri_newVector(size_t element_size, TInt64 len) {
	VectorDesc* v = mm_allocate(sizeof(VectorDesc));
	v->len = len;
	v->body = mm_allocate(element_size * len);
	
	memset(v->body, 0x0, element_size * len);
	return v;
}

TInt64 tri_lenVector(void* vd) {
	VectorDesc* v = vd;
	return v->len;	
}

TInt64 tri_vcheck(void* vd, TInt64 inx) {
	VectorDesc* v = vd;
	if (inx < 0 || inx >= v->len) {
		crash("vector index out of bounds");
	}
	
	return inx;
}

//==== class

typedef struct VTMini { size_t self_size; } VTMini;
typedef struct MetaMini { size_t object_size; } MetaMini;
typedef struct ClassMini { void* meta; } ClassMini;

void* tri_newObject(void* meta) {
	
	VTMini* vt = meta;
	size_t vt_sz = vt->self_size;

	MetaMini* m = meta + vt_sz;
	size_t o_sz = m->object_size;
	
	ClassMini* c = mm_allocate(o_sz);
	c->meta = meta;
	
	return c;
}

//==== conversions

TByte tri_TInt64_to_TByte(TInt64 x) {
	if (x < 0 || x > 255) {
		crash("conversion to byte out of range");
	}
	return (TByte)x;
}

TByte tri_TSymbol_to_TByte(TSymbol x) {
	if (x > 255) {
		crash("conversion to byte out of range");
	}
	return (TByte)x;
}

TInt64 tri_TFloat64_to_TInt64(TFloat64 x) {
	return llround(x);
}

#define MaxSymbol 0x10FFFF

TSymbol tri_TInt64_to_TSymbol(TInt64 x) {
	if (x < 0 || x > MaxSymbol) {
		crash("conversion to symbol out of range");
	}
	return (TSymbol)x;	
}

TString tri_TSymbol_to_TString(TSymbol x) {
	TByte buf[8];
	size_t bytes = encode_symbol(x, buf);
	
	return tri_newString(bytes, 1, (char *)buf);
}

//==== console

/*
void print_int(int i) {
  printf("%d", i);
}
*/

void print_int64(TInt64 i) {
  printf("%lld", i);
}

void print_float64(TFloat64 f) {
  printf("%g", f);
}

void print_symbol(TSymbol s) {
  printf("0x%x", s);
}	

void print_string(TString s) {
  printf("%s", s->body);
}	

void println() {
  printf("\n");
}

//==== other

void tri_welcome() {
  printf("Trivil!\n");
}
