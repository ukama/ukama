/* Copyright 2018 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <getopt.h>

#include "cgpt.h"
#include "vboot_host.h"

extern const char* progname;

static void Usage(void)
{
  printf("\nUsage: %s edit [OPTIONS] DRIVE\n\n"
         "Edit a drive's parameters.\n\n"
         "Options:\n"
         "  -D NUM       Size (in bytes) of the disk where partitions reside;\n"
         "                 default 0, meaning partitions and GPT structs are\n"
         "                 both on DRIVE\n"
         "  -u GUID      Drive Unique ID\n"
         "\n", progname);
}

int cmd_edit(int argc, char *argv[]) {

  CgptEditParams params;
  memset(&params, 0, sizeof(params));

  int c;
  int errorcnt = 0;
  char *e = 0;

  opterr = 0;                     // quiet, you
  while ((c=getopt(argc, argv, ":hu:D:")) != -1)
  {
    switch (c)
    {
    case 'D':
      params.drive_size = strtoull(optarg, &e, 0);
      errorcnt += check_int_parse(c, e);
      break;
    case 'u':
      params.set_unique = 1;
      if (CGPT_OK != StrToGuid(optarg, &params.unique_guid)) {
        Error("invalid argument to -%c: %s\n", c, optarg);
        errorcnt++;
      }
      break;
    case 'h':
      Usage();
      return CGPT_OK;
    case '?':
      Error("unrecognized option: -%c\n", optopt);
      errorcnt++;
      break;
    case ':':
      Error("missing argument to -%c\n", optopt);
      errorcnt++;
      break;
    default:
      errorcnt++;
      break;
    }
  }
  if (errorcnt)
  {
    Usage();
    return CGPT_FAILED;
  }

  if (optind >= argc)
  {
    Error("missing drive argument\n");
    return CGPT_FAILED;
  }

  params.drive_name = argv[optind];

  if (!params.set_unique)
  {
    Error("no parameters were edited\n");
    return CGPT_FAILED;
  }

  return CgptEdit(&params);
}
