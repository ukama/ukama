/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";

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
}: OverallUnitStatusCardProps) {
  return (
    <CardUI sx={{ height: {xs: "auto", sm: "220px"} }}>
      <Typography
        variant="overline"
        sx={{
          display: "flex",
          fontWeight: 700,
          alignItems: "center",
        }}
      >
        OVERALL UNIT STATUS
      </Typography>
      <Typography
        variant="subtitle1"
        sx={{
          fontWeight: 700,
          letterSpacing: 0.5,
          display: "flex",
          alignItems: "center",
          mb: 1,
        }}
      >
        {status}
      </Typography>
      <Typography variant="body2" sx={{ mb: 1, lineHeight: 1.6 }}>
        {description}
      </Typography>
      <Typography variant="body2" sx={{ mb: 2, lineHeight: 1.6 }}>
        {additionalInfo}
      </Typography>
      <Typography variant="body2" sx={{ fontWeight: 500, mb: 0.5 }}>
        {serviceImpact}
      </Typography>
      <Typography variant="body2" sx={{ mb: 2, lineHeight: 1.6 }}>
        {serviceImpactDetail}
      </Typography>
      <Box
        sx={{
          display: "flex",
          flexWrap: "wrap",
          gap: 2,
          fontSize: "0.875rem",
        }}
      >
        <span>Severity: {severity}</span>
        <span>Confidence: {confidence}</span>
        <span>Last evaluated: {lastEvaluated}</span>
      </Box>
    </CardUI>
  );
}
