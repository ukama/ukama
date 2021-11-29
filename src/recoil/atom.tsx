import { atom } from "recoil";
import { recoilPersist } from "recoil-persist";
const { persistAtom } = recoilPersist();

const isSkeltonLoading = atom({
    key: "isSkeltonLoading",
    default: false,
    effects_UNSTABLE: [persistAtom],
});

export { isSkeltonLoading };
