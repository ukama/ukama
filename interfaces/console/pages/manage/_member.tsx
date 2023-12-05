/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { MANAGE_TABLE_COLUMN } from '@/constants';
import { colors } from '@/styles/theme';
import EmptyView from '@/ui/molecules/EmptyView';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import { Search } from '@mui/icons-material';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import {
  Box,
  Button,
  Grid,
  Paper,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@mui/material';
import React, { useState } from 'react';

interface IMember {
  memberData: any[];
  invitationsData: any;
  search: string;
  setSearch: (value: string) => void;
  handleButtonAction: () => void;
  invitationTitle: string;
  onSearchChange?: (value: string) => void;
}

const Member: React.FC<IMember> = ({
  memberData,
  invitationsData,
  search,
  setSearch,
  handleButtonAction,
  onSearchChange,
}) => {
  const [tabIndex, setTabIndex] = useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabIndex(newValue);
  };

  const renderMemberDataTable = () => {
    if (memberData && memberData.length > 0) {
      return (
        <>
          <SimpleDataTable
            dataKey="uuid"
            dataset={memberData}
            columns={MANAGE_TABLE_COLUMN}
          />
        </>
      );
    } else {
      return (
        <Box
          sx={{
            width: '100%',
            mt: 20,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <EmptyView icon={PeopleAltIcon} title="No members yet!" />
        </Box>
      );
    }
  };

  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        overflow: 'hidden',
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
          {invitationsData && invitationsData.invitations && (
            <>
              <Typography variant="body1">Pending invitations</Typography>
              <SimpleDataTable
                dataKey="uuid"
                dataset={invitationsData.invitations}
                columns={MANAGE_TABLE_COLUMN}
              />
            </>
          )}
        </Box>
      )}
    </Paper>
  );
};

export default Member;
