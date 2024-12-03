/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { Role_Type } from '@/client/graphql/generated';
import Welcome from '@/components/Welcome';
import { useAppContext } from '@/context';
import { CenterContainer } from '@/styles/global';
import { roleEnumToString } from '@/utils';
import GradientWrapper from '@/wrappers/gradiantWrapper';
import { useRouter } from 'next/navigation';

const Page = () => {
  const { env, user } = useAppContext();
  const router = useRouter();
  return (
    <CenterContainer>
      <GradientWrapper>
        <Welcome
          orgName={user.orgName}
          handleNext={() => router.push('/configure')}
          operatingCountry={'Dominican Republic of Congo'}
          handleBack={() => router.push(`${env.AUTH_APP_URL}/user/logout`)}
          role={`${user.name} - ${roleEnumToString(user.role as Role_Type)}`}
        />
      </GradientWrapper>
    </CenterContainer>
  );
};

export default Page;
