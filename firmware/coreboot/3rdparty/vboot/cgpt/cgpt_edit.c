/* Copyright 2018 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "cgpt.h"
#include "cgptlib_internal.h"
#include "cgpt_params.h"
#include "vboot_host.h"

int CgptEdit(CgptEditParams *params) {
  struct drive drive;
  GptHeader *h;
  int gpt_retval;

  if (params == NULL)
    return CGPT_FAILED;

  if (CGPT_OK != DriveOpen(params->drive_name, &drive, O_RDWR,
                           params->drive_size))
    return CGPT_FAILED;

  if (GPT_SUCCESS != (gpt_retval = GptSanityCheck(&drive.gpt))) {
    Error("GptSanityCheck() returned %d: %s\n",
          gpt_retval, GptError(gpt_retval));
    goto bad;
  }

  if (CGPT_OK != CheckValid(&drive)) {
    Error("Please run 'cgpt repair' before changing settings.\n");
    goto bad;
  }

  h = (GptHeader *)drive.gpt.primary_header;
  if (params->set_unique) {
    memcpy(&h->disk_uuid, &params->unique_guid, sizeof(h->disk_uuid));
  }
  // Copy to secondary
  RepairHeader(&drive.gpt, MASK_PRIMARY);
  drive.gpt.modified |= (GPT_MODIFIED_HEADER1 | GPT_MODIFIED_HEADER2);

  UpdateCrc(&drive.gpt);

  // Write it all out.
  return DriveClose(&drive, 1);

bad:

  DriveClose(&drive, 0);
  return CGPT_FAILED;
}
