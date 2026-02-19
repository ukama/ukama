/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";

export interface RecentEventsCardProps {
  title?: string;
  events?: string[];
}

const DEFAULT_EVENTS = [
  "Backhaul capacity measured at 20.4 Mbps (1h ago)",
  "Backhaul RTT p95 elevated (last 30m)",
  "No critical service restarts detected in last 24h",
  "No thermal or power events detected",
];

export default function RecentEventsCard({
  title = "RECENT EVENTS",
  events = DEFAULT_EVENTS,
}: RecentEventsCardProps) {
  return (
    <CardUI>
      <Typography
        variant="overline"
        sx={{
          display: "flex",
          fontWeight: 700,
          alignItems: "center",
        }}
      >
        {title}
      </Typography>
      <Box
        component="ul"
        sx={{
          pl: 2,
          listStyle: "disc",
          m: 0,
          "& li": {
            lineHeight: 1.6,
            fontSize: "0.875rem",
            "&:last-child": { mb: 0 },
          },
        }}
      >
        {events.map((event) => (
          <li key={event}>{event}</li>
        ))}
      </Box>
    </CardUI>
  );
}
