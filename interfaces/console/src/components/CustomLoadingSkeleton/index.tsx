/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Skeleton, SxProps, Theme } from '@mui/material';

interface CustomSkeletonProps {
  width?: number | string;
  height?: number | string;
  variant?: 'text' | 'rectangular' | 'circular';
  sx?: SxProps<Theme>;
}

const CustomSLoadingSkeleton: React.FC<CustomSkeletonProps> = ({
  width = 100,
  height = 20,
  variant = 'rectangular',
  sx,
}) => (
  <Skeleton
    width={width}
    height={height}
    variant={variant}
    sx={{ borderRadius: '5px', ...sx }}
  />
);

export default CustomSLoadingSkeleton;
