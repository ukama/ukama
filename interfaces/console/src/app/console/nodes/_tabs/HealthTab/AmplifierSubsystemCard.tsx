/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Link, Stack } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatItem from "../../_components/StatItem";
import StatTitle from "../../_components/StatTitle";

export default function AmplifierSubsystemCard() {
  return (
    <CardUI>
      <StatTitle label="AMPLIFIER SUBSYSTEM" icon="📡" />
      <Stack direction="column" spacing={0.5}>
        <StatItem label="RF Output" value="Nominal" />
        <StatItem label="Temperature" value="61°C → stable" />
        <StatItem label="Power" value="48 W → stable" impact="No immediate impact" />
      </Stack>
      <Link
        href="#"
        sx={{
          display: "inline-flex",
          alignItems: "center",
          mt: 1,
          fontFamily: "monospace",
          fontSize: "0.75rem",
          textDecoration: "none",
          "&:hover": { textDecoration: "underline" },
        }}
      >
        ► View board-level details
      </Link>
    </CardUI>
  );
}
