/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
import { mergeTests } from '@playwright/test';
import { appTest } from './app/app_actions';
import { authTest } from './auth/auth_actions';
import { manageTest } from './manage/manage_actions';
import { onboardingTest } from './onboarding/onboarding_actions';

export const consoleTest = mergeTests(
  authTest,
  appTest,
  manageTest,
  onboardingTest,
);

export { expect } from '@playwright/test';
