/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * First-visit welcome page. proxy.ts routes users here when the session
 * token says isShowWelcome=true and routes them away once acknowledged, so
 * this page only renders for eligible users (redirect() is a safety net).
 */
import { redirect } from 'next/navigation';
import { getCurrentUser } from '@/lib/auth/server';
import WelcomeScreen from './_components/WelcomeScreen';
import './welcome.css';

export const metadata = { title: 'Welcome · Ukama Console' };

export default async function WelcomePage() {
  const user = await getCurrentUser();
  if (!user) redirect('/unauthorized');
  if (!user.isShowWelcome) redirect('/');
  return <WelcomeScreen user={user} />;
}
