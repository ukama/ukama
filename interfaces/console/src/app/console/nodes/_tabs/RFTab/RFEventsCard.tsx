/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Box, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatTitle from "../../_components/StatTitle";

const RF_EVENTS = [
  "No TX power reduction events",
  "No PA faults detected",
  "No RF thermal throttling",
];

export default function RFEventsCard() {
  return (
    <CardUI sx={{ fontFamily: "monospace" }}>
      <StatTitle label="RF EVENTS" icon="📜" />
      <Box component="ul" sx={{ m: 0, pl: 2.5, "& li": { py: 0.25 } }}>
        {RF_EVENTS.map((item) => (
          <Box component="li" key={item}>
            <Typography component="span" variant="body2">
              {item}
            </Typography>
          </Box>
        ))}
      </Box>
    </CardUI>
  );
}