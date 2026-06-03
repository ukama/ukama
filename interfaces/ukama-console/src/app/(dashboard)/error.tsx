/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Dashboard segment error boundary (BUILD-PLAN §13.3). */
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import Typography from '@mui/material/Typography';

export default function DashboardError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="page">
      <Card sx={{ p: '56px 24px', textAlign: 'center', mt: 3 }}>
        <Typography
          sx={{ fontFamily: 'var(--font-display)', fontSize: 18, fontWeight: 500 }}
        >
          We couldn’t load this
        </Typography>
        <Typography sx={{ fontSize: 13.5, mt: 0.75, color: 'text.secondary' }}>
          Something went wrong rendering this screen.
          {error.digest ? ` (ref ${error.digest})` : ''}
        </Typography>
        <Button variant="outlined" sx={{ mt: 2.5 }} onClick={() => reset()}>
          Try again
        </Button>
      </Card>
    </div>
  );
}
