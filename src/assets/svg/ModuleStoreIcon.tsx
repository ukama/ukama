import { SVGType } from "../../types";
const ModuleStoreIcon = ({ color }: SVGType, props: any) => (
    <svg
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        {...props}
    >
        <path
            d="M2 7H7V2H2V7ZM9.5 22H14.5V17H9.5V22ZM2 22H7V17H2V22ZM2 14.5H7V9.5H2V14.5ZM9.5 14.5H14.5V9.5H9.5V14.5ZM17 2V7H22V2H17ZM9.5 7H14.5V2H9.5V7ZM17 14.5H22V9.5H17V14.5ZM17 22H22V17H17V22Z"
            fill={color}
        />
    </svg>
);

export default ModuleStoreIcon;
