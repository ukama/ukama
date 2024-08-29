/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  ComponentContainer,
  GradiantBar,
  RootContainer,
} from '@/styles/global';
import { Box } from '@mui/material';

import { ReactNode } from 'react';

const GradientWrapper = ({ children }: { children: ReactNode }) => {
  return (
    <RootContainer maxWidth="sm" disableGutters>
      <GradiantBar />
      <Box sx={ComponentContainer}>{children}</Box>
    </RootContainer>
  );
};

export default GradientWrapper;
