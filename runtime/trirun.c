#include "trirun.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

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

//==== strings

TString tri_newLiteralString(TString* sptr, TInt64 bytes, char* body) {

	if (*sptr != NULL) return *sptr;

	size_t sz = sizeof(StringDesc) + bytes + 1; // +1 для 0x0
//	printf("sz=%lld\n", sz);
	void* mem = mm_allocate(sz);
//	printf("mem=%p\n", mem);

	TString s = mem;
	s->bytes = bytes;
	s->symbols = -1;
	s->body = mem + sizeof(StringDesc);
	memcpy(s->body, body, bytes);
	s->body[bytes] = 0x0;

	*sptr = s;
	
	return s;
}

TInt64 tri_lenString(TString s) {
	if (s->symbols >= 0) return s->symbols;
	crash("string len ni");
	return 0;
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
