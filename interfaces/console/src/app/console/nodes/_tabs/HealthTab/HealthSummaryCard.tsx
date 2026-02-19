/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatItem from "../../_components/StatItem";
import StatTitle from "../../_components/StatTitle";

export default function HealthSummaryCard() {
  return (
    <CardUI>
      <StatTitle label="HEALTH SUMMARY" icon="📊" />
      <Box sx={{ display: "flex", flexWrap: "wrap", gap: 2 }}>
        {["Compute", "Storage", "Thermal", "Power"].map((label) => (
          <StatItem key={label} label={label} value={`🟢 Normal`} />
        ))}
      </Box>
    </CardUI>
  );
}
