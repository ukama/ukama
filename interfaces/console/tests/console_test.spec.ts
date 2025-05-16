/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { consoleTest } from './test_base';

consoleTest.describe('Console End-to-End Test Flow', () => {
  consoleTest('1. Login Test', async ({ page, login }) => {
    await login();
  });

  consoleTest('2. Onboarding Test', async ({ page, onboarding, login }) => {
    await login();
    await onboarding();
  });

  consoleTest(
    '3. Sim Pool Management Test',
    async ({ page, uploadSims, login }) => {
      await login();
      await uploadSims();
    },
  );

  consoleTest(
    '4. Data Plan Management Test',
    async ({ page, createMonthlyPlan, createDailyPlan, editPlan, login }) => {
      await login();
      await createMonthlyPlan();
      await createDailyPlan();
      await editPlan();
    },
  );

  consoleTest(
    '5. Subscriber Management Test',
    async ({ page, createSubscriber, topupSubscriberData, login }) => {
      await login();
      await createSubscriber();
      await topupSubscriberData();
    },
  );

  consoleTest(
    '6. Network Management Test',
    async ({ page, createNetwork, login }) => {
      await login();
      await createNetwork();
    },
  );

  consoleTest('7. Rename Site Test', async ({ page, renameSite, login }) => {
    await login();
    await renameSite();
  });

  consoleTest('8. Rename Node Test', async ({ page, renameNode, login }) => {
    await login();
    await renameNode();
  });

  consoleTest('9. Logout Test', async ({ page, logout, login }) => {
    await login();
    await logout();
  });
});
