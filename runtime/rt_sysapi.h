#ifndef _rt_sysapi_H
#define _rt_sysapi_H

#include "rt_api.h"

struct SysFiles {
    _BaseObject _base;
    //FILE* handler;
    TString errcode;
};

// Возвращает дескриптор байтового массива или NULL, в случае ошибки
EXPORTED void* sysapi_fread(struct SysFiles* sf, TString filename);

#endif