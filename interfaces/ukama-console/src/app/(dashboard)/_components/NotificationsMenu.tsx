/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useState } from 'react';
import Box from '@mui/material/Box';
import Menu from '@mui/material/Menu';
import Typography from '@mui/material/Typography';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import NotificationsRounded from '@mui/icons-material/NotificationsRounded';
import { useTopBarAlertsQuery } from '@/client/graphql/network-home.generated';
import { useToast } from '@/components/ToastProvider';
import type { Alert } from '@/data';
import { POLL_OVERVIEW_MS, visiblePoll } from '@/lib/polling';
import { useUiPrefs } from '@/lib/store';
import AlertDialog from './AlertDialog';
import { Ic } from './icons';

/** NotificationsDto.type → alert severity. */
const sevFromType = (type: string): Alert['sev'] => {
  const t = type.toUpperCase();
  if (t.includes('CRITICAL') || t.includes('ERROR')) return 'critical';
  if (t.includes('WARNING')) return 'warning';
  return 'info';
};

const relativeAge = (iso: string): string => {
  const ms = Date.now() - new Date(iso).getTime();
  if (!Number.isFinite(ms) || ms < 0) return '';
  const mins = Math.floor(ms / 60_000);
  if (mins < 60) return `${mins}m`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h`;
  return `${Math.floor(hours / 24)}d`;
};

const SEV_ICON: Record<Alert['sev'], string> = {
  critical: 'error',
  warning: 'warning',
  info: 'info',
};
const SEV_COLOR: Record<Alert['sev'], string> = {
  critical: 'var(--uk-error)',
  warning: 'var(--uk-orange)',
  info: 'var(--uk-ac)',
};

export default function NotificationsMenu() {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const [openAlert, setOpenAlert] = useState<Alert | null>(null);
  const toast = useToast();
  const networkId = useUiPrefs((s) => s.networkId);

  // Polled notifications (v1: no subscriptions — BUILD-PLAN §5.1·4); shares
  // the networkOverview cache entry with the home screen.
  const { data } = useTopBarAlertsQuery({
    variables: { networkId },
    skip: !networkId,
    ...visiblePoll(POLL_OVERVIEW_MS),
  });
  const alerts: Alert[] = (
    data?.networkOverview.latestAlerts.notifications ?? []
  ).map((n) => ({
    id: n.id,
    sev: sevFromType(n.type),
    icon: 'info',
    title: n.title,
    detail: n.description,
    action: 'View',
    age: relativeAge(n.createdAt),
  }));
  const unread = alerts.filter((a) => a.sev !== 'info').length;

  const runAction = (a: Alert) => {
    setOpenAlert(null);
    toast(`${a.action} — done`);
  };

  return (
    <>
      <button
        type="button"
        className="topbar-icon"
        title="Notifications"
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        <NotificationsRounded sx={{ fontSize: 22 }} />
        {unread > 0 && <span className="badge-num">{unread}</span>}
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        slotProps={{ paper: { sx: { width: 344, mt: 0.5 } } }}
      >
        <Box
          sx={{
            px: 2,
            py: 1,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <Typography sx={{ fontWeight: 600, fontSize: 13.5 }}>
            Notifications
          </Typography>
          {unread > 0 && (
            <Typography sx={{ fontSize: 12, color: 'primary.main', fontWeight: 600 }}>
              {unread} unread
            </Typography>
          )}
        </Box>
        {alerts.length === 0 && (
          <Box sx={{ px: 2, py: 2, fontSize: 13, color: 'text.secondary' }}>
            No notifications.
          </Box>
        )}
        {alerts.map((a) => (
          <Box
            key={a.id}
            component="button"
            type="button"
            onClick={() => {
              setAnchor(null);
              setOpenAlert(a);
            }}
            sx={{
              display: 'flex',
              gap: 1.5,
              width: '100%',
              textAlign: 'left',
              alignItems: 'flex-start',
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              px: 2,
              py: 1.25,
              fontFamily: 'inherit',
              '&:hover': { bgcolor: 'action.hover' },
            }}
          >
            <Ic
              name={SEV_ICON[a.sev]}
              sx={{ fontSize: 20, mt: '1px', color: SEV_COLOR[a.sev], flex: 'none' }}
            />
            <Box sx={{ flex: 1, minWidth: 0 }}>
              <Typography sx={{ fontSize: 13, fontWeight: 600 }}>
                {a.title}
              </Typography>
              <Typography
                sx={{ fontSize: 12, color: 'text.secondary', mt: '1px' }}
              >
                {a.detail}
              </Typography>
              <Typography
                sx={{ fontSize: 11.5, color: 'text.disabled', mt: '3px' }}
              >
                {a.site ? `${a.site} · ` : ''}
                {a.age} ago
              </Typography>
            </Box>
            <ChevronRightRounded
              sx={{ fontSize: 18, color: 'text.disabled', alignSelf: 'center', flex: 'none' }}
            />
          </Box>
        ))}
      </Menu>
      {openAlert && (
        <AlertDialog alert={openAlert} onClose={() => setOpenAlert(null)} onAction={runAction} />
      )}
    </>
  );
}
