import { selector } from "recoil";
import { isSkeltonLoading, pageName } from "./atom";

const withIsSkeltonLoading = selector({
    key: "withIsSkeltonLoading",
    get: ({ get }) => get(isSkeltonLoading),
});

const withPageName = selector({
    key: "withPageName",
    get: ({ get }) => get(pageName),
});

export { withIsSkeltonLoading, withPageName };
