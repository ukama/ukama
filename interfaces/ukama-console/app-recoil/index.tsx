import {
  commonData,
  isDarkmode,
  isFirstVisit,
  isSkeltonLoading,
  pageName,
  snackbarMessage,
  user,
} from './atom';
import {
  withCommonData,
  withIsDarkMod,
  withIsFirstVisit,
  withIsSkeltonLoading,
  withPageName,
  withSnackbarMessage,
  withUser,
} from './withGlobalSelectors';

export {
  user,
  withUser,
  pageName,
  isDarkmode,
  commonData,
  withPageName,
  isFirstVisit,
  withIsDarkMod,
  withCommonData,
  snackbarMessage,
  isSkeltonLoading,
  withIsFirstVisit,
  withSnackbarMessage,
  withIsSkeltonLoading,
};
