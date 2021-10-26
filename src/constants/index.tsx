import {
    checkPasswordLength,
    checkPasswordSpecialCharacter,
    checkPasswordHasLowercaseLetters,
    checkPasswordHasUppercaseLetters,
} from "../utils";

const PasswordRules = [
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
        label: "Must have Uppercase letter",
        validator: checkPasswordHasUppercaseLetters,
    },
    {
        id: 4,
        label: "Must have Lowercase letter",
        validator: checkPasswordHasLowercaseLetters,
    },
];
export { PasswordRules };
