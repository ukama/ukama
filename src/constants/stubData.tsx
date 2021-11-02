import { MONTH_FILTER, TIME_FILTER } from ".";
import { DataBilling, DataUsage, UsersWithBG } from "../assets/svg";

const DashboardStatusCard = [
    {
        id: "statusUser",
        subtitle1: "13",
        Icon: UsersWithBG,
        options: TIME_FILTER,
        title: "Connected Users",
        subtitle2: "| 10 residents; 3 guests",
    },
    {
        id: "statusUsage",
        subtitle1: "0 GBs",
        Icon: DataUsage,
        options: TIME_FILTER,
        title: "Data usage",
        subtitle2: "/ unlimited",
    },
    {
        id: "statusBill",
        subtitle1: "$ 20",
        Icon: DataBilling,
        options: MONTH_FILTER,
        title: "Data Bill",
        subtitle2: "/ due in 8 days",
    },
];

const DashboardSliderData = [
    {
        id: 1,
        title: "Tryphena’s Node",
        subTitle: "Home node",
        users: "3",
        isConfigure: false,
    },
    {
        id: 2,
        title: "Tryphena’s Node",
        subTitle: "Home node",
        users: "4",
        isConfigure: false,
    },
    {
        id: 3,
        title: "Tryphena’s Node",
        subTitle: "Home node",
        users: "5",
        isConfigure: false,
    },
    {
        id: 4,
        title: "Tryphena’s Node",
        subTitle: "Home node",
        users: "6",
        isConfigure: false,
    },
];

export { DashboardStatusCard, DashboardSliderData };
