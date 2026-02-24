/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { MetricItem } from "@/types";
import { Grid2 } from "@mui/material";
import { useState } from "react";
import KeyMetricsCard from "../../_components/KeyMetricsCard";
import OverallUnitStatusCard from "./OverallUnitStatusCard";
import OverviewStatusGrid from "./OverviewStatusGrid";

const HEALTH_METRICS: MetricItem[] = [
  { label: "Health", value: "health" },
  { label: "Network", value: "network" },
  { label: "RF", value: "rf" },
  { label: "Software", value: "software" },
];

/** Viewport offset for scroll containers (header + chrome). */
const SCROLL_OFFSET = {
  /** Mobile: larger offset for header + mobile nav + padding */
  xs: 344,
  /** Desktop: header + sidebar + padding */
  sm: 280,
} as const;

export default function OverviewTab() {
  const [selectedMetric, setSelectedMetric] = useState<string>("health");

  const handleSelectMetric = (metric: string) => {
    setSelectedMetric(metric);
  };

  return (
    <Grid2
      mt={2}
      container
      rowSpacing={2}
      columnSpacing={2}
      sx={{
        height: "100%",
        overflowX: "hidden",
        overflowY: "auto",
        maxHeight: `calc(100vh - ${SCROLL_OFFSET.xs}px)`,
      }}
    >
      {/* <Grid2 size={{ xs: 12, md: 2 }}>
        <GlowieBg>
          <Image
            src="/images/tnode_anode.png"
            alt="Unit Deployment"
            width={220}
            height={400}
            style={{ objectFit: "contain" }}
          />
        </GlowieBg>
      </Grid2> */}
      
      <Grid2 size={12}>
        <OverallUnitStatusCard />
      </Grid2>
      <Grid2 size={{ xs: 12, md: 3 }}>  
        <KeyMetricsCard 
          selectedMetric={selectedMetric} 
          handleSelectMetric={handleSelectMetric} 
          metrics={HEALTH_METRICS}
        />
      </Grid2>
      <Grid2 size={{ xs: 12, md: 9 }}>  
        <OverviewStatusGrid selectedMetric={selectedMetric} />
      </Grid2>
    </Grid2>
  );
}
