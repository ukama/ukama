/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Phase-1 smoke screen — verifies theme (light/dark/accent/density), fonts,
 * tokens and the data layer end-to-end. Replaced by BizHome in Phase 4.
 */
import { useSyncExternalStore } from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';
import ToggleButton from '@mui/material/ToggleButton';
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup';
import Typography from '@mui/material/Typography';
import { useColorScheme } from '@mui/material/styles';
import { NETWORKS, STATUS_MAP } from '@/data';
import type { StatusTone } from '@/data';
import { useUiPrefs } from '@/lib/store';
import type { Accent, Density } from '@/theme/tokens';

const TONE_COLOR: Record<StatusTone, string> = {
  ok: 'var(--uk-success-bright)',
  warn: 'var(--uk-warning)',
  err: 'var(--uk-error)',
  neutral: 'var(--uk-ink-3)',
};

const emptySubscribe = () => () => {};

/** Hydration-safe mounted check (false on server/hydration, true after). */
function useMounted() {
  return useSyncExternalStore(
    emptySubscribe,
    () => true,
    () => false,
  );
}

function ModeToggle() {
  const { mode, setMode } = useColorScheme();
  const mounted = useMounted();
  if (!mounted) return null;

  return (
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
  );
}

export default function BusinessHomePage() {
  const { accent, density, setAccent, setDensity } = useUiPrefs();

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'background.default', p: 4 }}>
      <Stack spacing={3} sx={{ maxWidth: 900, mx: 'auto' }}>
        <Box>
          <Typography variant="h3">Ukama Console</Typography>
          <Typography variant="body1" color="text.secondary">
            Phase 1 scaffold — theme, fonts, tokens and data layer smoke test.
          </Typography>
        </Box>

        <Card>
          <CardContent>
            <Stack spacing={2}>
              <Typography variant="h6">Appearance</Typography>
              <Stack direction="row" spacing={3} flexWrap="wrap" useFlexGap>
                <ModeToggle />
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
              </Stack>
              <Stack direction="row" spacing={1.5}>
                <Button variant="contained">Primary action</Button>
                <Button variant="outlined">Secondary</Button>
                <Button variant="text">Tertiary</Button>
              </Stack>
            </Stack>
          </CardContent>
        </Card>

        <Card>
          <CardContent>
            <Stack spacing={2}>
              <Typography variant="h6">
                Networks{' '}
                <Typography component="span" color="text.disabled">
                  {NETWORKS.length}
                </Typography>
              </Typography>
              <Stack spacing={1}>
                {NETWORKS.map((n) => {
                  const meta = STATUS_MAP[n.status];
                  const tone: StatusTone = meta ? meta.tone : 'neutral';
                  return (
                    <Stack
                      key={n.id}
                      direction="row"
                      alignItems="center"
                      spacing={2}
                      sx={{
                        minHeight: 'var(--uk-row-h)',
                        px: 'calc(var(--uk-gap) * 0.75)',
                        borderRadius: 'var(--uk-r-sm)',
                        '&:hover': { bgcolor: 'var(--uk-hover)' },
                      }}
                    >
                      <Box
                        sx={{
                          width: 9,
                          height: 9,
                          borderRadius: '50%',
                          flex: 'none',
                          bgcolor: TONE_COLOR[tone],
                        }}
                      />
                      <Typography sx={{ fontWeight: 600, flex: 1 }}>
                        {n.name}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {n.region}
                      </Typography>
                      <Chip
                        size="small"
                        label={meta ? meta.label : n.status}
                        sx={{
                          color: TONE_COLOR[tone],
                          bgcolor: 'var(--uk-line-soft)',
                        }}
                      />
                    </Stack>
                  );
                })}
              </Stack>
            </Stack>
          </CardContent>
        </Card>

        <Typography
          variant="body2"
          color="text.disabled"
          className="tnum"
          sx={{ textAlign: 'center' }}
        >
          tokens: gap var(--uk-gap) · row var(--uk-row-h) · 0123456789
        </Typography>
      </Stack>
    </Box>
  );
}
