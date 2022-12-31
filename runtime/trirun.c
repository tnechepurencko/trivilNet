#include "trirun.h"
#include <stdio.h>
#include <stdlib.h>

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

TString makeString(uint8_t* body, TInt64 bytes) {
/*	
	TString s = mm_allocate(sizeof(TStringDesc));
	s->lenBytes = bytes;
	s->lenSymbols = -1;
*/	
	
	return NULL;
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


void println() {
  printf("\n");
}

//==== other

void tri_welcome() {
  printf("Trivil!\n");
}
