#include <stdio.h>
#include "rt_sysapi.h"

struct BytesDesc { TInt64 len; TInt64 capacity; TByte* body; };

EXPORTED void* sysapi_fread(struct SysFiles* sf, TString filename) {
    printf("fread %s\n", filename->body);
    sf->errcode = tri_emptyString();
    return NULL;
}