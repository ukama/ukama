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
