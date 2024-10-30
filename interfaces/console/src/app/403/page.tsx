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
import { colors } from '@/theme';
import { Container, Stack, Typography } from '@mui/material';
import Link from 'next/link';
import { Forbidden } from '../../../public/svg/403';

const Page = () => {
  const { env } = useAppContext();
  return (
    <Container maxWidth="md" sx={{ height: '100vh' }}>
      <CenterContainer>
        <Stack spacing={1} alignItems={'center'}>
          <Forbidden />
          <Typography
            variant="body1"
            fontFamily={'Work Sans, sans-serif'}
            textAlign={'center'}
          >
            Sorry, You don&apos;t have valid permission to access console.
            <br />
            If you think its wrong, Please contact your network owner to review
            your role.
          </Typography>
          <Link
            href={`${env.AUTH_APP_URL}/user/logout`}
            style={{
              fontSize: 16,
              fontWeight: 600,
              color: colors.primaryMain,
              fontFamily: 'Rubik, sans-serif',
            }}
          >
            Logout
          </Link>
        </Stack>
      </CenterContainer>
    </Container>
  );
};

export default Page;
