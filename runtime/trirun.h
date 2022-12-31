#include <stdint.h>
#include <stdbool.h>

typedef uint8_t TByte;
typedef int64_t TInt64;
typedef double TFloat64;
typedef _Bool TBool;


typedef uint32_t TSymbol;

typedef struct StringDesc {
  int64_t lenBytes;
  int64_t lenSymbols;
  uint8_t* body;
} StringDesc;

typedef StringDesc* TString;

//====

void tri_welcome();

//==== strings

TString makeString(uint8_t* body, TInt64 bytes);

//==== console

//void print_int(int i);
void print_int64(TInt64 i);

void println();