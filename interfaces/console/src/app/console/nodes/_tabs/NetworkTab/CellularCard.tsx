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

export default function CellularCard() {
  return (
    <CardUI>
      <StatTitle label="CELLULAR (Tower Node)" icon="📡" />
      <Stack direction="column" spacing={0.5}>
        <StatItem label="Status" value="🟢 Healthy" />
        <StatItem label="DL / UL" value="18.9 / 4.2 Mbps → stable" />
        <StatItem label="Active UEs" value="3 → steady" />
        <StatItem label="RRC Success" value="99.6%" />
        <StatItem label="ERAB Success" value="99.8%" />
        <StatItem label="RSRP" value="-89 dBm (Good)" />
        <StatItem label="Impact" value="No cellular limitation" />
      </Stack>
    </CardUI>
  );
}
