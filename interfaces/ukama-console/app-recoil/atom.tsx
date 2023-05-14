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
const networkId = atom({
  key: 'networkId',
  default: '',
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
  user,
  pageName,
  networkId,
  isDarkmode,
  isFirstVisit,
  snackbarMessage,
  isSkeltonLoading,
};
