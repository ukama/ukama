/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import AppSnackbar from '@/components/AppSnackbar/page';
import { CenterContainer } from '@/styles/global';
import GradientWrapper from '@/wrappers/gradiantWrapper';

const ConfigureLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  return (
    <CenterContainer>
      <GradientWrapper>{children}</GradientWrapper>
      <AppSnackbar />
    </CenterContainer>
  );
};

export default ConfigureLayout;
