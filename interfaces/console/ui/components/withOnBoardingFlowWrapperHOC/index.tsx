/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Box } from '@mui/material';
import { FunctionComponent } from 'react';
import { RootContainer, GradiantBar, ComponentContainer } from './style';

const withOnBoardingFlowWrapperHOC = (
  WrappedComponent: FunctionComponent<any>,
) => {
  return function HOC(props: any) {
    return (
      <RootContainer maxWidth="sm" disableGutters>
        <GradiantBar />

        <Box sx={ComponentContainer}>
          <WrappedComponent {...props} />
        </Box>
      </RootContainer>
    );
  };
};

export default withOnBoardingFlowWrapperHOC;
