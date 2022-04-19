
#ifndef U_PTR_LIB_H
#define U_PTR_LIB_H

#include <stdio.h>
#include <stdarg.h>
#include <string.h>
#include <stdlib.h>

#include "log.h"

#define ARRAY_LEN_INIT 10
#define ARRAY_LEN_INC  10

#define TRUE  1
#define FALSE 0

typedef void *uPointer; /* void*. */
typedef unsigned int uint;

typedef struct {
  uPointer *pdata; /* array of pointer, void ** */
  uint len; /* number of elements in the array */
  uint alloc; /* allocated elements for the array. allocated >= len */
}uPtrArray;

/* Declerations. */
extern void add_argv(uPtrArray *pArray, ...);
extern void end_argv(uPtrArray *pArray);
extern uPtrArray *create_new_array(int len);
extern int add_elem_to_array(uPtrArray *pArray, char *arg);
extern void free_array(uPtrArray *pArray, int flag);
extern void print_pArray(uPtrArray *pArray);

#endif /* U_PTR_LIB_H */
