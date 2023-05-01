#include <stdio.h>
#include <string.h>
#include <errno.h>
#include "rt_sysapi.h"

struct BytesDesc { TInt64 len; TInt64 capacity; TByte* body; };

TString error_id(int errcode) {
    char buf[80];
    
    switch (errcode) {
    case ENOENT: strcpy_s(buf, 80, "ФАЙЛ-НЕ-НАЙДЕН"); break;
    default:
        sprintf(buf, "ОШИБКА[%d]", errcode); 
    }
    return tri_newString(strlen(buf), -1, buf);
}

EXPORTED void* sysapi_fread(void* request, TString filename) {
    
    struct Request* req = request;
    
    FILE* fp;

    int errcode = fopen_s(&fp, (char *)filename->body, "rb");
    if (errcode != 0) {
        req->err_id = error_id(errcode);
        return NULL;
    }

    fseek(fp, 0, SEEK_END);
    size_t sz = ftell(fp);
    rewind(fp);
    
    struct BytesDesc* bytes = tri_newVector(sizeof(TByte), sz, 0); 
    
    size_t ret = fread(bytes->body, sizeof(TByte), sz, fp);
    
    if (ret != sz) {
            req->err_id = error_id(ferror(fp));
            // TODO: удалить bytes?
            fclose(fp);
            return NULL;
    }

    fclose(fp);
    req->err_id = NULL;
    return bytes;
}

EXPORTED void sysapi_fwrite(void* request, TString filename, void* bytes) {
    
    struct Request* req = request;
    
    FILE* fp;

    int errcode = fopen_s(&fp, (char *)filename->body, "wb");
    if (errcode != 0) {
        req->err_id = error_id(errcode);
        return;
    }
    
    struct BytesDesc* v = bytes;
    
    size_t ret = fwrite(v->body, sizeof(TByte), v->len, fp);
    
    if (ret != v->len) {
            req->err_id = error_id(ferror(fp));
            // TODO: удалить bytes?
            fclose(fp);
            return;
    }

    fclose(fp);
    req->err_id = NULL;
}
