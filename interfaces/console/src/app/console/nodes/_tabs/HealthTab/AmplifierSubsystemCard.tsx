/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Stack, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatItem from "../../_components/StatItem";

interface AmplifierSubsystemCardProps {
  selectedMetric: string;
}

export default function AmplifierSubsystemCard({ selectedMetric }: Readonly<AmplifierSubsystemCardProps>) {
  return (
    <CardUI>
      <Stack direction="column" spacing={1}>
        <Typography
          variant="overline"
          fontWeight={700}
        >
          AMPLIFIER SUBSYSTEM
        </Typography> 
        <Stack direction="column" spacing={0.5}>
          {selectedMetric === "compute" && <StatItem label="CPU" value="75% → stable" impact="No immediate impact" />}
          {selectedMetric === "memory" && <StatItem label="Memory" value="90% ↑ slowly" impact="Risk of control-plane instability if sustained" />}
          {selectedMetric === "storage" && <StatItem label="Storage" value="42% → stable" impact="No immediate impact" />}
          {selectedMetric === "thermal" && <StatItem label="Thermal" value="42% → stable" impact="No immediate impact" />}
          {selectedMetric === "power" && <StatItem label="Power" value="23 W → stable" impact="No immediate impact" />}
        </Stack>
      </Stack>
    </CardUI>
  );
}
