/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Stack } from "@mui/material";
import CoverageIndicatorsCard from "./CoverageIndicatorsCard";
import RFEventsCard from "./RFEventsCard";
import SignalQualityCard from "./SignalQualityCard";
import TransmissionPowerCard from "./TransmissionPowerCard";
import UnitRFStatusCard from "./UnitRFStatusCard";

export default function RFTab() {
  return (
    <Stack
      spacing={2}
      sx={{
        width: "100%",
        overflow: "auto",
        height: "calc(100vh - 272px)",
      }}
    >
      <UnitRFStatusCard />
      <TransmissionPowerCard />
      <SignalQualityCard />
      <CoverageIndicatorsCard />
      <RFEventsCard />
    </Stack>
  );
}
