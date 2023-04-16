#include <signal.h>
#include <stdio.h>
#include "rt_defs.h"

// ==============   linux     ==============

#ifdef __linux__



// ============== windows ==============

#else
    
#include <windows.h>

LONG WINAPI windows_crash_handler(EXCEPTION_POINTERS* ExceptionInfo) {
    switch (ExceptionInfo->ExceptionRecord->ExceptionCode) {
    case EXCEPTION_ACCESS_VIOLATION:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_ACCESS_VIOLATION");
        break;
    case EXCEPTION_ARRAY_BOUNDS_EXCEEDED:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_ARRAY_BOUNDS_EXCEEDED");
        break;
    case EXCEPTION_BREAKPOINT:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_BREAKPOINT");
        break;
    case EXCEPTION_DATATYPE_MISALIGNMENT:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_DATATYPE_MISALIGNMENT");
        break;
    case EXCEPTION_FLT_DENORMAL_OPERAND:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_FLT_DENORMAL_OPERAND");
        break;
    case EXCEPTION_FLT_DIVIDE_BY_ZERO:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_FLT_DIVIDE_BY_ZERO");
        break;
    case EXCEPTION_FLT_INEXACT_RESULT:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_FLT_INEXACT_RESULT");
        break;
    case EXCEPTION_FLT_INVALID_OPERATION:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_FLT_INVALID_OPERATION");
        break;
    case EXCEPTION_FLT_OVERFLOW:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_FLT_OVERFLOW");
        break;
    case EXCEPTION_FLT_STACK_CHECK:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_FLT_STACK_CHECK");
        break;
    case EXCEPTION_FLT_UNDERFLOW:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_FLT_UNDERFLOW");
        break;
    case EXCEPTION_ILLEGAL_INSTRUCTION:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_ILLEGAL_INSTRUCTION");
        break;
     case EXCEPTION_IN_PAGE_ERROR:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_IN_PAGE_ERROR");
        break;
    case EXCEPTION_INT_DIVIDE_BY_ZERO:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_INT_DIVIDE_BY_ZERO");
        break;
    case EXCEPTION_INT_OVERFLOW:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_INT_OVERFLOW");
        break;
    case EXCEPTION_INVALID_DISPOSITION:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_INVALID_DISPOSITION");
        break;
    case EXCEPTION_NONCONTINUABLE_EXCEPTION:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_NONCONTINUABLE_EXCEPTION");
        break;
    case EXCEPTION_PRIV_INSTRUCTION:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_PRIV_INSTRUCTION");
        break;
    case EXCEPTION_SINGLE_STEP:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_SINGLE_STEP");
        break;
    case EXCEPTION_STACK_OVERFLOW:
        fprintf(stderr, "Sys crash: %s\n", "EXCEPTION_STACK_OVERFLOW");
        break;
/*
    case :
        fprintf(stderr, "Sys crash: %s\n", "");
        break;        
*/        
    default:
        fprintf(stderr, "Unrecognized exception: %lu\n", ExceptionInfo->ExceptionRecord->ExceptionCode);
        break;
    }
    
    return EXCEPTION_EXECUTE_HANDLER;
}

EXPORTED void register_default_crash_handler() {
        SetUnhandledExceptionFilter(windows_crash_handler);
}

#endif
