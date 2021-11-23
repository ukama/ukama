import * as defaultCasual from "casual";
import {
    ALERT_TYPE,
    CONNECTED_USER_TYPE,
    DATA_PLAN_TYPE,
    GET_STATUS_TYPE,
    NODE_TYPE,
} from "../../constants";
import { AlertDto } from "../../modules/alert/types";
import { DataBillDto, DataUsageDto } from "../../modules/data/types";
import { EsimDto } from "../../modules/esim/types";
import { NodeDto } from "../../modules/node/types";

import { GetUserDto, UserDto } from "../../modules/user/types";

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
}

defaultCasual.define("randomArray", randomArray);
defaultCasual.define("user", user);
defaultCasual.define("dataUsage", dataUsage);
defaultCasual.define("dataBill", dataBill);
defaultCasual.define("alert", alert);
defaultCasual.define("node", node);
defaultCasual.define("esim", esim);
defaultCasual.define("getUser", getUser);

const casual = defaultCasual as Generators & functions & Casual.Casual;

export default casual;
