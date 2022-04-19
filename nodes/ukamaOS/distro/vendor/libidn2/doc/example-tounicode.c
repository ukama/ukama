/* example-tounicode.c --- Example ToUnicode() code showing how to use Libidn2.
 *
 * This code is placed under public domain.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <idn2.h>		/* idn2_to_unicode_8z8z() */

/*
 * Compiling using pkg-config is recommended:
 *
 * $ cc -o example-to-unicode example-to-unicode.c $(pkg-config --cflags --libs libidn2)
 * $ ./example-tounicode
 * Input domain (ACE) encoded as `UTF-8': xn--nxasmm1c.com
 *
 * Read string (length 16): 78 6e 2d 2d 6e 78 61 73 6d 6d 31 63 2e 63 6f 6d
 * ACE label (length 14): 'βόλος.com'
 *
 */

int
main (void)
{
  char buf[BUFSIZ];
  char *p;
  int rc;
  size_t i;

  if (!fgets (buf, BUFSIZ, stdin))
    perror ("fgets");
  buf[strlen (buf) - 1] = '\0';

  printf ("Read string (length %ld): ", (long int) strlen (buf));
  for (i = 0; i < strlen (buf); i++)
    printf ("%02x ", (unsigned) buf[i] & 0xFF);
  printf ("\n");

  rc = idn2_to_unicode_8z8z (buf, &p, 0);
  if (rc != IDNA_SUCCESS)
    {
      printf ("ToUnicode() failed (%d): %s\n", rc, idn2_strerror (rc));
      return EXIT_FAILURE;
    }

  printf ("ACE label (length %ld): '%s'\n", (long int) strlen (p), p);

  free (p);			/* or idn2_free() */

  return 0;
}
