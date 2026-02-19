/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Stack } from "@mui/material";
import CardUI from "../../_components/CardUI";
import StatItem from "../../_components/StatItem";
import StatTitle from "../../_components/StatTitle";

export default function CoverageIndicatorsCard() {
  return (
    <CardUI sx={{ fontFamily: "monospace" }}>
      <StatTitle label="COVERAGE INDICATORS" icon="🗺️" />
      <Stack direction="column" spacing={0.5}>
        <StatItem
          label="UE Distance"
          value="Near 40% | Mid 45% | Far 15%"
        />
        <StatItem label="Coverage vs baseline" value="Stable" />
        <StatItem label="Service impact" value="Coverage unchanged" />
      </Stack>
    </CardUI>
  );
}
