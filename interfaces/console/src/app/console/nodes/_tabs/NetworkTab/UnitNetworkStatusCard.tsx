/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatTitle from "../../_components/StatTitle";

export default function UnitNetworkStatusCard() {
  return (
    <CardUI>
      <StatTitle label="UNIT NETWORK STATUS" icon="🧠" />
      <Typography variant="subtitle1" sx={{ fontWeight: 700, letterSpacing: 0.5, mb: 1 }}>
        🟡 BACKHAUL LIMITED
      </Typography>
      <Typography variant="body2" sx={{ mb: 1, lineHeight: 1.6 }}>
        Throughput capped at ~20 Mbps. LTE operating normally.
      </Typography>
      <Typography variant="body2" sx={{ mb: 1 }}>
        Root cause: Backhaul capacity limit
      </Typography>
      <Typography variant="body2" sx={{ fontWeight: 500, mb: 0.5 }}>
        Service impact (customer-facing):
      </Typography>
      <Typography variant="body2" sx={{ mb: 1.5, lineHeight: 1.6 }}>
        Users experience reduced download speeds; latency may increase under load.
      </Typography>
      <Box sx={{ display: "flex", flexWrap: "wrap", justifyContent: "space-between", gap: 1 }}>
        <Typography variant="body2">Severity: Degraded</Typography>
        <Typography variant="body2">Confidence: High (0.92)</Typography>
      </Box>
      <Typography variant="body2" sx={{ mt: 0.5 }}>
        Last evaluated: 3 min ago
      </Typography>
    </CardUI>
  );
}
