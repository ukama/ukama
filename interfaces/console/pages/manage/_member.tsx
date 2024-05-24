/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  INVITATION_TABLE_COLUMN,
  MEMBER_TABLE_COLUMN,
  MEMBER_TABLE_MENU,
} from '@/constants';
import { colors } from '@/styles/theme';
import DataTableWithOptions from '@/ui/molecules/DataTableWithOptions';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import { Search } from '@mui/icons-material';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import {
  Box,
  Button,
  Grid,
  Paper,
  Stack,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@mui/material';
import React, { useState } from 'react';

interface IMember {
  search: string;
  memberData: any;
  invitationsData: any;
  invitationTitle: string;
  handleButtonAction: () => void;
  setSearch: (value: string) => void;
  handleMemberAction: (id: string, type: string) => void;
  handleDeleteInviteAction: (uuid: string) => void;
}

const Member: React.FC<IMember> = ({
  search,
  setSearch,
  memberData,
  invitationsData,
  handleButtonAction,
  handleMemberAction,
  handleDeleteInviteAction,
}) => {
  const [tabIndex, setTabIndex] = useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabIndex(newValue);
  };

  const renderMemberDataTable = () => (
    <Stack direction={'column'}>
      <Typography variant="h6" fontWeight={500}>
        Members
      </Typography>
      <DataTableWithOptions
        dataset={memberData || []}
        icon={PeopleAltIcon}
        isRowClickable={false}
        columns={MEMBER_TABLE_COLUMN}
        menuOptions={MEMBER_TABLE_MENU}
        emptyViewLabel={'No members yet!'}
        onMenuItemClick={handleMemberAction}
      />
    </Stack>
  );

  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        overflow: 'scroll',
        borderRadius: '10px',
        height: 'calc(100vh - 200px)',
      }}
    >
      <Tabs value={tabIndex} onChange={handleTabChange}>
        <Tab label="team members" />
      </Tabs>
      {tabIndex === 0 && (
        <Box sx={{ width: '100%', mt: 4 }}>
          <Grid container spacing={2}>
            <Grid item xs={6}>
              <TextField
                id="subscriber-search"
                variant="outlined"
                size="small"
                placeholder="Search"
                defaultValue={search}
                fullWidth
                onChange={(e) => setSearch(e.target.value)}
                InputLabelProps={{
                  shrink: false,
                }}
                InputProps={{
                  endAdornment: <Search htmlColor={colors.black54} />,
                }}
              />
            </Grid>
            <Grid item xs={6} container justifyContent="flex-end">
              <Button
                variant="contained"
                color="primary"
                fullWidth
                sx={{ width: { xs: '100%', md: 'fit-content' } }}
                onClick={handleButtonAction}
              >
                INVITE MEMBER
              </Button>
            </Grid>
          </Grid>

          <br />
          {renderMemberDataTable()}
          <br />
          <br />
          {invitationsData.length > 0 && (
            <Stack direction={'column'}>
              <Typography variant="h6" fontWeight={500}>
                Pending/Declined Invitations
              </Typography>
              <SimpleDataTable
                dataKey="id"
                dataset={invitationsData}
                columns={INVITATION_TABLE_COLUMN}
                handleDeleteElement={(id: string) =>
                  handleDeleteInviteAction(id)
                }
              />
            </Stack>
          )}
        </Box>
      )}
    </Paper>
  );
};

export default Member;
