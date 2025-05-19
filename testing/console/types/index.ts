/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

export type ConsoleTests = {
  login: () => Promise<void>;
  onboarding: () => Promise<void>;
  createMonthlyPlan: () => Promise<void>;
  createDailyPlan: () => Promise<void>;
  editPlan: () => Promise<void>;
  uploadSims: () => Promise<void>;
  createNetwork: () => Promise<void>;
  renameNode: () => Promise<void>;
  renameSite: () => Promise<void>;
  createSubscriber: () => Promise<void>;
  topupSubscriberData: () => Promise<void>;
  restartNode: () => Promise<void>;
  nodeRFOff: () => Promise<void>;
  nodeRFOn: () => Promise<void>;
  logout: () => Promise<void>;
};
