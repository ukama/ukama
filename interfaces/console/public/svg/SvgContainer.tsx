import React, { ReactNode } from 'react';

type SvgContainerProps = {
  children: ReactNode;
  borderColor?: string;
  translateX?: number;
  translateY?: number;
  id?: string; // Optional id prop
  className?: string; // Optional class prop
};

export const SvgContainer = ({
  children,
  borderColor,
  translateX = 30,
  translateY = 25,
  id,
  className,
}: SvgContainerProps) => (
  <svg
    width="80"
    height="80"
    viewBox="0 0 80 80"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    id={id} // Pass the id prop to the svg element
    className={className} // Pass the class prop to the svg element
  >
    <rect
      x="0"
      y="0"
      width="80"
      height="80"
      rx="20"
      fill="#F2F8F6"
      stroke={borderColor}
      strokeWidth={borderColor ? 4 : 0}
    />
    <g transform={`translate(${translateX}, ${translateY})`}>{children}</g>
  </svg>
);
