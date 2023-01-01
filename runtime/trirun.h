#include <stdint.h>
#include <stdbool.h>

typedef uint8_t TByte;
typedef int64_t TInt64;
typedef double TFloat64;
typedef _Bool TBool;


typedef uint32_t TSymbol;

typedef struct StringDesc {
  int64_t bytes;
  int64_t symbols;
  uint8_t* body;
} StringDesc;

typedef StringDesc* TString;

//==== strings

TString tri_newLiteralString(TString* sptr, TInt64 bytes, TInt64 symbols, char* body);

TInt64 tri_lenString(TString s);

//==== console

//void print_int(int i);
void print_int64(TInt64 i);

void print_string(TString s);

void println();

//==== other

void tri_welcome();