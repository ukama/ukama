/* 
 * Simple test program for the uPtrLib
 *
 */

#include "uPtrLib.h"

int main(void){

  /* run few tests and print screen the contents. */

  uPtrArray *p1;

  p1 = create_new_array(20);

  add_argv(p1, "TEST0", "TEST1", "TEST2", NULL);
  
  add_elem_to_array(p1, "TEST1");
  add_elem_to_array(p1, "TEST2");
  add_elem_to_array(p1, "TEST3");

  add_argv(p1, "TEST0.0", "TEST1.0", "TEST2.0", NULL);

  
  print_pArray(p1);

  free_array(p1, 1);
  
  return 0;
}
