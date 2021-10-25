import { selector } from "recoil";
import { isLoginAtom } from "./atom";

const withIsLogin = selector({
    key: "withIsLogin",
    get: ({ get }) => get(isLoginAtom),
});

export { withIsLogin };
