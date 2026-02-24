/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { HorizontalContainerJustify } from "@/styles/global";
import { Stack, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatItem from "../../_components/StatItem";

export interface OverallUnitStatusCardProps {
  status?: string;
  description?: string;
  additionalInfo?: string;
  serviceImpact?: string;
  serviceImpactDetail?: string;
  severity?: string;
  confidence?: string;
  lastEvaluated?: string;
}

export default function OverallUnitStatusCard({
  status = "🟡 BACKHAUL LIMITED",
  description = "Throughput is capped by the backhaul (~20 Mbps).",
  additionalInfo = "LTE logic, RF transmission, unit health, and software are operating normally.",
  serviceImpact = "Service impact (customer-facing):",
  serviceImpactDetail = "Users may experience reduced download speeds even with good LTE signal.",
  severity = "Degraded",
  confidence = "High (0.92)",
  lastEvaluated = "2 min ago",
}: Readonly<OverallUnitStatusCardProps>) {
  return (
    <CardUI>
      <Stack direction="column" spacing={1}>
      <HorizontalContainerJustify>
        <Typography
          variant="overline"
          fontWeight={700}
        >
          OVERALL UNIT STATUS
        </Typography>
        <StatItem label="Last evaluated" value={lastEvaluated} labelVariant="caption" valueVariant="caption" />
      </HorizontalContainerJustify>
      
      <Stack direction="column" spacing={1}>
        <Typography
          variant="subtitle1"
          sx={{
            fontWeight: 700,
            letterSpacing: 0.5,
          }}
        >
          {status}
        </Typography>
        <Typography variant="body2">
          {description}&nbsp;{additionalInfo}
        </Typography>
      </Stack>
      <StatItem label="Severity" value={severity} labelVariant="body2" valueVariant="body2" />
      <StatItem label="Confidence" value={confidence} labelVariant="body2" valueVariant="body2" />
      <StatItem label={serviceImpact} value={serviceImpactDetail} labelVariant="body2" valueVariant="body2" />
      </Stack>
    </CardUI>
  );
}
