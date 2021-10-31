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
export type AlertItemType = {
    id: number;
    Icon: any;
    title: string;
    date: string;
    description: string;
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
