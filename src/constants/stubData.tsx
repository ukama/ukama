import { CloudOffIcon } from "../assets/svg";
const RechartsData = [
    {
        name: "10 AM",
        uv: 590,
        pv: 800,
        amt: 1400,
        cnt: 490,
    },
    {
        name: "8 PM",
        uv: 868,
        pv: 967,
        amt: 1506,
        cnt: 590,
    },
    {
        name: "12 PM",
        uv: 1397,
        pv: 1098,
        amt: 989,
        cnt: 350,
    },
    {
        name: "4 AM",
        uv: 1480,
        pv: 1200,
        amt: 1228,
        cnt: 480,
    },
    {
        name: "6 AM",
        uv: 1520,
        pv: 1108,
        amt: 1100,
        cnt: 460,
    },
    {
        name: "10 PM",
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

const SimCardData = [
    {
        id: 1,
        title: "ESIM # 1",
        serial: "# 32019-08-12-2023-999990",
        isActive: true,
    },
    {
        id: 2,
        title: "ESIM # 2",
        serial: "# 32019-08-12-2023-0225016",
        isActive: false,
    },
    {
        id: 3,
        title: "ESIM # 3",
        serial: "# 32019-08-12-2023-8888800",
        isActive: false,
    },
    {
        id: 4,
        title: "ESIM # 4",
        serial: "# 32019-08-12-2023-6666660",
        isActive: false,
    },
];

export {
    SimCardData,
    RechartsData,
    ALERT_INFORMATION,
    DashboardSliderData,
    DashboardResidentsTable,
};
