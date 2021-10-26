const UppercaseRegex = /(?=.*[A-Z])/g;
const LowercaseRegex = /(?=.*[a-z])/g;
//eslint-disable-next-line
const specialCharactersRegex = /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;

const checkPasswordHasLowercaseLetters = (password: string) =>
    LowercaseRegex.test(password);

const checkPasswordHasUppercaseLetters = (password: string) =>
    UppercaseRegex.test(password);

const checkPasswordSpecialCharacter = (password: string) =>
    specialCharactersRegex.test(password);

const checkPasswordLength = (password: string) => password.length > 7;
export {
    checkPasswordLength,
    checkPasswordSpecialCharacter,
    checkPasswordHasLowercaseLetters,
    checkPasswordHasUppercaseLetters,
};
