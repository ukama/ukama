/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

import { test } from '@playwright/test';

import {
  applyClaimSimsPatch,
  applyCreateNetworkPatch,
  applyCreateSubscriberPatch,
  applyDataPlanCreationPatch,
  applyLoginPatch,
  applyLogoutPatch,
  applyNodeRfOffPatch,
  applyNodeRfOnPatch,
  applyOnboardingPatch,
  applyRenameNodePatch,
  applyRenameSitePatch,
  applyResetNodePatch,
  applyTopupSubscriberPatch,
} from '../patches';

test('apply login patch', async () => {
  await applyLoginPatch();
});

test('apply logout patch', async () => {
  await applyLogoutPatch();
});

test('apply data plan creation patch', async () => {
  await applyDataPlanCreationPatch();
});

test('apply create network patch', async () => {
  await applyCreateNetworkPatch();
});

test('apply claim sims patch', async () => {
  await applyClaimSimsPatch();
});

test('apply rename node patch', async () => {
  await applyRenameNodePatch();
});

test('apply rename site patch', async () => {
  await applyRenameSitePatch();
});

test('apply reset node patch', async () => {
  await applyResetNodePatch();
});

test('apply node rf off patch', async () => {
  await applyNodeRfOffPatch();
});

test('apply node rf on patch', async () => {
  await applyNodeRfOnPatch();
});

test('apply create subscriber patch', async () => {
  await applyCreateSubscriberPatch();
});

test('apply topup subscriber patch', async () => {
  await applyTopupSubscriberPatch();
});

test('apply onboarding patch', async () => {
  await applyOnboardingPatch();
});
