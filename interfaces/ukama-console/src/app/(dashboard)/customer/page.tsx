/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { redirect } from 'next/navigation';

/** The agent lens lands on Customers (BUILD-PLAN §2). */
export default function CustomerRootPage() {
  redirect('/customer/customers');
}
