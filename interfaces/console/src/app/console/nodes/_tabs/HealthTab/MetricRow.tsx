/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box, Typography } from "@mui/material";

const BODY_SX = {
  fontFamily: "monospace",
  fontSize: "0.8125rem",
  lineHeight: 1.6,
} as const;

export default function MetricRow({
  label,
  value,
  impact,
}: Readonly<{ label: string; value: string; impact?: string }>) {
  return (
    <Box sx={{ mb: 1.5, "&:last-of-type": { mb: 0 } }}>
      <Typography sx={BODY_SX}>
        {label}: {value}
      </Typography>
      {impact && (
        <Typography sx={{ ...BODY_SX, fontSize: "0.75rem", mt: 0.5 }}>
          Service impact: {impact}
        </Typography>
      )}
    </Box>
  );
}
