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

// Возвращает дескриптор байтового массива или NULL, в случае ошибки
// Выставляет код ошибки
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
    
    printf("size %zu\n", sz);    
    
    struct BytesDesc* bytes = tri_newVector(sizeof(TByte), sz, 0); 
    
    size_t ret = fread(bytes->body, sizeof(TByte), sz, fp);
    printf("read %zu\n", ret);    
    
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