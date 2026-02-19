/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Stack } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatItem from "../../_components/StatItem";
import StatTitle from "../../_components/StatTitle";

export default function SignalQualityCard() {
  return (
    <CardUI sx={{ fontFamily: "monospace" }}>
      <StatTitle label="SIGNAL QUALITY (LAST 15 MIN)" icon="📶" />
      <Stack direction="column" spacing={0.5}>
        <StatItem label="Median RSRP" value="-89 dBm (Good) → stable" />
        <StatItem label="Median SINR" value="14 dB (Normal) → stable" />
        <StatItem
          label="UE Distribution"
          value="Good 67% | Fair 25% | Poor 8%"
        />
        <StatItem label="Service impact" value="No coverage degradation" />
      </Stack>
    </CardUI>
  );
}
