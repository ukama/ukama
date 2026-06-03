/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Placeholder screen used by the Phase-2 routing skeleton. */
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import Typography from '@mui/material/Typography';
import DashboardCustomizeRounded from '@mui/icons-material/DashboardCustomizeRounded';
import PageHeader from './PageHeader';

export default function PageStub({
  title,
  crumb,
  count,
  sub,
  phase,
}: {
  title: string;
  crumb?: string[];
  count?: string | number;
  sub?: string;
  phase: string;
}) {
  return (
    <div className="page">
      <PageHeader crumb={crumb} title={title} count={count} sub={sub} />
      <Card sx={{ p: '56px 24px', textAlign: 'center' }}>
        <DashboardCustomizeRounded
          sx={{ fontSize: 42, color: 'text.disabled' }}
        />
        <Typography
          sx={{
            fontFamily: 'var(--font-display)',
            fontSize: 18,
            fontWeight: 500,
            mt: 1.5,
          }}
        >
          {title}
        </Typography>
        <Box sx={{ fontSize: 13.5, mt: 0.75, color: 'text.secondary' }}>
          This screen is being built in Phase {phase}.
        </Box>
      </Card>
    </div>
  );
}
