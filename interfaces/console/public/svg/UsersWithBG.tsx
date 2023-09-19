import { SVGType } from '@/types';

export const UsersWithBG = ({
  color = '#2190F6',
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
      d="M28 23C29.66 23 30.99 21.66 30.99 20C30.99 18.34 29.66 17 28 17C26.34 17 25 18.34 25 20C25 21.66 26.34 23 28 23ZM20 23C21.66 23 22.99 21.66 22.99 20C22.99 18.34 21.66 17 20 17C18.34 17 17 18.34 17 20C17 21.66 18.34 23 20 23ZM20 25C17.67 25 13 26.17 13 28.5V31H27V28.5C27 26.17 22.33 25 20 25ZM28 25C27.71 25 27.38 25.02 27.03 25.05C28.19 25.89 29 27.02 29 28.5V31H35V28.5C35 26.17 30.33 25 28 25Z"
      fill={color}
    />
  </svg>
);
