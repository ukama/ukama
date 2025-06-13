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
  applyBackhaulPortTogglePatch,
  applySolarPortTogglePatch,
  applyNodePortTogglePatch,
} from '../patches/v1.0.0';

test('v1.0.0 - Apply login patch', async () => {
  await applyLoginPatch();
});

test('v1.0.0 - Apply node port toggle patch', async () => {
  await applyNodePortTogglePatch();
});
test('v1.0.0 - Apply backhaul port toggle patch', async () => {
  await applyBackhaulPortTogglePatch();
});

test('v1.0.0 - Apply solar port toggle patch', async () => {
  await applySolarPortTogglePatch();
});
test('v1.0.0 - Apply logout patch', async () => {
  await applyLogoutPatch();
});

test('v1.0.0 - Apply data plan creation patch', async () => {
  await applyDataPlanCreationPatch();
});

test('v1.0.0 - Apply create network patch', async () => {
  await applyCreateNetworkPatch();
});

test('v1.0.0 - Apply claim sims patch', async () => {
  await applyClaimSimsPatch();
});

test('v1.0.0 - Apply rename node patch', async () => {
  await applyRenameNodePatch();
});

test('v1.0.0 - Apply rename site patch', async () => {
  await applyRenameSitePatch();
});

test('v1.0.0 - Apply reset node patch', async () => {
  await applyResetNodePatch();
});

test('v1.0.0 - Apply node rf off patch', async () => {
  await applyNodeRfOffPatch();
});

test('v1.0.0 - Apply node rf on patch', async () => {
  await applyNodeRfOnPatch();
});

test('v1.0.0 - Apply create subscriber patch', async () => {
  await applyCreateSubscriberPatch();
});

test('v1.0.0 - Apply topup subscriber patch', async () => {
  await applyTopupSubscriberPatch();
});

test('v1.0.0 - Apply onboarding patch', async () => {
  await applyOnboardingPatch();
});
