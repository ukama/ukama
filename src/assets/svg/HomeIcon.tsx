import { SVGType } from "../../types";

const HomeIcon = ({
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
            d="M11.6668 23.3333V16.3333H16.3335V23.3333H22.1668V14H25.6668L14.0002 3.5L2.3335 14H5.8335V23.3333H11.6668Z"
            fill={color}
        />
    </svg>
);

export default HomeIcon;
