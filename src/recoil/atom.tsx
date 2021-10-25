import { atom } from "recoil";

const isLoginAtom = atom({
    key: "isLoginAtom",
    default: false,
});

export { isLoginAtom };
