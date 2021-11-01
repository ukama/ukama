const UppercaseRegex = /[A-Z]/;
const LowercaseRegex = /[a-z]/;
//eslint-disable-next-line
const specialCharactersRegex = /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;
const alphabetRegex = /[a-zA-Z]/g;
const checkPasswordHasLowercaseLetters = (password: string) =>
    LowercaseRegex.test(password);

const checkPasswordHasUppercaseLetters = (password: string) =>
    UppercaseRegex.test(password);
const getTitleFromPath = (path: string) => {
    switch (path) {
        case "/":
            return "Home";
        case "/home":
            return "Home";
        case "/nodes":
            return "Nodes";
        case "/user":
            return "User";
        case "/billing":
            return "Billing";
        case "/store":
            return "Module Store";
        default:
            return "Home";
    }
};

const checkPasswordHasLetters = (password: string) =>
    alphabetRegex.test(password);

const checkPasswordSpecialCharacter = (password: string) =>
    specialCharactersRegex.test(password);

const checkPasswordLength = (password: string) => password.length > 7;
export {
    checkPasswordHasLowercaseLetters,
    checkPasswordHasUppercaseLetters,
    alphabetRegex,
    getTitleFromPath,
    checkPasswordLength,
    specialCharactersRegex,
    checkPasswordHasLetters,
    checkPasswordSpecialCharacter,
};
