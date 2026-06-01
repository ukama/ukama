/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { useEnvContext } from '@/context';
import { CenterContainer } from '@/styles/global';
import { CircularProgress } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

const Page = () => {
  const { env } = useEnvContext();
  const router = useRouter();

  useEffect(() => {
    setTimeout(() => {
      router.push(`${env.AUTH_APP_URL}/user/logout`);
    }, 1000);
  }, []);

  return (
    <CenterContainer>
      <CircularProgress />
    </CenterContainer>
  );
};

export default Page;
