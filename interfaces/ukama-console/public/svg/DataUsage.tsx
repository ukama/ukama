import { SVGType } from '@/types';

export const DataUsage = ({
  color = '#6974F8',
  width = '48',
  height = '48',
}: SVGType) => (
  <svg
    width={width}
    height={height}
    viewBox="0 0 48 48"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect width="48" height="48" rx="4" fill={color} fillOpacity="0.1" />
    <path
      d="M25 14.05V17.08C28.39 17.57 31 20.47 31 24C31 24.9 30.82 25.75 30.52 26.54L33.12 28.07C33.68 26.83 34 25.45 34 24C34 18.82 30.05 14.55 25 14.05ZM24 31C20.13 31 17 27.87 17 24C17 20.47 19.61 17.57 23 17.08V14.05C17.94 14.55 14 18.81 14 24C14 29.52 18.47 34 23.99 34C27.3 34 30.23 32.39 32.05 29.91L29.45 28.38C28.17 29.98 26.21 31 24 31Z"
      fill={color}
    />
  </svg>
);
