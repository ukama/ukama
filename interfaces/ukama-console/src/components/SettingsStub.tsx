/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Settings screen placeholder: live Appearance card now, full tabs in Phase 7. */
import Card from '@mui/material/Card';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import AppearanceSettings from './AppearanceSettings';
import PageHeader from './PageHeader';

export default function SettingsStub() {
  return (
    <div className="page">
      <PageHeader
        title="Settings"
        sub="Manage your account, organization and billing."
      />
      <Stack spacing={2} sx={{ maxWidth: 720 }}>
        <AppearanceSettings />
        <Card sx={{ p: '32px 24px', textAlign: 'center' }}>
          <Typography sx={{ fontSize: 13.5, color: 'text.secondary' }}>
            Account, organization, notification and billing settings are being
            built in Phase 7.
          </Typography>
        </Card>
      </Stack>
    </div>
  );
}
