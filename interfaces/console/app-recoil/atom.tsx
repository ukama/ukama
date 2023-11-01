/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
    orgId: '',
    userId: '',
    orgName: '',
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
  commonData,
  isDarkmode,
  isFirstVisit,
  isSkeltonLoading,
  pageName,
  snackbarMessage,
  user,
};
