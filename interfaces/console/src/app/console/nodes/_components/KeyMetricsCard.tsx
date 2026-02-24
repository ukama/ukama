/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { HorizontalContainerJustify } from "@/styles/global";
import { MetricItem } from "@/types";
import { Button, Stack, Typography } from "@mui/material";
import CardUI from "./CardUI";

export interface KeyMetricsCardProps {
  title?: string;
  metrics?: MetricItem[];
  selectedMetric: string;
  handleSelectMetric: (metric: string) => void;
}

export default function KeyMetricsCard({
  metrics = [],
  selectedMetric,
  handleSelectMetric,
  title = "KEY METRICS",
}: Readonly<KeyMetricsCardProps>) {
  return (
    <CardUI isBorderLeft={true} sx={{ maxHeight: "fit-content" }}>
      <Stack direction="column" spacing={1}>  
        <HorizontalContainerJustify>  
          <Typography
            variant="overline"
            fontWeight={700}
          >
            {title}
          </Typography>
        </HorizontalContainerJustify>

        <Stack
          spacing={1}
          direction="column"
          alignItems="flex-start"
        >
          {metrics.map((metric) => (
            <Button 
              size="small" 
              variant="text" 
              key={metric.label} 
              onClick={() => handleSelectMetric(metric.value)} 
              sx={{ 
                typography: "body2",
                textTransform: "capitalize", 
                fontWeight: selectedMetric === metric.value ? 500 : 400,
                color: selectedMetric === metric.value ? "primary.main" : "text.primary", 
              }}
            >
              {metric.icon}&nbsp;&nbsp;{metric.label}
            </Button>
          ))}
        </Stack>
      </Stack>
    </CardUI>
  );
}
