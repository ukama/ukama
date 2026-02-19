/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import CardUI from "../../_components/CardUI";
import StatItem from "../../_components/StatItem";
import StatTitle from "../../_components/StatTitle";

export default function PowerUnitCard() {
  return (
    <CardUI>
      <StatTitle label="POWER (UNIT)" icon="⚡" />
      <StatItem label="Unit Power" value="71 W → stable" impact="No immediate impact" />
      <StatItem label="Peak (24h)" value="84 W" />
    </CardUI>
  );
}
