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
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import BadgeRounded from '@mui/icons-material/BadgeRounded';
import MailRounded from '@mui/icons-material/MailRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import PersonAddRounded from '@mui/icons-material/PersonAddRounded';
import PersonRemoveRounded from '@mui/icons-material/PersonRemoveRounded';
import { useTeamListQuery } from '@/client/graphql/team.generated';
import { EmptyState } from '@/components/EmptyState';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { ROLE_DESC } from '@/data';
import InviteMemberDialog from './InviteMemberDialog';

/** Team row view-model from the membersView composite (members +
 *  invitations merged server-side; status: Active | Deactivated | Invited). */
interface TeamRow {
  id: string;
  name?: string | null;
  email?: string | null;
  role: string;
  status: string;
  memberSince?: string | null;
}

function MemberMenu({ m }: { m: TeamRow }) {
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
  const [showInvite, setShowInvite] = useState(false);
  const { data, loading, refetch } = useTeamListQuery();
  const teamSection = data?.membersView.team;
  const rows: TeamRow[] = teamSection?.rows ?? [];

  return (
    <div className="page">
      <PageHeader
        title="Members"
        count={rows.length}
        sub="People with access to this organization."
        actions={
          <Button
            variant="contained"
            startIcon={<PersonAddRounded />}
            onClick={() => setShowInvite(true)}
          >
            Invite member
          </Button>
        }
      />
      <div className="card card-pad">
        <div className="tbl-wrap">
          {loading ? (
            <SkeletonTable cols={5} rows={5} lead />
          ) : teamSection?.error ? (
            <EmptyState
              art="error"
              title="Couldn't load members"
              sub={teamSection.error.message}
              cta="Try again"
              onCta={() => refetch()}
            />
          ) : rows.length === 0 ? (
            <EmptyState
              art="people"
              title="No members yet"
              sub="Invite teammates to give them access."
            />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Member</TableCell>
                  <TableCell>Role</TableCell>
                  <TableCell>Last active</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell sx={{ width: 44 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {rows.map((m) => (
                  <TableRow key={m.id}>
                    <TableCell>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 11 }}>
                        <span className="av-sm">
                          {(m.name ?? m.email ?? '?')
                            .split(' ')
                            .map((x) => x[0])
                            .join('')}
                        </span>
                        <div>
                          <div style={{ fontWeight: 600 }}>{m.name ?? '—'}</div>
                          <div className="muted" style={{ fontSize: 12 }}>
                            {m.email ?? '—'}
                          </div>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div style={{ fontWeight: 600 }}>{m.role}</div>
                      <div className="muted" style={{ fontSize: 12 }}>
                        {(ROLE_DESC as Record<string, string>)[m.role] ?? ''}
                      </div>
                    </TableCell>
                    <TableCell className="muted">{m.memberSince ?? '—'}</TableCell>
                    <TableCell>
                      <StatusBadge
                        status={
                          m.status === 'Active'
                            ? 'active'
                            : m.status === 'Invited'
                              ? 'pending'
                              : 'inactive'
                        }
                      >
                        {m.status}
                      </StatusBadge>
                    </TableCell>
                    <TableCell>
                      <MemberMenu m={m} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
        {!loading && !teamSection?.error && (
          <TableFooter count={rows.length} noun="members" />
        )}
      </div>
      {showInvite && <InviteMemberDialog onClose={() => setShowInvite(false)} />}
    </div>
  );
}
