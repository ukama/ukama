import { selector } from "recoil";
import { isSkeltonLoading } from "./atom";

const withIsSkeltonLoading = selector({
    key: "withIsSkeltonLoading",
    get: ({ get }) => get(isSkeltonLoading),
});

export { withIsSkeltonLoading };
