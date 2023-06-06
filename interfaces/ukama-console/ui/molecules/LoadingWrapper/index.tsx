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
  const borderRadius =
    radius === 'medium' ? '10px' : radius === 'small' ? '4px' : '0px';
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
        height: height ? height : 'inherit',
        width: width ? width : 'inherit',
        ...cstyle,
      }}
    >
      {children}
    </div>
  );
};

export default LoadingWrapper;
