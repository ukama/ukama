#ifndef TEST_STUB_USYS_MEM_H
#define TEST_STUB_USYS_MEM_H
#include <stdlib.h>
#endif
#ifndef usys_free
#define usys_free(ptr) do { free(ptr); } while (0)
#endif
