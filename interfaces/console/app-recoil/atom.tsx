import { atom } from 'recoil';
import { recoilPersist } from 'recoil-persist';
const { persistAtom } = recoilPersist();

const isFirstVisit = atom({
  key: 'isFirstVisit',
  default: true,
  effects_UNSTABLE: [persistAtom],
});

const isSkeltonLoading = atom({
  key: 'isSkeltonLoading',
  default: false,
  effects_UNSTABLE: [persistAtom],
});

const isDarkmode = atom({
  key: 'isDarkmode',
  default: false,
  effects_UNSTABLE: [persistAtom],
});

const pageName = atom({
  key: 'pageName',
  default: 'Home',
  effects_UNSTABLE: [persistAtom],
});
const commonData = atom({
  key: 'commonData',
  default: {
    networkId: '',
    networkName: '',
    orgId: 'aac1ed88-2546-4f9c-a808-fb9c4d0ef24b',
    userId: '018688fa-d861-4e7b-b119-ffc5e1637ba8',
    orgName: 'ukama',
  },
  effects_UNSTABLE: [persistAtom],
});
const snackbarMessage = atom({
  key: 'snackbarMessage',
  default: { id: 'message-id', message: '', type: 'info', show: false },
});

const user = atom({
  key: 'user',
  default: { id: '', name: '', email: '', role: '', isFirstVisit: false },
  effects_UNSTABLE: [persistAtom],
});

export {
  commonData, isDarkmode, isFirstVisit, isSkeltonLoading, pageName, snackbarMessage, user
};

