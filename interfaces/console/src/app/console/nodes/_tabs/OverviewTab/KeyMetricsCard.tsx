/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box, Stack, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";

export interface MetricItem {
  label: string;
  value: string;
  trend?: string;
}

export interface KeyMetricsCardProps {
  title?: string;
  metrics?: MetricItem[];
}

const DEFAULT_METRICS: MetricItem[] = [
  { label: "Unit Health", value: "🟢 Normal" },
  { label: "Active UEs", value: "3", trend: "steady" },
  { label: "Cellular DL", value: "18.9 Mbps", trend: "stable" },
  { label: "Cellular UL", value: "4.2 Mbps", trend: "stable" },
  {
    label: "Backhaul RTT (p95)",
    value: "820 ms",
    trend: "↑ increasing",
  },
];

export default function KeyMetricsCard({
  title = "KEY METRICS (LAST 15 MIN)",
  metrics = DEFAULT_METRICS,
}: KeyMetricsCardProps) {
  return (
    <CardUI sx={{ height: {xs: "auto", sm: "220px"} }}>
      <Typography
        variant="overline"
        sx={{
          display: "flex",
          fontWeight: 700,
          alignItems: "center",
          mb: 1,
        }}
      >
        {title}
      </Typography>
      <Stack
        spacing={2}
        direction="column"
      >
        {metrics.map((metric) => (
          <Box
            key={metric.label}
            sx={{
              display: "flex",
              alignItems: "baseline",
              flexWrap: "wrap",
            }}
          >
            <Typography variant="body2" sx={{ fontWeight: 500 }}>
              {metric.label}: &nbsp;
            </Typography>
            <Typography variant="body2">
              {metric.value}
              {metric.trend && ` → ${metric.trend}`}
            </Typography>
          </Box>
        ))}
      </Stack>
    </CardUI>
  );
}
