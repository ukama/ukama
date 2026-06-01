/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { SkeletonRoundedCard } from '@/styles/global';
import React from 'react';

interface ILoadingWrapper {
  cstyle?: React.CSSProperties;
  width?: string | number;
  height?: string | number;
  children: React.ReactNode;
  radius?: 'small' | 'medium' | 'none';
  isLoading: boolean | undefined;
  variant?: 'text' | 'rectangular' | 'circular';
}

const LoadingWrapper = ({
  children,
  width = '100%',
  height = '100%',
  radius = 'medium',
  variant = 'rectangular',
  isLoading = false,
  cstyle = {},
}: ILoadingWrapper) => {
  const smr = radius === 'small' ? '4px' : '0px';
  const borderRadius = radius === 'medium' ? '10px' : smr;
  if (isLoading)
    return (
      <SkeletonRoundedCard
        width={width}
        height={height}
        variant={variant}
        animation="wave"
        sx={{ borderRadius: borderRadius }}
        style={{ ...cstyle }}
      />
    );

  return (
    <div
      style={{
        height: height ?? 'inherit',
        width: width ?? 'inherit',
        ...cstyle,
      }}
    >
      {children}
    </div>
  );
};

export default LoadingWrapper;
