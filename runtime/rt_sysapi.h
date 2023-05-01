#ifndef _rt_sysapi_H
#define _rt_sysapi_H

#include "rt_api.h"

// Используется 
struct Request {
    _BaseObject _base;
    //FILE* handler;
    TString err_id;
};

// Читает файл, возвращает дескриптор байтового вектора.
// В случае ошибки, возвращает NULL и выставляет код ошибки в запросе 
EXPORTED void* sysapi_fread(void* request, TString filename);

// Записывает в файл байтовый вектор.
// В случае ошибки выставляет код ошибки в запросе
EXPORTED void sysapi_fwrite(void* request, TString filename, void* bytes);

// Выдает список имен в папке - список строк []Строка
// В случае ошибки выставляет код ошибки в запросе
EXPORTED void* sysapi_dirnames(void* request, TString filename) ;

#endif
