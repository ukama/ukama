import * as defaultCasual from "casual";
import { CONNECTED_USER_TYPE } from "../../constants";
import { DataUsageDto } from "../../modules/data/types";

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

interface Generators extends Casual.Generators {
    _randomArray: <T>(
        minLength: number,
        maxLength: number,
        elementGenerator: () => T
    ) => Array<T>;

    _user: () => UserDto;
    _dataUsage: () => DataUsageDto;
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
}

defaultCasual.define("randomArray", randomArray);
defaultCasual.define("user", user);
defaultCasual.define("dataUsage", dataUsage);

const casual = defaultCasual as Generators & functions & Casual.Casual;

export default casual;
