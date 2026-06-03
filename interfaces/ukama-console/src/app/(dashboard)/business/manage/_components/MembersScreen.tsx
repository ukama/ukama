/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Members — people with access to this organization (screens-manage.jsx). */
import { useState } from 'react';
import Button from '@mui/material/Button';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import BadgeRounded from '@mui/icons-material/BadgeRounded';
import MailRounded from '@mui/icons-material/MailRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import PersonAddRounded from '@mui/icons-material/PersonAddRounded';
import PersonRemoveRounded from '@mui/icons-material/PersonRemoveRounded';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { MEMBERS, ROLE_DESC } from '@/data';
import type { Member } from '@/data';
import { useFirstLoad } from '@/lib/useFirstLoad';

function MemberMenu({ m }: { m: Member }) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const toast = useToast();
  return (
    <>
      <IconButton
        size="small"
        aria-label="More actions"
        sx={{ color: 'var(--uk-ink-3)' }}
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        <MoreVertRounded sx={{ fontSize: 20 }} />
      </IconButton>
      <Menu anchorEl={anchor} open={!!anchor} onClose={() => setAnchor(null)}>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Change role for ${m.name}`);
          }}
        >
          <BadgeRounded sx={{ fontSize: 18 }} /> Change role
        </MenuItem>
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25 }}
          onClick={() => {
            setAnchor(null);
            toast(`Invite resent to ${m.email}`);
          }}
        >
          <MailRounded sx={{ fontSize: 18 }} /> Resend invite
        </MenuItem>
        <Divider />
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25, color: 'var(--uk-error)' }}
          onClick={() => {
            setAnchor(null);
            toast(`${m.name} removed`, {
              action: { label: 'Undo', fn: () => toast(`${m.name} restored`) },
            });
          }}
        >
          <PersonRemoveRounded sx={{ fontSize: 18 }} /> Remove
        </MenuItem>
      </Menu>
    </>
  );
}

export default function MembersScreen() {
  const loading = useFirstLoad('members');
  const toast = useToast();

  return (
    <div className="page">
      <PageHeader
        title="Members"
        count={MEMBERS.length}
        sub="People with access to this organization."
        actions={
          <Button
            variant="contained"
            startIcon={<PersonAddRounded />}
            onClick={() => toast('Invite member — dialog lands in the overlays phase')}
          >
            Invite member
          </Button>
        }
      />
      <div className="card card-pad">
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={5} rows={5} lead />
          ) : (
            <table className="tbl">
              <thead>
                <tr className="static">
                  <th>Member</th>
                  <th>Role</th>
                  <th>Last active</th>
                  <th>Status</th>
                  <th style={{ width: 40 }} />
                </tr>
              </thead>
              <tbody>
                {MEMBERS.map((m) => (
                  <tr key={m.id} className="static">
                    <td>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 11 }}>
                        <span className="av-sm">
                          {m.name
                            .split(' ')
                            .map((x) => x[0])
                            .join('')}
                        </span>
                        <div>
                          <div style={{ fontWeight: 600 }}>{m.name}</div>
                          <div className="muted" style={{ fontSize: 12 }}>
                            {m.email}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div style={{ fontWeight: 600 }}>{m.role}</div>
                      <div className="muted" style={{ fontSize: 12 }}>
                        {ROLE_DESC[m.role]}
                      </div>
                    </td>
                    <td className="muted">{m.last}</td>
                    <td>
                      <StatusBadge status={m.status === 'active' ? 'active' : 'pending'}>
                        {m.status === 'active' ? 'Active' : 'Pending'}
                      </StatusBadge>
                    </td>
                    <td>
                      <MemberMenu m={m} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {!loading && <TableFooter count={MEMBERS.length} noun="members" />}
      </div>
    </div>
  );
}
