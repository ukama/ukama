import { atom } from "recoil";
import { recoilPersist } from "recoil-persist";
const { persistAtom } = recoilPersist();

const isFirstVisit = atom({
    key: "isFirstVisit",
    default: true,
    effects_UNSTABLE: [persistAtom],
});

const isSkeltonLoading = atom({
    key: "isSkeltonLoading",
    default: false,
    effects_UNSTABLE: [persistAtom],
});

const pageName = atom({
    key: "pageName",
    default: "Home",
    effects_UNSTABLE: [persistAtom],
});

const organizationId = atom<string | undefined>({
    key: "organizationId",
    default: undefined,
    effects_UNSTABLE: [persistAtom],
});

export { isSkeltonLoading, pageName, organizationId, isFirstVisit };
