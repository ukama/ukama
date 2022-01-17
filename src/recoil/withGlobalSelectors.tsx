import { selector } from "recoil";
import {
    pageName,
    isDarkmode,
    isFirstVisit,
    organizationId,
    isSkeltonLoading,
} from "./atom";

const withIsSkeltonLoading = selector({
    key: "withIsSkeltonLoading",
    get: ({ get }) => get(isSkeltonLoading),
});

const withPageName = selector({
    key: "withPageName",
    get: ({ get }) => get(pageName),
});

const withOrganizationId = selector({
    key: "withOrganizationId",
    get: ({ get }) => get(organizationId),
});

const withIsFirstVisit = selector({
    key: "withIsFirstVisit",
    get: ({ get }) => get(isFirstVisit),
});

const withIsDarkMod = selector({
    key: "withIsDarkMod",
    get: ({ get }) => get(isDarkmode),
});

export {
    withPageName,
    withIsDarkMod,
    withIsFirstVisit,
    withOrganizationId,
    withIsSkeltonLoading,
};
