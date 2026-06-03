/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Appearance controls (mode · accent · density) — BUILD-PLAN §7.2 F. */
import { useSyncExternalStore } from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Stack from '@mui/material/Stack';
import ToggleButton from '@mui/material/ToggleButton';
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import Typography from '@mui/material/Typography';
import { useColorScheme } from '@mui/material/styles';
import { useUiPrefs } from '@/lib/store';
import type { Accent, Density } from '@/theme/tokens';

const emptySubscribe = () => () => {};
const useMounted = () =>
  useSyncExternalStore(
    emptySubscribe,
    () => true,
    () => false,
  );

function Row({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <Stack
      direction="row"
      alignItems="center"
      justifyContent="space-between"
      sx={{ py: 1.5, flexWrap: 'wrap', gap: 1.5 }}
    >
      <Typography sx={{ fontSize: 14, fontWeight: 600 }}>{label}</Typography>
      {children}
    </Stack>
  );
}

export default function AppearanceSettings() {
  const { mode, setMode } = useColorScheme();
  const { accent, density, setAccent, setDensity } = useUiPrefs();
  const mounted = useMounted();
  if (!mounted) return null;

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Appearance
        </Typography>
        <Row label="Theme">
          <ToggleButtonGroup
            size="small"
            exclusive
            value={mode ?? 'light'}
            onChange={(_, v: 'light' | 'dark' | 'system' | null) => {
              if (v) setMode(v);
            }}
            aria-label="Color mode"
          >
            <ToggleButton value="light">Light</ToggleButton>
            <ToggleButton value="dark">Dark</ToggleButton>
            <ToggleButton value="system">System</ToggleButton>
          </ToggleButtonGroup>
        </Row>
        <Row label="Accent">
          <ToggleButtonGroup
            size="small"
            exclusive
            value={accent}
            onChange={(_, v: Accent | null) => {
              if (v) setAccent(v);
            }}
            aria-label="Accent"
          >
            <ToggleButton value="blue">Blue</ToggleButton>
            <ToggleButton value="indigo">Indigo</ToggleButton>
            <ToggleButton value="teal">Teal</ToggleButton>
          </ToggleButtonGroup>
        </Row>
        <Row label="Density">
          <ToggleButtonGroup
            size="small"
            exclusive
            value={density}
            onChange={(_, v: Density | null) => {
              if (v) setDensity(v);
            }}
            aria-label="Density"
          >
            <ToggleButton value="comfortable">Comfortable</ToggleButton>
            <ToggleButton value="compact">Compact</ToggleButton>
          </ToggleButtonGroup>
        </Row>
      </CardContent>
    </Card>
  );
}
