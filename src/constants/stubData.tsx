const RechartsData = [
    {
        name: "Page A",
        uv: 4000,
        pv: 2400,
        vx: 1000,
    },
    {
        name: "Page B",
        uv: 3000,
        pv: 1398,
        vx: 1000,
    },
    {
        name: "Page C",
        uv: 2000,
        pv: 9800,
        vx: 4000,
    },
    {
        name: "Page D",
        uv: 2780,
        pv: 3908,
        vx: 9000,
    },
    {
        name: "Page E",
        uv: 1890,
        pv: 4800,
        vx: 2400,
    },
    {
        name: "Page F",
        uv: 2390,
        pv: 3800,
        vx: 4500,
    },
    {
        name: "Page G",
        uv: 3490,
        pv: 4300,
        vx: 8000,
    },
];

const CREDIT_CARD = [
    {
        id: 1,

        card_experintionDetails: "Card is ending in 20 days",
    },
];
const CurrentBillingData = [
    {
        id: 1,
        name: "Tryphena Nelson 1",
        dataUsage: " 20 GB  ",
        rate: "$5 / 20GB",
        subTotal: 10,
    },
    {
        id: 2,
        name: "Tryphena Nelson 2",
        dataUsage: " 20 GB  ",
        rate: "$5 / 20GB",
        subTotal: 10,
    },
    {
        id: 3,
        name: "Tryphena Nelson 3",
        dataUsage: " 20 GB  ",
        rate: "$5 / 20GB",
        subTotal: 10,
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
const UserData = [
    {
        id: 1,
        iccid: "983098214-329-323",
        email: "pcassinga@gmail.com",
        phone: "111-111-111-111",
        roaming: true,
        name: "ESIM # 1",
        eSimNumber: "32019-08-12-2023-999990",
        dataPlan: 10,
        dataUsage: 2,
    },
    {
        id: 2,
        iccid: "233098214-329-323",
        email: "joe@gmail.com",
        phone: "111-111-111-111",
        roaming: true,
        name: "ESIM # 1",
        eSimNumber: "32019-08-12-2023-999990",
        dataPlan: 10,
        dataUsage: 2,
    },
    {
        id: 3,
        iccid: "983098214-329-323",
        email: "ukama@gmail.com",
        phone: "111-111-111-111",
        roaming: true,
        name: "ESIM # 1",
        eSimNumber: "32019-08-12-2023-999990",
        dataPlan: 10,
        dataUsage: 2,
    },
    {
        id: 4,
        iccid: "983098214-329-323",
        email: "BillJohn@gmail.com",
        phone: "111-111-111-111",
        roaming: true,
        name: "ESIM # 1",
        eSimNumber: "32019-08-12-2023-999990",
        dataPlan: 10,
        dataUsage: 2,
    },
];
const CurrentBilling = [
    {
        id: 1,
        name: "Tryphena Nelson ",
        dataUsage: " 20 GB  ",
        rate: "$5 / 20GB",
        subTotal: 10,
    },
    {
        id: 2,
        name: "Tryphena Nelson ",
        dataUsage: " 20 GB  ",
        rate: "$5 / 20GB",
        subTotal: 10,
    },
    {
        id: 3,
        name: "Tryphena Nelson ",
        dataUsage: " 20 GB  ",
        rate: "$5 / 20GB",
        subTotal: 10,
    },
    {
        id: 4,
        name: "Tryphena Nelson ",
        dataUsage: " 20 GB  ",
        rate: "$5 / 20GB",
        subTotal: 10,
    },
];

const NodesData = [
    {
        name: "Node X",
        statusType: "ONLINE",
        duration: "21 days 5 hours 1 minute",
    },
    {
        name: "Node Y",
        statusType: "BEING_CONFIGURED",
        duration: "21 days 5 hours 1 minute",
    },
    {
        name: "Node Z",
        statusType: "",
        duration: "21 days 5 hours 1 minute",
    },
];

const NODE_PROPERTIES8 = [
    {
        name: "Model type",
        value: "Home Node",
    },
    {
        name: "Serial #",
        value: "1111111111111111111",
    },
    {
        name: "MAC address",
        value: "1111111111111111111",
    },
    { name: "OS version", value: "1.0" },
    {
        name: "Manufacturing ",
        value: "1209391023209103",
    },
    { name: "Ukama OS", value: "1.0" },
    { name: "Hardware", value: "1.0" },
    {
        name: "Description",
        value: "Home node is a xyz.",
    },
];
const NODE_PROPERTIES2 = [
    {
        name: "Model type",
        value: "Home Node",
    },
    {
        name: "Serial #",
        value: "1111111111111111111",
    },
];
const NODE_PROPERTIES3 = [
    {
        name: "Model type",
        value: "Home Node",
    },
    {
        name: "Serial #",
        value: "1111111111111111111",
    },
    {
        name: "Description",
        value: "Home node is a xyz.",
    },
];

const NODE_PROPERTIES4 = [
    {
        name: "Model type",
        value: "Home Node",
    },
    {
        name: "Serial #",
        value: "1111111111111111111",
    },
    {
        name: "Description",
        value: "Home node is a xyz.",
    },
    { name: "Ukama OS", value: "1.0" },
];

const NODES = [
    {
        id: "1",
        totalUser: 4,
        title: "Node 1",
        status: "PENDING",
        description: "Node 1 description",
    },
    {
        id: "2",
        totalUser: 2,
        title: "Node 2",
        status: "PENDING",
        description: "Node 2 description",
    },
    {
        id: "3",
        totalUser: 6,
        title: "Node 3",
        status: "ONGOING",
        description: "Node 4 description",
    },
];

export {
    SimCardData,
    CurrentBillingData,
    CurrentBilling,
    RechartsData,
    UserData,
    CREDIT_CARD,
    NodesData,
    NODES,
    NODE_PROPERTIES2,
    NODE_PROPERTIES3,
    NODE_PROPERTIES8,
    NODE_PROPERTIES4,
};
