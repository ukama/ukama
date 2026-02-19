/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Grid2, Stack } from "@mui/material";
import AmplifierSubsystemCard from "./AmplifierSubsystemCard";
import HealthSummaryCard from "./HealthSummaryCard";
import PowerUnitCard from "./PowerUnitCard";
import TowerSubsystemCard from "./TowerSubsystemCard";
import UnitHealthCard from "./UnitHealthCard";

export default function HealthTab() {
  return (
    <Stack spacing={2} sx={{ width: "100%", overflow: "auto", height: "calc(100vh - 272px)" }}>
      <Grid2 container spacing={2}>
        <Grid2 size={{ xs: 12  }}>
          <HealthSummaryCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 6 }}>
          <UnitHealthCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 6 }}>
          <TowerSubsystemCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 6 }}>
          <AmplifierSubsystemCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 6 }}>
          <PowerUnitCard />
        </Grid2>
      </Grid2>
    </Stack>
  );
}
