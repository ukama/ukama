import { selector } from "recoil";
import { isSkeltonLoading, pageName, organizationId } from "./atom";

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

export { withIsSkeltonLoading, withPageName, withOrganizationId };
