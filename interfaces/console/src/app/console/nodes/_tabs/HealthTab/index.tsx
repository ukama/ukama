/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { MetricItem } from "@/types";
import { Grid2, Stack } from "@mui/material";
import { useState } from "react";
import KeyMetricsCard from "../../_components/KeyMetricsCard";
import AmplifierSubsystemCard from "./AmplifierSubsystemCard";
import TowerSubsystemCard from "./TowerSubsystemCard";
import UnitHealthCard from "./UnitHealthCard";

const HEALTH_METRICS: MetricItem[] = [
  { label: "Compute", value: "compute", icon: "🟢" },
  { label: "Memory", value: "memory", icon: "🟡" },
  { label: "Storage", value: "storage", icon: "🟢" },
  { label: "Thermal", value: "thermal", icon: "🔴" },
  { label: "Power", value: "power", icon: "🟢" },
];

export default function HealthTab() {
  const [selectedMetric, setSelectedMetric] = useState<string>("compute");

  const handleSelectMetric = (metric: string) => {
    setSelectedMetric(metric);
  };
  
  return (
    <Stack spacing={2}>
      <Grid2 container spacing={2}>
        {/* <Grid2 size={{ xs: 12  }}>
          <HealthSummaryCard />
        </Grid2> */}
        <Grid2 size={12}>
          <UnitHealthCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 3 }}>  
          <KeyMetricsCard 
            selectedMetric={selectedMetric} 
            handleSelectMetric={handleSelectMetric} 
            metrics={HEALTH_METRICS}/>
        </Grid2>
        <Grid2 size={{ xs: 12, md: 9 }} container spacing={2}>
          <Grid2 size={12}>
            <TowerSubsystemCard selectedMetric={selectedMetric} />
          </Grid2>
          <Grid2 size={12}>
            <AmplifierSubsystemCard selectedMetric={selectedMetric} />
          </Grid2>
        </Grid2>
      </Grid2>
    </Stack>
  );
}
