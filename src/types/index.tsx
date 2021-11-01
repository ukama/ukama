export type MenuItemType = {
    Icon: any;
    id: number;
    title: string;
    route: string;
};

export type HeaderMenuItemType = {
    id: string;
    Icon: any;
    title: string;
};

export type SelectItemType = {
    id: number;
    label: string;
    value: string;
};

export type SVGType = {
    color?: string;
    width?: string;
    height?: string;
};

export type ColumnsWithOptions = {
    id: "name" | "usage";
    label: string;
    minWidth?: number;
    align?: "right";
    format?: (value: number) => string;
};
