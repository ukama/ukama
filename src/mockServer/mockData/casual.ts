import * as defaultCasual from "casual";
import {
    ALERT_TYPE,
    CONNECTED_USER_TYPE,
    ORG_NODE_STATE,
    GET_STATUS_TYPE,
    NETWORK_STATUS,
    NODE_TYPE,
} from "../../constants";
import { AlertDto } from "../../modules/alert/types";
import { BillHistoryDto, CurrentBillDto } from "../../modules/billing/types";
import { DataBillDto, DataUsageDto } from "../../modules/data/types";
import { EsimDto } from "../../modules/esim/types";
import { NetworkDto } from "../../modules/network/types";
import {
    CpuUsageMetricsDto,
    ThroughputMetricsDto,
    IOMetricsDto,
    NodeDetailDto,
    NodeDto,
    NodeMetaDataDto,
    NodePhysicalHealthDto,
    NodeRFDto,
    TemperatureMetricsDto,
    MemoryUsageMetricsDto,
    UpdateNodeResponse,
} from "../../modules/node/types";

import {
    UsersAttachedMetricsDto,
    DeactivateResponse,
    GetUserDto,
    UserDto,
    UserResponse,
} from "../../modules/user/types";
import { getUniqueTimeStamp } from "../../utils";

function randomArray<T>(
    minLength: number,
    maxLength: number,
    elementGenerator: (index?: number, length?: number) => T
): T[] {
    const length = casual.integer(minLength, maxLength);
    const result = [];
    for (let i = 0; i < length; i++) {
        result.push(elementGenerator(i + 1, maxLength));
    }
    return result;
}
const user = (): UserDto => {
    return {
        id: defaultCasual._uuid(),
        name: defaultCasual._name(),
        email: defaultCasual._email(),
        type: defaultCasual.random_value(CONNECTED_USER_TYPE),
    };
};
const dataUsage = (): DataUsageDto => {
    return {
        id: defaultCasual._uuid(),
        dataConsumed: defaultCasual.integer(1, 39),
        dataPackage: `${defaultCasual.integer(5, 60)} GB free left`,
    };
};

const dataBill = (): DataBillDto => {
    return {
        id: defaultCasual._uuid(),
        dataBill: defaultCasual.integer(1, 39),
        billDue: defaultCasual.integer(1, 29),
    };
};

const alert = (): AlertDto => {
    return {
        id: defaultCasual._uuid(),
        type: defaultCasual.random_value(ALERT_TYPE),
        title: defaultCasual._title(),
        description: defaultCasual._short_description(),
        alertDate: new Date(),
    };
};

const node = (): NodeDto => {
    return {
        id: defaultCasual._uuid(),
        title: defaultCasual._title(),
        description: `${defaultCasual.random_value(NODE_TYPE)} node`,
        status: defaultCasual.random_value(ORG_NODE_STATE),
        totalUser: defaultCasual.integer(1, 99),
    };
};
const updateNode = (
    id: string,
    name: string,
    serialNo: string
): UpdateNodeResponse => {
    return {
        id: id,
        name: name ?? defaultCasual._name(),
        serialNo: serialNo ?? `#${defaultCasual.integer(1111111, 9999999)}`,
    };
};

const getUser = (): GetUserDto => {
    return {
        id: defaultCasual._uuid(),
        status: defaultCasual.random_value(GET_STATUS_TYPE),
        name: defaultCasual._name(),
        eSimNumber: `# ${defaultCasual.integer(
            11111,
            99999
        )}-${defaultCasual.date("DD-MM-YYYY")}-${defaultCasual.integer(
            1111111,
            9999999
        )}`,
        iccid: `${defaultCasual.integer(11111, 99999)}${defaultCasual.integer(
            11010,
            99999
        )}${defaultCasual.integer(11010, 99999)}`,
        email: defaultCasual._email(),
        phone: defaultCasual._phone(),
        roaming: defaultCasual.random_value([true, false]),
        dataPlan: defaultCasual.integer(3, 8),
        dataUsage: defaultCasual.integer(1, 5),
    };
};
const esim = (): EsimDto => {
    const boolean = {
        true: true,
        false: false,
    };
    return {
        esim: `# ${defaultCasual.integer(11111, 99999)}-${defaultCasual.date(
            "DD-MM-YYYY"
        )}-${defaultCasual.integer(1111111, 9999999)}`,
        active: defaultCasual.random_value(boolean),
    };
};
const network = (): NetworkDto => {
    const status = defaultCasual.random_value(NETWORK_STATUS);

    if (status === NETWORK_STATUS.BEING_CONFIGURED)
        return {
            id: defaultCasual._uuid(),
            status,
        };
    return {
        id: defaultCasual._uuid(),
        status,
        description: "21 days 5 hours 1 minute",
    };
};

const currentBill = (): CurrentBillDto => {
    const data = defaultCasual.integer(1, 10);
    const rate = defaultCasual.integer(3, 6);
    const subtotal = data * rate;
    return {
        id: defaultCasual._uuid(),
        name: defaultCasual._name(),
        dataUsed: data,
        rate: rate,
        subtotal: subtotal,
    };
};

const billHistory = (): BillHistoryDto => {
    const totalUsage = defaultCasual.integer(1, 10);
    const subtotal = totalUsage * 3;
    return {
        id: defaultCasual._uuid(),
        date: defaultCasual.date("MM-DD-2021"),
        description: `Bill for month`,
        totalUsage,
        subtotal: subtotal,
    };
};
const updateUser = (
    id: string,
    firstName: string,
    lastName: string,
    eSimNumber: string,
    email: string,
    phone: string
): UserResponse => {
    return {
        id: id,
        name: `${firstName ?? defaultCasual._first_name()} ${
            lastName ?? defaultCasual._last_name()
        }`,
        sim:
            eSimNumber ??
            `# ${defaultCasual.integer(11111, 99999)}-${defaultCasual.date(
                "DD-MM-2023"
            )}-${defaultCasual.integer(1111111, 9999999)}`,
        email: email ?? defaultCasual._email(),
        phone: phone ?? defaultCasual._phone(),
    };
};
const deleteRes = (id: string): DeactivateResponse => {
    return {
        id: id,
        success: true,
    };
};
const nodeDetail = (): NodeDetailDto => {
    return {
        id: defaultCasual._uuid(),
        modelType: `${defaultCasual.random_value(NODE_TYPE)} Node`,
        serial: defaultCasual.integer(1111111111111111111, 9999999999999999999),
        macAddress: defaultCasual.integer(
            1111111111111111111,
            9999999999999999999
        ),
        osVersion: defaultCasual.integer(1, 9),
        manufacturing: defaultCasual.integer(
            1111111111111111,
            9999999999999999
        ),
        ukamaOS: defaultCasual.integer(1, 9),
        hardware: defaultCasual.integer(1, 9),
        description: `${defaultCasual.random_value(NODE_TYPE)} node is a xyz`,
    };
};
const nodeMetaData = (): NodeMetaDataDto => {
    return {
        throughput: defaultCasual.integer(1, 19),
        usersAttached: defaultCasual.integer(1, 9),
    };
};
const nodePhysicalHealth = (): NodePhysicalHealthDto => {
    return {
        temperature: defaultCasual.integer(1, 19),
        Memory: defaultCasual.integer(1, 19),
        cpu: defaultCasual.integer(1, 19),
        io: defaultCasual.integer(1, 19),
    };
};
const nodeNetwork = (): NetworkDto => {
    return {
        id: defaultCasual._uuid(),
        status: NETWORK_STATUS.ONLINE,
        description: "21 days 5 hours 1 minute",
    };
};
const throughputMetrics = (): ThroughputMetricsDto => {
    const time = {
        AM: "AM",
        PM: "PM",
    };
    return {
        uv: defaultCasual.integer(50, 500),
        pv: defaultCasual.integer(50, 500),
        amt: defaultCasual.integer(150, 1500),
        time: `${defaultCasual.integer(1, 12)} ${defaultCasual.random_value(
            time
        )}`,
    };
};

const usersAttachedMetrics = (
    index: number,
    length: number
): UsersAttachedMetricsDto => {
    return {
        id: defaultCasual._uuid(),
        users: defaultCasual.integer(1, 50),
        timestamp: getUniqueTimeStamp(index, length),
    };
};

const cpuUsageMetrics = (index: number, length: number): CpuUsageMetricsDto => {
    return {
        id: defaultCasual._uuid(),
        usage: defaultCasual.integer(1, 200),
        timestamp: getUniqueTimeStamp(index, length),
    };
};

const nodeRF = (index: number, length: number): NodeRFDto => {
    return {
        qam: defaultCasual.integer(1, 19),
        rfOutput: defaultCasual.integer(1, 19),
        rssi: defaultCasual.integer(1, 19),
        timestamp: getUniqueTimeStamp(index, length),
    };
};

const temperatureMetrics = (
    index: number,
    length: number
): TemperatureMetricsDto => {
    return {
        id: defaultCasual._uuid(),
        temperature: defaultCasual.integer(1, 150),
        timestamp: getUniqueTimeStamp(index, length),
    };
};

const getIOMetrics = (index: number, length: number): IOMetricsDto => {
    return {
        id: defaultCasual._uuid(),
        input: defaultCasual.integer(1, 200),
        output: defaultCasual.integer(1, 200),
        timestamp: getUniqueTimeStamp(index, length),
    };
};

const memoryUsageMetrics = (
    index: number,
    length: number
): MemoryUsageMetricsDto => {
    return {
        id: defaultCasual._uuid(),
        usage: defaultCasual.integer(1, 4096),
        timestamp: getUniqueTimeStamp(index, length),
    };
};

interface Generators extends Casual.Generators {
    _randomArray: <T>(
        minLength: number,
        maxLength: number,
        elementGenerator: (index?: number) => T
    ) => Array<T>;

    _user: () => UserDto;
    _dataUsage: () => DataUsageDto;
    _dataBill: () => DataBillDto;
    _alert: () => AlertDto;
    _node: () => NodeDto;
    _esim: () => EsimDto;
    _getUser: () => GetUserDto;
    _currentBill: () => CurrentBillDto;
    _billHistory: () => BillHistoryDto;
    _network: () => NetworkDto;
    _updateNode: (
        id: string,
        name: string,
        serialNo: string
    ) => UpdateNodeResponse;
    _updateUser: (
        id: string,
        firstName: string,
        lastName: string,
        eSimNumber: string,
        email: string,
        phone: string
    ) => UserResponse;
    _deleteRes: (id: string) => DeactivateResponse;
    _nodeDetail: () => NodeDetailDto;
    _nodeMetaData: () => NodeMetaDataDto;
    _nodePhysicalHealth: () => NodePhysicalHealthDto;
    _nodeRF: () => NodeRFDto;
    _nodeNetwork: () => NetworkDto;
    _throughputMetrics: () => ThroughputMetricsDto;
    _cpuUsageMetrics: () => CpuUsageMetricsDto;
    _memoryUsageMetrics: () => MemoryUsageMetricsDto;
    _usersAttachedMetrics: () => UsersAttachedMetricsDto;
    _temperatureMetrics: () => TemperatureMetricsDto;
    _ioMetrics: () => IOMetricsDto;
    functions(): Functions;
}
interface Functions extends Casual.functions {
    randomArray: <T>(
        minLength: number,
        maxLength: number,
        elementGenerator: (index?: number) => T
    ) => Array<T>;
    user: () => UserDto;
    dataUsage: () => DataUsageDto;
    dataBill: () => DataBillDto;
    alert: () => AlertDto;
    node: () => NodeDto;
    getUser: () => GetUserDto;
    esim: () => EsimDto;
    currentBill: () => CurrentBillDto;
    billHistory: () => BillHistoryDto;
    network: () => NetworkDto;
    updateNode: (
        id: string,
        name: string,
        serialNo: string
    ) => UpdateNodeResponse;
    updateUser: (
        id: string,
        firstName: string,
        lastName: string,
        eSimNumber: string,
        email: string,
        phone: string
    ) => UserResponse;
    deleteRes: (id: string) => DeactivateResponse;
    nodeDetail: () => NodeDetailDto;
    nodeMetaData: () => NodeMetaDataDto;
    nodePhysicalHealth: () => NodePhysicalHealthDto;
    nodeRF: () => NodeRFDto;
    nodeNetwork: () => NetworkDto;
    throughputMetrics: () => ThroughputMetricsDto;
    cpuUsageMetrics: () => CpuUsageMetricsDto;
    memoryUsageMetrics: () => MemoryUsageMetricsDto;
    usersAttachedMetrics: () => UsersAttachedMetricsDto;
    temperatureMetrics: () => TemperatureMetricsDto;
    ioMetrics: () => IOMetricsDto;
}

defaultCasual.define("randomArray", randomArray);
defaultCasual.define("user", user);
defaultCasual.define("dataUsage", dataUsage);
defaultCasual.define("dataBill", dataBill);
defaultCasual.define("alert", alert);
defaultCasual.define("node", node);
defaultCasual.define("esim", esim);
defaultCasual.define("getUser", getUser);
defaultCasual.define("currentBill", currentBill);
defaultCasual.define("billHistory", billHistory);
defaultCasual.define("network", network);
defaultCasual.define("updateNode", updateNode);
defaultCasual.define("updateUser", updateUser);
defaultCasual.define("deleteRes", deleteRes);
defaultCasual.define("nodeDetail", nodeDetail);
defaultCasual.define("nodeMetaData", nodeMetaData);
defaultCasual.define("nodePhysicalHealth", nodePhysicalHealth);
defaultCasual.define("nodeRF", nodeRF);
defaultCasual.define("nodeNetwork", nodeNetwork);
defaultCasual.define("throughputMetrics", throughputMetrics);
defaultCasual.define("cpuUsageMetrics", cpuUsageMetrics);
defaultCasual.define("temperatureMetrics", temperatureMetrics);
defaultCasual.define("memoryUsageMetrics", memoryUsageMetrics);
defaultCasual.define("usersAttachedMetrics", usersAttachedMetrics);
defaultCasual.define("ioMetrics", getIOMetrics);
const casual = defaultCasual as Generators & Functions & Casual.Casual;

export default casual;
