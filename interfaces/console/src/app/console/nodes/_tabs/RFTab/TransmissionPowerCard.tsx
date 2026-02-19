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

export default function TransmissionPowerCard() {
  return (
    <CardUI sx={{ fontFamily: "monospace" }}>
      <StatTitle label="TRANSMISSION & POWER" icon="📡" />
      <Stack direction="column" spacing={0.5}>
        <StatItem label="TX State" value="ON" />
        <StatItem label="Target TX Power" value="23 dBm" />
        <StatItem label="Actual TX Power" value="22.8 dBm → stable" />
        <StatItem label="PA State" value="Operational" />
        <StatItem label="Service impact" value="No immediate impact" />
      </Stack>
    </CardUI>
  );
}
