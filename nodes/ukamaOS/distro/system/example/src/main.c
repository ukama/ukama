/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * example.d -- example system module using capp
 *
 */

#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "example.h"

int main(int argc, char **argv, char **envp) {

  int loop=FALSE, interval;
  char **env;
  FILE *fp;

  if (argc>1) {
    if (strcmp(argv[1], "-f")==0) { /* loop forever. */
      loop = TRUE;
    }
  }

  /* List all the environment variables */
  for (env = envp; *env != 0; env++) {
    char *thisEnv = *env;
    fprintf(stdout, "%s\n", thisEnv);
  }

  /* Open the config file and read the sleep interval */
  fp = fopen(CONFIG_FILE, "r");
  if (fp == NULL) {
    fprintf(stderr, "Error opening config file at: %s \n", CONFIG_FILE);
    return 1;
  }
  fscanf(fp, "%d", &interval);
  fclose(fp);

  do {
    fprintf(stdout, "Hello world from %s\n", UKAMA_STR);
    sleep(interval);
  } while (loop);

  fprintf(stderr, "... Exiting \n");

  return 1;
}
