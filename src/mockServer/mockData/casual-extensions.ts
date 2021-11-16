import * as defaultCasual from "casual";
import { ALERT_TYPE, CONNECTED_USER_TYPE } from "../../constants";
import { AlertDto } from "../../modules/alert/types";
import { DataBillDto, DataUsageDto } from "../../modules/data/types";
import { NodeDto } from "../../modules/node/types";
import { ResidentDto } from "../../modules/resident/types";

import { UserDto } from "../../modules/user/types";

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
        dataConsumed: `${defaultCasual.integer(1, 999)}GBs`,
        dataPackage: "Unlimited",
    };
};

const dataBill = (): DataBillDto => {
    return {
        id: defaultCasual._uuid(),
        dataBill: `${defaultCasual.integer(1, 999)}$`,
        billDue: `${defaultCasual.integer(1, 29)} days`,
    };
};

const alert = (): AlertDto => {
    return {
        id: defaultCasual._uuid(),
        type: defaultCasual.random_value(ALERT_TYPE),
        title: defaultCasual._title(),
        description: defaultCasual._description(),
        alertDate: new Date(),
    };
};

const node = (): NodeDto => {
    return {
        id: defaultCasual._uuid(),
        title: defaultCasual._title(),
        description: defaultCasual._description(),
        totalUser: defaultCasual.integer(1, 99),
    };
};
const resident = (): ResidentDto => {
    return {
        id: defaultCasual._uuid(),
        name: defaultCasual._name(),
        usage: `${defaultCasual.integer(1, 999)}GB`,
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
    _resident: () => ResidentDto;
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
    resident: () => ResidentDto;
}

defaultCasual.define("randomArray", randomArray);
defaultCasual.define("user", user);
defaultCasual.define("dataUsage", dataUsage);
defaultCasual.define("dataBill", dataBill);
defaultCasual.define("alert", alert);
defaultCasual.define("node", node);
defaultCasual.define("resident", resident);

const casual = defaultCasual as Generators & functions & Casual.Casual;

export default casual;
