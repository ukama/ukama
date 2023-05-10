import { selector } from 'recoil';
import {
  isDarkmode,
  isFirstVisit,
  isSkeltonLoading,
  networkName,
  pageName,
  snackbarMessage,
  user,
} from './atom';

const withIsSkeltonLoading = selector({
  key: 'withIsSkeltonLoading',
  get: ({ get }) => get(isSkeltonLoading),
});

const withPageName = selector({
  key: 'withPageName',
  get: ({ get }) => get(pageName),
});

const withIsFirstVisit = selector({
  key: 'withIsFirstVisit',
  get: ({ get }) => get(isFirstVisit),
});
const withNetworkName = selector({
  key: 'withNetworkName',
  get: ({ get }) => get(networkName),
});
const withIsDarkMod = selector({
  key: 'withIsDarkMod',
  get: ({ get }) => get(isDarkmode),
});

const withSnackbarMessage = selector({
  key: 'withSnackbarMessage',
  get: ({ get }) => get(snackbarMessage),
});
const withUser = selector({
  key: 'withUser',
  get: ({ get }) => get(user),
});

export {
  withUser,
  withNetworkName,
  withPageName,
  withIsDarkMod,
  withIsFirstVisit,
  withSnackbarMessage,
  withIsSkeltonLoading,
};
