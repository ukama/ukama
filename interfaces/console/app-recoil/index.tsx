/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
