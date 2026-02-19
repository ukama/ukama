/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Grid2, Stack } from "@mui/material";
import BackhaulCard from "./BackhaulCard";
import CellularCard from "./CellularCard";
import EvidenceCard from "./EvidenceCard";
import UnitNetworkStatusCard from "./UnitNetworkStatusCard";

export default function NetworkTab() {
  return (
    <Stack spacing={2} sx={{ width: "100%", overflow: "auto", height: "calc(100vh - 272px)" }}>
      <Grid2 container spacing={2}>
        <Grid2 size={{ xs: 12 }}>
          <UnitNetworkStatusCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 6 }}>
          <CellularCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 6 }}>
          <BackhaulCard />
        </Grid2>
        <Grid2 size={{ xs: 12 }}>
          <EvidenceCard />
        </Grid2>
      </Grid2>
    </Stack>
  );
}
