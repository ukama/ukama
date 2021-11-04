import { MONTH_FILTER, TIME_FILTER } from ".";
import {
    DataBilling,
    DataUsage,
    UsersWithBG,
    CloudOffIcon,
} from "../assets/svg";
const RechartsData = [
    {
        name: "Page A",
        uv: 590,
        pv: 800,
        amt: 1400,
        cnt: 490,
    },
    {
        name: "Page B",
        uv: 868,
        pv: 967,
        amt: 1506,
        cnt: 590,
    },
    {
        name: "Page C",
        uv: 1397,
        pv: 1098,
        amt: 989,
        cnt: 350,
    },
    {
        name: "Page D",
        uv: 1480,
        pv: 1200,
        amt: 1228,
        cnt: 480,
    },
    {
        name: "Page E",
        uv: 1520,
        pv: 1108,
        amt: 1100,
        cnt: 460,
    },
    {
        name: "Page F",
        uv: 1400,
        pv: 680,
        amt: 1700,
        cnt: 380,
    },
];
const ALERT_INFORMATION = [
    {
        id: 1,
        Icon: CloudOffIcon,
        title: "Software error",
        description: "Short description of alert.",
        date: "08/16/21 1PM",
    },
    {
        id: 2,
        Icon: CloudOffIcon,
        title: "Software error",
        description: "Short description of alert.",
        date: "08/16/21 1PM",
    },
    {
        id: 3,
        Icon: CloudOffIcon,
        title: "Software error",
        description: "Short description of alert.",
        date: "08/16/21 1PM",
    },
    {
        id: 4,
        Icon: CloudOffIcon,
        title: "Software error",
        description: "Short description of alert.",
        date: "08/16/21 1PM",
    },
    {
        id: 5,
        Icon: CloudOffIcon,
        title: "Software error",
        description: "Short description of alert.",
        date: "08/16/21 1PM",
    },
    {
        id: 6,
        Icon: CloudOffIcon,
        title: "Software error",
        description: "Short description of alert.",
        date: "08/16/21 1PM",
    },
];
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

const DashboardResidentsTable = [
    {
        name: "Tryphena Nelson 1",
        usage: "1 GB",
        actions: "",
    },
    {
        name: "Tryphena Nelson 2",
        usage: "1 GB",
        actions: "",
    },
    {
        name: "Tryphena Nelson 3",
        usage: "1 GB",
        actions: "",
    },
    {
        name: "Tryphena Nelson 4",
        usage: "1 GB",
        actions: "",
    },
    {
        name: "Tryphena Nelson 5",
        usage: "1 GB",
        actions: "",
    },
];

const DashboardSliderData = [
    {
        id: 1,
        title: "Tryphena’s Node 1",
        subTitle: "Home node",
        users: "3",
        isConfigure: false,
    },
    {
        id: 2,
        title: "Tryphena’s Node 2",
        subTitle: "Home node",
        users: "4",
        isConfigure: false,
    },
    {
        id: 3,
        title: "Tryphena’s Node 3",
        subTitle: "Home node",
        users: "5",
        isConfigure: false,
    },
    {
        id: 4,
        title: "Tryphena’s Node 4",
        subTitle: "Home node",
        users: "6",
        isConfigure: false,
    },
];

export {
    DashboardStatusCard,
    RechartsData,
    ALERT_INFORMATION,
    DashboardSliderData,
    DashboardResidentsTable,
};
