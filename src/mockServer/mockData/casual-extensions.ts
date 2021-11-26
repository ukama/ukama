import * as defaultCasual from "casual";
import {
    ALERT_TYPE,
    CONNECTED_USER_TYPE,
    DATA_PLAN_TYPE,
    GET_STATUS_TYPE,
    NETWORK_STATUS,
    NODE_TYPE,
} from "../../constants";
import { AlertDto } from "../../modules/alert/types";
import { BillHistoryDto, CurrentBillDto } from "../../modules/billing/types";
import { DataBillDto, DataUsageDto } from "../../modules/data/types";
import { EsimDto } from "../../modules/esim/types";
import { NetworkDto } from "../../modules/network/types";
import { NodeDto, UpdateNodeResponse } from "../../modules/node/types";

import {
    DeleteUserResponse,
    GetUserDto,
    UserDto,
    UserResponse,
} from "../../modules/user/types";

function randomArray<T>(
    minLength: number,
    maxLength: number,
    elementGenerator: () => T
): T[] {
    const length = casual.integer(minLength, maxLength);
    const result = [];
    for (let i = 0; i < length; i++) {
        result.push(elementGenerator());
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
        dataConsumed: defaultCasual.integer(1, 999),
        dataPackage: "Unlimited",
    };
};

const dataBill = (): DataBillDto => {
    return {
        id: defaultCasual._uuid(),
        dataBill: defaultCasual.integer(1, 999),
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
        status: defaultCasual.random_value(GET_STATUS_TYPE),
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
    const node = {
        Default: "Default",
        Intermediate: "Intermediate",
    };
    return {
        id: defaultCasual._uuid(),
        status: defaultCasual.random_value(GET_STATUS_TYPE),
        name: defaultCasual._name(),
        node: `${defaultCasual.random_value(node)} Data Plan`,
        dataPlan: defaultCasual.random_value(DATA_PLAN_TYPE),
        dataUsage: defaultCasual.integer(1, 199),
        dlActivity: "Table cell",
        ulActivity: "Table cell",
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
const deleteUser = (id: string): DeleteUserResponse => {
    return {
        id: id,
        success: true,
    };
};

interface Generators extends Casual.Generators {
    _randomArray: <T>(
        minLength: number,
        maxLength: number,
        elementGenerator: () => T
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
    _deleteUser: (id: string) => DeleteUserResponse;
    functions(): functions;
}
interface functions extends Casual.functions {
    randomArray: <T>(
        minLength: number,
        maxLength: number,
        elementGenerator: () => T
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
    deleteUser: (id: string) => DeleteUserResponse;
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
defaultCasual.define("deleteUser", deleteUser);

const casual = defaultCasual as Generators & functions & Casual.Casual;

export default casual;
