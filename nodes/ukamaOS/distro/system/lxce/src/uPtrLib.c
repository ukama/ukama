/*
 * Generic helper function to handle array of pointers.
 *
 */

#include "uPtrLib.h"

/* add_argv -- add each args to the Array of pointer. Last arg needs to be
 *             NULL.
 *
 */

void add_argv(uPtrArray *pArray, ...) {
  
  va_list args;
  char *arg;
  
  va_start(args, pArray);
  
  while ((arg = va_arg(args, char *)))
    add_elem_to_array(pArray, arg);
  
  va_end(args);
}

void end_argv(uPtrArray *pArray) {
  add_elem_to_array(pArray, NULL);
}

/*
 * create_new_array -- Initialize and create new array pointer of size len.
 *
 *
 */

uPtrArray *create_new_array(int len) {
  
  int elm;
  uPtrArray *pArray = NULL;
  
  if (!len) {
    elm = ARRAY_LEN_INIT;
  } else {
    elm = len;
  }

  pArray = (uPtrArray *)malloc(sizeof(uPtrArray));

  if (!pArray) {
    fprintf(stderr, "Error allocating new array of size: %d\n",
	    sizeof(uPtrArray));
    return NULL;
  }
   
  pArray->pdata = (uPointer *)malloc(sizeof(uPointer)*elm);

  if (!pArray->pdata) {
    fprintf(stderr, "Error allocating new array of size: %d\n",
	    sizeof(uPointer)*elm);
    free(pArray);
    return NULL;
  }

  pArray->alloc = elm;
  pArray->len   = 0; /* Empty but pre-allocated elm */

  return pArray;
}

/*
 * add_elem_to_array -- Add new element to the array. 
 *                      memory will be re-allocated, if needed, in which case
 *                      the pointer will be updated and old allocation will be 
 *                      freed.
 *
 */

int add_elem_to_array(uPtrArray *pArray, char *arg) {

  int len;
  char *dest = NULL;

  /* Check to see if we need to reallocate */
  if (pArray->alloc == pArray->len) {
    /* XXX fix me. */
    
    uPtrArray *newpArray = create_new_array(pArray->alloc * ARRAY_LEN_INC);

    for (int i=0; i<pArray->len; i++) {
      newpArray->pdata[i] = pArray->pdata[i];
    }
    newpArray->len = pArray->len;

    free_array(pArray, FALSE);
    pArray = newpArray;
  }

  if (arg == NULL) { /* null terminated */
    pArray->pdata[pArray->len] = NULL;
  } else {
    
    len = strlen(arg)+1;
    dest = (char *)malloc(len);
    
    if (!dest) {
      fprintf(stderr, "Error allocating memory of size: %d\n", len);
      return FALSE;
    }
    
    strncpy(dest, arg, len);
    pArray->pdata[pArray->len] = dest;
  }
  
  pArray->len++;
  
  return TRUE;
}

/*
 * free_array -- allocated the pointer array. If flag is set, also free 
 *               uPointer.
 */

void free_array(uPtrArray *pArray, int flag) {

  if (flag) {
    for (int i=0; i<pArray->len; i++) {
      free(pArray->pdata[i]);
    }
  }
  
  free(pArray->pdata);
  free(pArray);

  pArray = NULL;
}

/*
 * print_pArray -- dump on the screen.
 *
 */

void print_pArray(uPtrArray *pArray) {

  uPointer *ptr = pArray->pdata;
  int count;

  for (count=0; count<pArray->len; count++) {
    fprintf(stdout, "%d: %s\n", count, (char *)ptr[count]);
  }

  fprintf(stdout, "total alloc: %d len: %d \n", pArray->alloc, pArray->len);
  
}
