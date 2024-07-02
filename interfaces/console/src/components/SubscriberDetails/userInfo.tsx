/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import LoadingWrapper from '@/components/LoadingWrapper';
import { colors } from '@/theme';
import EditIcon from '@mui/icons-material/Edit';
import { IconButton, Stack, TextField, Typography } from '@mui/material';
import React from 'react';

interface UserInfoProps {
  subscriberLoading: boolean;
  onEditName: boolean;
  firstName: string;
  handleEditName: () => void;
  onEditEmail: boolean;
  email: string;
  handleSimEdit: () => void;
  setOnEditName: any;
  setOnEditEmail: any;
}
const UserInfo: React.FC<UserInfoProps> = ({
  subscriberLoading,
  onEditName,
  firstName,
  handleEditName,
  onEditEmail,
  email,
  handleSimEdit,
  setOnEditName,
  setOnEditEmail,
}) => (
  <LoadingWrapper
    radius="small"
    width={'100%'}
    isLoading={subscriberLoading}
    cstyle={{
      overflow: 'auto',
      backgroundColor: false ? colors.white : 'transparent',
    }}
  >
    <Stack direction="column" spacing={2}>
      <Typography variant="body1" sx={{ color: colors.black }}>
        Name
      </Typography>
      <Stack
        direction="row"
        spacing={2}
        alignItems={'center'}
        justifyContent={'space-between'}
        sx={{ pr: 2 }}
      >
        <TextField
          id="outlined-basic"
          value={firstName}
          variant="standard"
          disabled={!onEditName}
          size="small"
          onChange={handleEditName}
          sx={{ width: '100%' }}
        />

        <IconButton
          size="small"
          color="primary"
          onClick={() => setOnEditName(!onEditName)}
        >
          <EditIcon fontSize="small" />
        </IconButton>
      </Stack>

      <Typography variant="body1" sx={{ color: colors.black }}>
        Email
      </Typography>
      <Stack
        direction="row"
        spacing={2}
        alignItems={'center'}
        justifyContent={'space-between'}
        sx={{ pr: 2 }}
      >
        <TextField
          id="outlined-basic"
          value={email}
          variant="standard"
          disabled={!onEditEmail}
          size="small"
          onChange={handleSimEdit}
          sx={{ width: '100%' }}
        />

        <IconButton
          size="small"
          color="primary"
          onClick={() => setOnEditEmail(!onEditEmail)}
        >
          <EditIcon fontSize="small" />
        </IconButton>
      </Stack>
    </Stack>
  </LoadingWrapper>
);

export default UserInfo;
