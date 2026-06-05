/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shown when the session is valid but the account can't be resolved to a
 * complete console identity (no ukama user for the Kratos identity, no org
 * membership, no role, or incomplete claims — see BFF SessionValidationError
 * steps). The only ways out are logging out or contacting support.
 */
import Button from '@mui/material/Button';
import LogoutRounded from '@mui/icons-material/LogoutRounded';
import SupportAgentRounded from '@mui/icons-material/SupportAgentRounded';

import { env } from '@/env';

const SUPPORT_EMAIL = 'support@ukama.com';

export default function UnauthorizedPage() {
  const logoutUrl = new URL('/auth/logout', env.NEXT_PUBLIC_AUTH_APP_URL).toString();
  return (
    <main
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'var(--uk-page)',
        padding: 24,
      }}
    >
      <div
        className="card"
        style={{ maxWidth: 480, padding: '40px 36px', textAlign: 'center' }}
      >
        <div
          style={{
            fontFamily: 'var(--font-display)',
            fontSize: 22,
            fontWeight: 500,
            marginBottom: 10,
          }}
        >
          Your account isn&apos;t set up for this console
        </div>
        <p
          style={{
            fontSize: 13.5,
            color: 'var(--uk-ink-2)',
            lineHeight: 1.6,
            margin: '0 0 24px',
            textWrap: 'pretty',
          }}
        >
          You&apos;re signed in, but we couldn&apos;t link your account to an
          organization, role, or complete profile. This usually means your user
          isn&apos;t registered with an organization yet. Please log out and
          sign in with a different account, or contact us so we can finish
          setting you up.
        </p>
        <div style={{ display: 'flex', gap: 10, justifyContent: 'center' }}>
          <Button
            variant="contained"
            startIcon={<LogoutRounded />}
            href={logoutUrl}
          >
            Log out
          </Button>
          <Button
            variant="outlined"
            startIcon={<SupportAgentRounded />}
            href={`mailto:${SUPPORT_EMAIL}?subject=Console%20access%20issue`}
          >
            Contact us
          </Button>
        </div>
      </div>
    </main>
  );
}
