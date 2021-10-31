import { SVGType } from "../../types";

const BatteryIcon = ({
    color = "#03744B",
    width = "20px",
    height = "10px",
}: SVGType) => (
    <svg
        width={width}
        height={height}
        viewBox="0 0 20 10"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
    >
        <path
            d="M18 6.25L17.25 6.25L17.25 7L17.25 8.67C17.25 8.98579 16.9858 9.25 16.67 9.25L1.33 9.25C1.01786 9.25 0.75 8.98944 0.75 8.66L0.75 1.33C0.75 1.01786 1.01056 0.749999 1.34 0.749999L16.67 0.75C16.9858 0.75 17.25 1.01421 17.25 1.33L17.25 3L17.25 3.75L18 3.75L19.25 3.75L19.25 6.25L18 6.25Z"
            stroke={color}
            strokeWidth="1.5"
        />
    </svg>
);

export default BatteryIcon;
