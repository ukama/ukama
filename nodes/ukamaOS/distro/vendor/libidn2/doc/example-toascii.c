/* example-toascii.c --- Example ToASCII() code showing how to use Libidn2.
 *
 * This code is placed under public domain.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <idn2.h>		/* idn2_to_ascii_8z() */

/*
 * Compiling using pkg-config is recommended:
 *
 * $ cc -o example-toascii example-toascii.c $(pkg-config --cflags --libs libidn2)
 * $ ./example-toascii
 * Input domain encoded as `UTF-8': βόλος.com
 * Read string (length 15): ce b2 cf 8c ce bb ce bf cf 82 2e 63 6f 6d 0a
 * ACE label (length 17): 'xn--nxasmm1c.com'
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

  /* Use non-transitional IDNA2008 */
  rc = idn2_to_ascii_8z (buf, &p, IDN2_NONTRANSITIONAL);
  if (rc != IDNA_SUCCESS)
    {
      printf ("ToASCII() failed (%d): %s\n", rc, idn2_strerror (rc));
      return EXIT_FAILURE;
    }

  printf ("ACE label (length %ld): '%s'\n", (long int) strlen (p), p);

  free (p);			/* or idn2_free() */

  return 0;
}
