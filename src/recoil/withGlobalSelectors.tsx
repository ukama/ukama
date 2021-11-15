import { selector } from "recoil";
import { isLoginAtom, isSkeltonLoading } from "./atom";

const withIsLogin = selector({
    key: "withIsLogin",
    get: ({ get }) => get(isLoginAtom),
});

const withIsSkeltonLoading = selector({
    key: "withIsSkeltonLoading",
    get: ({ get }) => get(isSkeltonLoading),
});

export { withIsLogin, withIsSkeltonLoading };
