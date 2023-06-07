#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>
#include "rt_sysapi.h"

struct BytesDesc { TInt64 len; TInt64 capacity; TByte* body; };


//=== вещественные (временно)

TBool sysapi_string_to_float64(TString s, TFloat64* res)  {
    
    char *eptr;
    *res = strtod((char *)s->body, &eptr);
    if (eptr != NULL) return false;

    /* If the result is 0, test for an error */
    if (*res == 0)
    {
        /* If the value provided was out of range, display a warning message */
        if (errno == ERANGE || errno == EINVAL) return false;
    }
    return true;
}

//==== коды ошибок общие ===

TString error_id(int errcode) {
    char buf[80];
    
    switch (errcode) {
    case ENOENT: strcpy_s(buf, 80, "ФАЙЛ-НЕ-НАЙДЕН"); break;
    default:
        sprintf(buf, "ОШИБКА[%d]", errcode); 
    }
    return tri_newString(strlen(buf), -1, buf);
}

//==== папки ====

EXPORTED TString sysapi_exec_path() {
    TString folder = tri_arg(0);
    
#if defined(_WIN32) || defined(_WIN64)
    size_t len = folder->bytes;
    char* s = nogc_alloc(len + 1);
    strncpy_s(s, len+1, (char*)folder->body, len);    
    s[len] = 0; 

    for (int i = 0; i < folder->bytes; i++) {
        if (s[i] == '\\') s[i] = '/';
    }   
    folder =  tri_newString(len, -1, s); 
    nogc_free(s);
#endif

    return folder;
}

//==== чтение/запись ====

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

// ==============   linux     ==============

#ifndef _WIN32

EXPORTED TBool sysapi_is_dir(void* request, TString filename)  {
    struct Request* req = request;    
    
    char buf[80];
    sprintf(buf, "sysapi_is_dir не реализована"); 
    req->err_id = tri_newString(strlen(buf), -1, buf);
    
    return false;
}

EXPORTED void* sysapi_dirnames(void* request, TString filename)  {
    
    struct Request* req = request;    
    
    char buf[80];
    sprintf(buf, "sysapi_dirnames не реализована"); 

    req->err_id = tri_newString(strlen(buf), -1, buf);
    return NULL;
}

EXPORTED TString sysapi_abs_path(void* request, TString filename) {
    struct Request* req = request;    
    
    char buf[80];
    sprintf(buf, "sysapi_abs_path не реализована"); 

    req->err_id = tri_newString(strlen(buf), -1, buf);
    return NULL;
}

// ============== windows ==============
#else
    
#include <windows.h>

TString win_error_id(int errcode) {
    char buf[80];
    
    switch (errcode) {
//    case ENOENT: strcpy_s(buf, 80, "ФАЙЛ-НЕ-НАЙДЕН"); break;
    default:
        sprintf(buf, "ОШИБКА[%d]", errcode); 
    }
    return tri_newString(strlen(buf), -1, buf);
}

EXPORTED TBool sysapi_is_dir(void* request, TString filename)  {
    struct Request* req = request;    
    
    WIN32_FIND_DATA FindFileData;
    
    HANDLE hFind = FindFirstFileA((char*)filename->body, &FindFileData);
    if (hFind == INVALID_HANDLE_VALUE) {
        req->err_id = win_error_id(GetLastError());
        return false;
    } 
    
    if (FindFileData.dwFileAttributes & FILE_ATTRIBUTE_DIRECTORY) {
        return true;
    }
    return false;
}

EXPORTED void* sysapi_dirnames(void* request, TString filename)  {
    
     //printf ("dirnames %s\n", filename->body);    
    
    struct Request* req = request;    
    
    // подготовить образец
    size_t len = filename->bytes;
    char* pattern = nogc_alloc(len + 3);
    
    strncpy_s(pattern, len+3, (char*)filename->body, len);
    pattern[len+0] = '\\';
    pattern[len+1] = '*';
    pattern[len+2] = 0;
    
    for (int i = 0; i < len; i++) {
        if (pattern[i] == '/') pattern[i] = '\\';
    }    

//    printf ("!!!%s\n", pattern);    
    
    //=== Считаю число имен
    WIN32_FIND_DATA FindFileData;
    TInt64 count = 0;
    
    HANDLE hFind = FindFirstFileA(pattern, &FindFileData);
    if (hFind == INVALID_HANDLE_VALUE) {
        req->err_id = win_error_id(GetLastError());
        nogc_free(pattern);
        return NULL;
    } 

    do {
        if (strcmp(FindFileData.cFileName, ".") == 0 || strcmp(FindFileData.cFileName, "..") == 0) {
            // игнорирую
        } else {
            count++;
        }
        //printf("!name: %s\n", FindFileData.cFileName);
    } while (FindNextFile(hFind, &FindFileData) != 0);
    FindClose(hFind);

    //=== Собираю вектор
    hFind = FindFirstFileA(pattern, &FindFileData);
    if (hFind == INVALID_HANDLE_VALUE) {
        req->err_id = win_error_id(GetLastError());
        nogc_free(pattern);
        return NULL;
    } 
    
    void* list = tri_newVector(sizeof(TString), 0, count);

    do {
        if (strcmp(FindFileData.cFileName, ".") == 0 || strcmp(FindFileData.cFileName, "..") == 0) {
            // игнорирую
        } else {
            TString s = tri_newString(strlen(FindFileData.cFileName), -1, FindFileData.cFileName);
            tri_vectorAppend(list, sizeof(TString), 1, &s);            
        }
        //printf("!!name: %s\n", FindFileData.cFileName);
    } while (FindNextFile(hFind, &FindFileData) != 0);
    FindClose(hFind);

    req->err_id = NULL;
    nogc_free(pattern);
    return list;
}

EXPORTED TString sysapi_abs_path(void* request, TString filename) {
    struct Request* req = request;    
    
    DWORD retval = GetFullPathNameA(
        (char*)filename->body,
        0,
        NULL,
        NULL);
    
    if (retval == 0) {
        req->err_id = win_error_id(GetLastError());
        return NULL;
    }

    char* buf = nogc_alloc(retval);

    retval = GetFullPathNameA(
        (char*)filename->body,
        retval,
        buf,
        NULL);

    if (retval == 0) {
        tri_crash("sysapi_abs_path: assert 2nd call of GetFullPathNameA returns error", "");
        return NULL;
    }
    
    for (int i = 0; i < retval; i++) {
        if (buf[i] == '\\') buf[i] = '/';
    }       
    
    TString full =  tri_newString(retval, -1, buf); 
    nogc_free(buf);
    
    return full;
}

#endif
