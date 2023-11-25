#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <regex.h>

#define ARRAY_SIZE(arr) (sizeof((arr)) / sizeof((arr)[0]))

static const char *str = "http://www.google.com";

//static const char *re = "John.*o";
static const char *re =  "((http|https)://)(www.)?[a-zA-Z0-9@:%._\\+~#?&//=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%._\\+~#?&//=]*)";


int main(void)
{
  const char *s = str;
  regex_t     regex;
  regmatch_t  pmatch[1];
  regoff_t    off, len;
  
  if (regcomp(&regex, re, REG_NOSUB))
    exit(EXIT_FAILURE);
  
  printf("String = \"%s\"\n", str);
  printf("Matches:\n");

  if (regexec(&regex, s, 1, pmatch, 0) == 0) {
    printf("Match \n");
  } else {
    printf("No match \n");
  }

#if 0  
  for (int i = 0; ; i++) {
    if (regexec(&regex, s, ARRAY_SIZE(pmatch), pmatch, 0))
      break;
    
    off = pmatch[0].rm_so + (s - str);
    len = pmatch[0].rm_eo - pmatch[0].rm_so;
    printf("#%d:\n", i);
    printf("offset = %jd; length = %jd\n", (intmax_t) off,
	   (intmax_t) len);
    printf("substring = \"%.*s\"\n", len, s + pmatch[0].rm_so);
    
    s += pmatch[0].rm_eo;
  }
#endif
	  
  exit(EXIT_SUCCESS);
}
