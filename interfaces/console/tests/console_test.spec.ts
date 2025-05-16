/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { login } from './helpers';
import { consoleTest } from './test_base';

consoleTest.describe('Console End-to-End Test Flow', () => {
  consoleTest('1. Login Test', async ({ page }) => {
    await login(page);
  });

  consoleTest('2. Onboarding Test', async ({ page, onboarding }) => {
    await login(page);
    await onboarding();
  });

  consoleTest('3. Sim Pool Management Test', async ({ page, uploadSims }) => {
    await login(page);
    await uploadSims();
  });

  consoleTest(
    '4. Data Plan Management Test',
    async ({ page, createMonthlyPlan, createDailyPlan, editPlan }) => {
      await login(page);
      await createMonthlyPlan();
      await createDailyPlan();
      await editPlan();
    },
  );

  consoleTest(
    '5. Subscriber Management Test',
    async ({ page, createSubscriber, topupSubscriberData }) => {
      await login(page);
      await createSubscriber();
      await topupSubscriberData();
    },
  );

  consoleTest('6. Network Management Test', async ({ page, createNetwork }) => {
    await login(page);
    await createNetwork();
  });

  consoleTest('7. Logout Test', async ({ page, logout }) => {
    await login(page);
    await logout();
  });
});
