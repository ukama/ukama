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

export default function UnitHealthCard() {
  return (
    <CardUI>
      <StatTitle label="UNIT HEALTH" icon="🔧" />
      <StatItem label="Overall" value="🟢 Healthy" />
      <StatItem label="Tower Subsystem" value="🟢 Normal" />
      <StatItem label="Amplifier Subsystem" value="🟢 Normal" />
      <StatItem label="Severity" value="Healthy" />
      <StatItem label="Confidence" value="High (0.94)" />
      <StatItem label="Last evaluated" value="2 min ago" />
      <StatItem label="Service impact (customer-facing)" value="No physical, compute, thermal, storage, or power constraints affecting service." />
    </CardUI>
  );
}
