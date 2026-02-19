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

export default function TowerSubsystemCard() {
  return (
    <CardUI>
      <StatTitle label="TOWER SUBSYSTEM (TRX + COM)" icon="📡" />
      <Stack direction="column" spacing={0.5}>
        <StatItem label="CPU" value="75% → stable" impact="No immediate impact" />
        <StatItem label="Memory" value="90% ↑ slowly" impact="Risk of control-plane instability if sustained" />
        <StatItem label="Storage" value="42% → stable" impact="No immediate impact" />
        <StatItem label="Power" value="23 W → stable" impact="No immediate impact" />
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
