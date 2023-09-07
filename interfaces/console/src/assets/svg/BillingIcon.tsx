import { SVGType } from "../../types";
export const BillingIcon = ({
    color = "black",
    width = "24",
    height = "24",
}: SVGType) => (
    <svg
        width={width}
        height={height}
        viewBox="0 0 28 28"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
    >
        <path
            d="M23.3335 4.66675H4.66683C3.37183 4.66675 2.34516 5.70508 2.34516 7.00008L2.3335 21.0001C2.3335 22.2951 3.37183 23.3334 4.66683 23.3334H23.3335C24.6285 23.3334 25.6668 22.2951 25.6668 21.0001V7.00008C25.6668 5.70508 24.6285 4.66675 23.3335 4.66675ZM23.3335 21.0001H4.66683V14.0001H23.3335V21.0001ZM23.3335 9.33341H4.66683V7.00008H23.3335V9.33341Z"
            fill={color}
        />
    </svg>
);
