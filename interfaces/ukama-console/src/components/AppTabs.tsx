/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Tab strip — MUI Tabs themed to the design (.tabs/.tab). */
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';

export default function AppTabs({
  tabs,
  value,
  onChange,
  scrollable,
}: {
  tabs: readonly string[];
  value: string;
  onChange: (value: string) => void;
  scrollable?: boolean;
}) {
  return (
    <Tabs
      value={value}
      onChange={(_, v: string) => onChange(v)}
      {...(scrollable ? { variant: 'scrollable' as const, scrollButtons: false } : {})}
    >
      {tabs.map((t) => (
        <Tab key={t} value={t} label={t} />
      ))}
    </Tabs>
  );
}
