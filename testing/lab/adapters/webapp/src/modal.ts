/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { Page } from 'playwright';

// MUI also gives the customer drawer role="dialog". It stays open after a
// nested form is submitted or cancelled.
export function visibleModal(page: Page) {
  return page.locator('[role="dialog"]:not(.MuiDrawer-paper)').filter({ visible: true });
}
