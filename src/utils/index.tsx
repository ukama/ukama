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
//eslint-disable-next-line
const specialCharactersRegex = /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;
//eslint-disable-next-line
const alphabetRegex = /[a-zA-Z]/g;
export {
    alphabetRegex,
    getTitleFromPath,
    checkPasswordLength,
    specialCharactersRegex,
    checkPasswordHasLetters,
    checkPasswordSpecialCharacter,
};
