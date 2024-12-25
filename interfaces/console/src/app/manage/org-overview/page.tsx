/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

'use client';
import { useGetOrgTreeQuery } from '@/client/graphql/generated';
import { OrgTree } from '@/components/OrgTree';
import { CenterContainer } from '@/styles/global';
import { CircularProgress } from '@mui/material';

const Page = () => {
  const { data, loading } = useGetOrgTreeQuery();

  if (loading) {
    return (
      <CenterContainer>
        <CircularProgress />
      </CenterContainer>
    );
  }

  return <OrgTree data={data?.getOrgTree} />;
};

export default Page;
