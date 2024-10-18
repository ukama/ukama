/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { Role_Type } from '@/client/graphql/generated/subscriptions';
import { useAppContext } from '@/context';
import { Skeleton } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

const Page = () => {
  const router = useRouter();
  const { user } = useAppContext();

  useEffect(() => {
    if (user.role === Role_Type.RoleOwner) {
      router.push('/manage/members');
    }
  }, [user.role]);

  return (
    <Skeleton
      variant="rectangular"
      sx={{
        width: '100%',
        borderRadius: '10px',
        height: 'calc(100vh - 400px)',
      }}
    />
  );
};

export default Page;
