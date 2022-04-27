/* key-value test code. */

#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#define WIMC_EP_CONTAINER     "/container"
#define WIMC_QUERY_KEY_NAME   "name"
#define WIMC_MAX_CON_NAME_LEN 128 

#define TRUE 1
#define FALSE 0

int get_key_value(char *str, char *key, char *value) {

  char *k, *v;

  k = strtok(str, "=");
  v = strtok(NULL, "=");

  if (k == NULL || v == NULL) {
    goto failure;
  }
  
  if (strcasecmp(WIMC_QUERY_KEY_NAME, k)==0) {

    if (strlen(v) >= WIMC_MAX_CON_NAME_LEN) {
      return FALSE;
    }

    strncpy(value, v, strlen(v));
    strncpy(key, k, strlen(k));
    
    return TRUE;
  } 

 failure:
  return FALSE;
}


int main(int argc, char **argv) {

  char str1[] = "name=value1";
  char str2[]  = "name=valueeeee";

  char key[128], value[128];

  if (get_key_value(str1, &key[0], &value[0])) {
    printf("Key: %s \nValue: %s\n", key, value);
  } else {
    printf("Error\n");
  }

  return 0;
}
