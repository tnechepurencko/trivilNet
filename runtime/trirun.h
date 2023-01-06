#include <stdint.h>
#include <stdbool.h>

typedef uint8_t TByte;
typedef int64_t TInt64;
typedef double TFloat64;
typedef _Bool TBool;
typedef uint32_t TSymbol;

typedef struct StringDesc {
//TODO meta
  int64_t bytes;
  int64_t symbols;
  uint8_t* body;
} StringDesc;

typedef StringDesc* TString;

//==== strings

TString tri_newLiteralString(TString* sptr, TInt64 bytes, TInt64 symbols, char* body);

TInt64 tri_lenString(TString s);

//==== vector

void* tri_newVector(size_t element_size, TInt64 len);

TInt64 tri_lenVector(void* vd);

TInt64 tri_vcheck(void* vd, TInt64 inx);

//==== class

void* tri_newObject(void* meta);


//==== conversions

TByte tri_TInt64_to_TByte(TInt64 x);
TByte tri_TSymbol_to_TByte(TSymbol x);
TInt64 tri_TFloat64_to_TInt64(TFloat64 x);
TSymbol tri_TInt64_to_TSymbol(TInt64 x);

TString tri_TSymbol_to_TString(TSymbol x);

//==== console

//void print_int(int i);
void print_int64(TInt64 i);
void print_float64(TFloat64 f);

void print_symbol(TSymbol s);
void print_string(TString s);

void println();

//==== other

void tri_welcome();