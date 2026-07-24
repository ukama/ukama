/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { notFound } from 'next/navigation';
// import BillingScreen from '../_components/BillingScreen';

// Billing is hidden for now. The route is disabled (returns 404) but the
// BillingScreen code + UI are kept for when we finish it. To re-enable:
// restore the nav item in _config/nav.ts, uncomment the import + render below,
// and delete the notFound() guard.
export default function BizBillingPage() {
  notFound();

  // return <BillingScreen />;
}
