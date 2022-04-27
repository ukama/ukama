import { SVGType } from "../../types";

export const DataBilling = ({ width = "48", height = "48" }: SVGType) => (
    <svg
        height={height}
        viewBox="0 0 48 48"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
    >
        <rect
            width={width}
            height={height}
            rx="4"
            fill="#00D3EB"
            fillOpacity="0.1"
        />
        <path
            d="M32 16H16C14.89 16 14.01 16.89 14.01 18L14 30C14 31.11 14.89 32 16 32H32C33.11 32 34 31.11 34 30V18C34 16.89 33.11 16 32 16ZM32 30H16V24H32V30ZM32 20H16V18H32V20Z"
            fill="#00D3EB"
        />
    </svg>
);
