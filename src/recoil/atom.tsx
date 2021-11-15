import { atom } from "recoil";
import { recoilPersist } from "recoil-persist";
const { persistAtom } = recoilPersist();

const isLoginAtom = atom({
    key: "isLoginAtom",
    default: false,
    effects_UNSTABLE: [persistAtom],
});

const isSkeltonLoading = atom({
    key: "isSkeltonLoading",
    default: false,
    effects_UNSTABLE: [persistAtom],
});

export { isLoginAtom, isSkeltonLoading };
