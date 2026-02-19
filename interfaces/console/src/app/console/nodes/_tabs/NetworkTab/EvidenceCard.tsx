/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";

const EVIDENCE_ITEMS = [
  "Cellular throughput flat at ~20 Mbps for 25 minutes",
  "Backhaul utilization matches cellular traffic",
  "LTE attach and bearer success > 99%",
];

export default function EvidenceCard() {
  return (
    <CardUI sx={{ height: "100%" }}>
      <Typography variant="overline" sx={{ fontWeight: 700, display: "block", mb: 1 }}>
        EVIDENCE
      </Typography>
      <Box component="ul" sx={{ m: 0, pl: 2.5, "& li": { py: 0.25 } }}>
        {EVIDENCE_ITEMS.map((item) => (
          <Box component="li" key={item}>
            <Typography component="span" variant="body2">
              {item}
            </Typography>
          </Box>
        ))}
      </Box>
    </CardUI>
  );
}
