#include <stdint.h>
#include <stdbool.h>
#include "rt_defs.h"

typedef uint8_t TByte;
typedef int64_t TInt64;
typedef double TFloat64;
typedef uint64_t TWord64;
typedef _Bool TBool;
typedef uint32_t TSymbol;

// для преобразования с сохранением битов
typedef union {TFloat64 f; TInt64 i; TWord64 w; } TUnion64;

typedef struct StringDesc {
//TODO meta
  int64_t bytes;
  int64_t symbols;
  TByte* body; // использовать смещение, убрать лишнее обращение к памяти
} StringDesc;

typedef StringDesc* TString;

//==== strings

EXPORTED TString tri_newLiteralString(TString* sptr, TInt64 bytes, TInt64 symbols, char* body);

EXPORTED TInt64 tri_lenString(TString s);

EXPORTED TString tri_emptyString();

EXPORTED TBool tri_equalStrings(TString s1, TString s2); 

//==== vector

EXPORTED void* tri_newVector(size_t element_size, TInt64 len);

EXPORTED TInt64 tri_lenVector(void* vd);

EXPORTED TInt64 tri_indexcheck(TInt64 inx, TInt64 len);

EXPORTED void tri_vectorAppend(void* vd, size_t element_size, TInt64 len, void* body);

// Добавляет символ к []Байт
// Используется строковой библиотекой, не используется компилятором
EXPORTED void tri_vectorAppend_TSymbol_to_Bytes(void *vd, TSymbol x);

//==== class

/*
  object -> vtable ---> vtable size
			fields		(vtable fn)*
						------------ meta info
						object size (for allocation)
						other meta info

*/

EXPORTED void* tri_newObject(void* class_desc);

EXPORTED void* tri_checkClassType(void* object, void* class_desc);

//==== conversions

EXPORTED TByte tri_TInt64_to_TByte(TInt64 x);
EXPORTED TByte tri_TSymbol_to_TByte(TSymbol x);
EXPORTED TInt64 tri_TFloat64_to_TInt64(TFloat64 x);
EXPORTED TSymbol tri_TInt64_to_TSymbol(TInt64 x);

EXPORTED TString tri_TSymbol_to_TString(TSymbol x);

// Параметр []Байт
EXPORTED TString tri_Bytes_to_TString(void* vd);
// Параметр []Символ
EXPORTED TString tri_Symbols_to_TString(void* vd);

// Возвращает []Байт
EXPORTED void* tri_TString_to_Bytes(TString s);
EXPORTED void* tri_TSymbol_to_Bytes(TSymbol x);

// Возвращает []Символ
EXPORTED void* tri_TString_to_Symbols(TString s);

//==== tags

EXPORTED TWord64 tri_tagTByte();
EXPORTED TWord64 tri_tagTInt64();
EXPORTED TWord64 tri_tagTFloat64();
EXPORTED TWord64 tri_tagTWord64();
EXPORTED TWord64 tri_tagTBool();
EXPORTED TWord64 tri_tagTSymbol();
EXPORTED TWord64 tri_tagTString();

//==== console

//void print_int(int i);
EXPORTED void print_byte(TByte i);
EXPORTED void print_int64(TInt64 i);
EXPORTED void print_float64(TFloat64 f);
EXPORTED void print_word64(TWord64 w);

EXPORTED void print_symbol(TSymbol s);
EXPORTED void print_string(TString s);
EXPORTED void print_bool(TBool b);

EXPORTED void println();

//==== crash

EXPORTED void tri_crash(char* msg, char* pos);

//==== init

EXPORTED void tri_init();
