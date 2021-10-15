// This file will containt app constants
import {
    checkPasswordSpecialCharacter,
    checkPasswordLength,
    checkPasswordHasLetters,
} from "../utils";
const passwordRules = [
    {
        id: 1,
        label: "Be a minimum of 8 characters",
        validator: checkPasswordLength,
    },
    {
        id: 2,
        label: "At least one special character",
        validator: checkPasswordSpecialCharacter,
    },
    {
        id: 3,
        label: "Upper & lowercase letters ",
        validator: checkPasswordHasLetters,
    },
];
export { passwordRules };
