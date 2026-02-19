/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Typography } from "@mui/material";

export default function StatTitle({ label, icon }: { label: string; icon: React.ReactNode }) {
  return <Typography variant="overline" sx={{ fontWeight: 700 }}>
    <span>{icon}</span> &nbsp; &nbsp;
    {label}
  </Typography>;
}