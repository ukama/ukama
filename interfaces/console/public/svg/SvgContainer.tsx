import React, { ReactNode } from 'react';
type SvgContainerProps = {
  children: ReactNode;
};
export const SvgContainer = ({ children }: SvgContainerProps) => (
  <svg
    width="80"
    height="80"
    viewBox="0 0 80 80"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect x="0" y="0" width="80" height="80" rx="20" fill="#F2F8F6" />
    <g transform="translate(30, 25)">{children}</g>
  </svg>
);
