import { selector } from 'recoil';
import {
  commonData,
  isDarkmode,
  isFirstVisit,
  isSkeltonLoading,
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
const withCommonData = selector({
  key: 'withCommonData',
  get: ({ get }) => get(commonData),
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
  withPageName,
  withIsDarkMod,
  withCommonData,
  withIsFirstVisit,
  withSnackbarMessage,
  withIsSkeltonLoading,
};
