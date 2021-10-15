const checkPasswordHasLetters = (password: string) =>
    alphabetRegex.test(password);

const checkPasswordSpecialCharacter = (password: string) =>
    specialCharactersRegex.test(password);

const checkPasswordLength = (password: string) => password.length > 7;
//eslint-disable-next-line
const specialCharactersRegex = /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;
//eslint-disable-next-line
const alphabetRegex = /[a-zA-Z]/g;
export {
    specialCharactersRegex,
    alphabetRegex,
    checkPasswordHasLetters,
    checkPasswordSpecialCharacter,
    checkPasswordLength,
};
