import {
    checkPasswordLength,
    checkPasswordSpecialCharacter,
    checkPasswordHasLowercaseLetters,
    checkPasswordHasUppercaseLetters,
} from "../utils";

const PasswordRules = [
    {
        id: 1,
        idLabel: "checkLength",
        label: "Be a minimum of 8 characters",
        validator: checkPasswordLength,
    },
    {
        id: 2,
        idLabel: "checkSpecialCaracter",
        label: "At least one special character",
        validator: checkPasswordSpecialCharacter,
    },
    {
        id: 3,
        idLabel: "checkUpperCase",
        label: "Must have Uppercase letter",
        validator: checkPasswordHasUppercaseLetters,
    },
    {
        id: 4,
        idLabel: "checkLowerCase",
        label: "Must have Lowercase letter",
        validator: checkPasswordHasLowercaseLetters,
    },
];
export { PasswordRules };
