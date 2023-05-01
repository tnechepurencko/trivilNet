#ifndef _rt_sysapi_H
#define _rt_sysapi_H

#include "rt_api.h"

// Используется 
struct Request {
    _BaseObject _base;
    //FILE* handler;
    TString err_id;
};

// Возвращает дескриптор байтового массива или NULL, в случае ошибки
// Выставляет код ошибки
EXPORTED void* sysapi_fread(void* request, TString filename);

#endif
