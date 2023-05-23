#include <signal.h>
#include <stdio.h>
#include "rt_defs.h"

// ==============   linux     ==============

#ifndef _WIN32

#include <err.h>
#include <execinfo.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdarg.h>

void die(const char *fmt, ...) {
  va_list ap;

  va_start(ap, fmt);
  vfprintf(stderr, fmt, ap);
  va_end(ap);
  exit(1);
}

#define BT_BUF_SIZE 100

void posix_print_stack_trace() {
  int nptrs;
  void *buffer[BT_BUF_SIZE];
  char **strings;

  nptrs = backtrace(buffer, BT_BUF_SIZE);

  strings = backtrace_symbols(buffer, nptrs);
  if (strings == NULL) {
    die("backtrace_symbols");
    exit(EXIT_FAILURE);
  }

  // Пропускаем вызовы posix_print_stack_trace() и системные вызовы.
  size_t start = 0;
  if (nptrs > 3) {
    start = 3;
  }

  for (size_t j = start; j < nptrs; j++) {
    printf("%s\n", strings[j]);
  }

  free(strings);
}

void posix_signal_handler(int sig, siginfo_t *siginfo, void *context) {
  (void)context;
  switch(sig) {
  case SIGSEGV:
    fputs("Caught SIGSEGV: Segmentation Fault\n", stderr);
    break;
  case SIGINT:
    fputs("Caught SIGINT: Interactive attention signal, (usually ctrl+c)\n", stderr);
    break;
  case SIGFPE:
    switch(siginfo->si_code) {
    case FPE_INTDIV:
      fputs("Caught SIGFPE: (integer divide by zero)\n", stderr);
      break;
    case FPE_INTOVF:
      fputs("Caught SIGFPE: (integer overflow)\n", stderr);
      break;
    case FPE_FLTDIV:
      fputs("Caught SIGFPE: (floating-point divide by zero)\n", stderr);
      break;
    case FPE_FLTOVF:
      fputs("Caught SIGFPE: (floating-point overflow)\n", stderr);
      break;
    case FPE_FLTUND:
      fputs("Caught SIGFPE: (floating-point underflow)\n", stderr);
      break;
    case FPE_FLTRES:
      fputs("Caught SIGFPE: (floating-point inexact result)\n", stderr);
      break;
    case FPE_FLTINV:
      fputs("Caught SIGFPE: (floating-point invalid operation)\n", stderr);
      break;
    case FPE_FLTSUB:
      fputs("Caught SIGFPE: (subscript out of range)\n", stderr);
      break;
    default:
      fputs("Caught SIGFPE: Arithmetic Exception\n", stderr);
      break;
    }
    break;
  case SIGILL:
    switch(siginfo->si_code) {
    case ILL_ILLOPC:
      fputs("Caught SIGILL: (illegal opcode)\n", stderr);
      break;
    case ILL_ILLOPN:
      fputs("Caught SIGILL: (illegal operand)\n", stderr);
      break;
    case ILL_ILLADR:
      fputs("Caught SIGILL: (illegal addressing mode)\n", stderr);
      break;
    case ILL_ILLTRP:
      fputs("Caught SIGILL: (illegal trap)\n", stderr);
      break;
    case ILL_PRVOPC:
      fputs("Caught SIGILL: (privileged opcode)\n", stderr);
      break;
    case ILL_PRVREG:
      fputs("Caught SIGILL: (privileged register)\n", stderr);
      break;
    case ILL_COPROC:
      fputs("Caught SIGILL: (coprocessor error)\n", stderr);
      break;
    case ILL_BADSTK:
      fputs("Caught SIGILL: (internal stack error)\n", stderr);
      break;
    default:
      fputs("Caught SIGILL: Illegal Instruction\n", stderr);
      break;
    }
    break;
  case SIGTERM:
    fputs("Caught SIGTERM: a termination request was sent to the program\n", stderr);
    break;
  case SIGABRT:
    fputs("Caught SIGABRT: usually caused by an abort() or assert()\n", stderr);
    break;
  default:
    break;
  }
  posix_print_stack_trace();
  exit(1);
}

static uint8_t alternate_stack[SIGSTKSZ];
void set_signal_handler() {
  /* setup alternate stack */
  {
    stack_t ss = {};
    /* malloc is usually used here, I'm not 100% sure my static allocation
       is valid but it seems to work just fine. */
    ss.ss_sp = (void*)alternate_stack;
    ss.ss_size = SIGSTKSZ;
    ss.ss_flags = 0;

    if (sigaltstack(&ss, NULL) != 0) {
      err(1, "sigaltstack");
    }
  }

  /* register our signal handlers */
  {
    struct sigaction sig_action = {};
    sig_action.sa_sigaction = posix_signal_handler;
    sigemptyset(&sig_action.sa_mask);

#ifdef __APPLE__
    /* for some reason we backtrace() doesn't work on osx
       when we use an alternate stack */
    sig_action.sa_flags = SA_SIGINFO;
#else
    sig_action.sa_flags = SA_SIGINFO | SA_ONSTACK;
#endif

    if (sigaction(SIGSEGV, &sig_action, NULL) != 0) { err(1, "sigaction"); }
    if (sigaction(SIGFPE,  &sig_action, NULL) != 0) { err(1, "sigaction"); }
    if (sigaction(SIGINT,  &sig_action, NULL) != 0) { err(1, "sigaction"); }
    if (sigaction(SIGILL,  &sig_action, NULL) != 0) { err(1, "sigaction"); }
    if (sigaction(SIGTERM, &sig_action, NULL) != 0) { err(1, "sigaction"); }
    if (sigaction(SIGABRT, &sig_action, NULL) != 0) { err(1, "sigaction"); }
  }
}

EXPORTED void register_default_crash_handler() {
  set_signal_handler();
}

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
