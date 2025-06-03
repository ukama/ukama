/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Role_Type } from '@/client/graphql/generated';
import { MEMBER_ROLES } from '@/constants';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  Grid,
  IconButton,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
} from '@mui/material';
import { useState } from 'react';

type InviteMemberDialogProps = {
  title: string;
  isOpen: boolean;
  labelSuccessBtn?: string;
  labelNegativeBtn?: string;
  handleCloseAction: () => void;
  handleSuccessAction?: (member: any) => void;
  invitationLoading?: boolean;
};

const InviteMemberDialog = ({
  title,
  isOpen,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
  invitationLoading,
}: InviteMemberDialogProps) => {
  const [member, setMember] = useState({
    role: '',
    email: '',
    name: '',
  });
  return (
    <Dialog
      fullWidth
      open={isOpen}
      maxWidth="sm"
      onClose={() => handleCloseAction()}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <DialogTitle>{title}</DialogTitle>
        <IconButton
          onClick={() => handleCloseAction()}
          sx={{ position: 'relative', right: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </Stack>

      <DialogContent>
        <Grid
          container
          rowSpacing={2}
          gridAutoRows={2}
          gridAutoColumns={1}
          alignItems={'center'}
          justifyContent={'center'}
          spacing={2}
        >
          <Grid item xs={6}>
            <TextField
              fullWidth
              required
              label="Name"
              value={member.name}
              id={'invite-member-name'}
              InputLabelProps={{
                shrink: true,
              }}
              onChange={(e) => setMember({ ...member, name: e.target.value })}
            />
          </Grid>
          <Grid item xs={6}>
            <TextField
              fullWidth
              required
              label="Email"
              value={member.email}
              id={'invite-member-email'}
              InputLabelProps={{
                shrink: true,
              }}
              onChange={(e) => setMember({ ...member, email: e.target.value })}
            />
          </Grid>

          <Grid item xs={12}>
            <FormControl fullWidth>
              <InputLabel id={'invite-member-role-label'} shrink>
                Role
              </InputLabel>
              <Select
                notched
                required
                label="Role"
                value={member.role}
                id={'invite-member-role'}
                labelId="invite-member-role-label"
                onChange={(e) =>
                  setMember({ ...member, role: e.target.value as Role_Type })
                }
              >
                {MEMBER_ROLES.map(({ id, label, value }) => (
                  <MenuItem key={id} value={value}>
                    {label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
        </Grid>
      </DialogContent>

      <DialogActions>
        <Stack direction={'row'} alignItems="center" spacing={2}>
          {labelNegativeBtn && (
            <Button
              variant="text"
              color={'primary'}
              onClick={() => handleCloseAction()}
            >
              {labelNegativeBtn}
            </Button>
          )}
          {labelSuccessBtn && (
            <Button
              disabled={
                !member.role ||
                !member.email ||
                !member.name ||
                invitationLoading
              }
              variant="contained"
              onClick={() =>
                handleSuccessAction
                  ? handleSuccessAction(member)
                  : handleCloseAction()
              }
            >
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default InviteMemberDialog;
