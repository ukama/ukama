/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { useAppContext } from '@/context';
import '@/styles/console.css';
import { CenterContainer } from '@/styles/global';
import { Stack, Typography } from '@mui/material';
import Link from 'next/link';

const Page = () => {
  const { env } = useAppContext();
  return (
    <CenterContainer>
      <Stack spacing={0.5} alignItems={'center'}>
        <Typography variant="body1">
          {"Sorry, You don't have permissions to view this page"}
        </Typography>
        <Link href="/logout">Log me out</Link>
      </Stack>
    </CenterContainer>
  );
};

export default Page;
