#ifndef TEST_STUB_USYS_LOG_H
#define TEST_STUB_USYS_LOG_H
#include <stdio.h>
#define usys_log_debug(...) do { } while (0)
#define usys_log_info(...)  do { } while (0)
#define usys_log_error(...) do { fprintf(stderr, __VA_ARGS__); fprintf(stderr, "\n"); } while (0)
#endif
