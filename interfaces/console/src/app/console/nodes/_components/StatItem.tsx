/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Stack, Typography } from "@mui/material";

interface StatItemProps {
  label: string;
  value: string;
  impact?: string;
}

export default function StatItem({ label, value, impact }: StatItemProps) {
  return (
    <Stack direction="column" spacing={0.2}>
        <Stack direction="row" spacing={2} mb={0.8}>
            <Typography variant="body2" sx={{ fontWeight: 500 }}>{label}:</Typography>
            <Typography variant="body2">{value}</Typography>
        </Stack>
        {impact && (
            <Stack direction="row" spacing={2} mb={0.8}>
                <Typography variant="body2" sx={{ fontWeight: 500 }}>Service impact:</Typography>
                <Typography variant="body2">{value}</Typography>
            </Stack>
        )}
    </Stack>);
}