/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import GlowieBg from "@/components/GlowieBg";
import { Grid2 } from "@mui/material";
import Image from "next/image";
import KeyMetricsCard from "./KeyMetricsCard";
import OverallUnitStatusCard from "./OverallUnitStatusCard";
import OverviewStatusGrid from "./OverviewStatusGrid";
import RecentEventsCard from "./RecentEventsCard";

/** Viewport offset for scroll containers (header + chrome). */
const SCROLL_OFFSET = {
  /** Mobile: larger offset for header + mobile nav + padding */
  xs: 344,
  /** Desktop: header + sidebar + padding */
  sm: 280,
} as const;

export default function OverviewTab() {
  return (
    <Grid2
      mt={2}
      container
      rowSpacing={{ xs: 2, md: 8 }}
      sx={{
        height: "100%",
        overflowX: "hidden",
        overflowY: { xs: "auto", sm: "hidden" },
        maxHeight: { xs: `calc(100vh - ${SCROLL_OFFSET.xs}px)`, sm: "none" },
      }}
    >
      <Grid2 size={{ xs: 12, md: 2 }}>
        <GlowieBg>
          <Image
            src="/images/tnode_anode.png"
            alt="Unit Deployment"
            width={220}
            height={400}
            style={{ objectFit: "contain" }}
          />
        </GlowieBg>
      </Grid2>

      <Grid2
        container
        size={{ xs: 12, md: 10 }}
        rowSpacing={2}
        columnSpacing={2}
        sx={{
          overflowX: "hidden",
          overflowY: { xs: "hidden", sm: "auto" },
          maxHeight: { xs: "none", sm: `calc(100vh - ${SCROLL_OFFSET.sm}px)` },
          minHeight: 0, // Allow flex child to scroll
        }}
      >
        <Grid2 size={{ xs: 12, md: 6 }}>
          <OverallUnitStatusCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 6 }}>  
          <KeyMetricsCard />
        </Grid2>
        <Grid2 size={{ xs: 12, md: 12 }}>  
          <OverviewStatusGrid />
        </Grid2>
        <Grid2 size={12}>  
          <RecentEventsCard />
        </Grid2>
      </Grid2>
    </Grid2>
  );
}
