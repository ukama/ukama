export type MenuItemType = {
    Icon: any;
    id: number;
    title: string;
    route: string;
};
export type StatsItemType = {
    id: number;
    label: string;
    value: string;
};
export type statsPeriodItemType = {
    id: number;
    label: string;
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
    id: "name" | "actions" | "dataUsage";
    label: string;
    minWidth?: number;
    align?: "right";
};

export type SimActivateFormType = {
    email: string;
    phone: string;
    number: string;
    lastName: string;
    firstName: string;
};
