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

export default function BackhaulCard() {
  return (
    <CardUI>
      <StatTitle label="BACKHAUL" icon="🌐" />
      <Stack direction="column" spacing={0.5}>
        <StatItem label="Status" value="🟡 Degraded (Capacity Limited)" />
        <StatItem label="Latency" value="640 ms / 820 ms ↑" />
        <StatItem label="Utilization" value="19.8 / 4.1 Mbps" />
        <StatItem label="Estimated Capacity" value="~20.4 Mbps" />
        <StatItem label="Impact" value="Backhaul caps speed" />
        <StatItem label="Confidence" value="High" />
      </Stack>
    </CardUI>
  );
}
