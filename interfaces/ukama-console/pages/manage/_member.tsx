import { MANAGE_TABLE_COLUMN } from '@/constants';
import EmptyView from '@/ui/molecules/EmptyView';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import { Paper, Tabs, Button, TextField, Grid, Box } from '@mui/material';
import { Search } from '@mui/icons-material';
import { colors } from '@/styles/theme';
import React, { useState } from 'react';

interface IMember {
  data: any;
  search: string;
  setSearch: (value: string) => void;
  handleButtonAction: () => void;
  invitationTitle: string;
  onSearchChange?: Function;
}

const Member = ({
  data,
  search,
  setSearch,
  handleButtonAction,
  invitationTitle,
  onSearchChange,
}: IMember) => {
  const [tabIndex, setTabIndex] = useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabIndex(newValue);
  };
  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        overflow: 'hidden',
        borderRadius: '5px',
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
                value={search}
                fullWidth
                onChange={(e) =>
                  onSearchChange && onSearchChange(e.target.value)
                }
                InputLabelProps={{
                  shrink: false,
                }}
                InputProps={{
                  endAdornment: <Search htmlColor={colors.black54} />,
                }}
              />
            </Grid>
            <Grid item xs={6} container justifyContent={'flex-end'}>
              <Button
                variant="contained"
                color="primary"
                fullWidth
                sx={{ width: { xs: '100%', md: 'fit-content' } }}
                onClick={() => handleButtonAction()}
              >
                {`INVITE MEMBER`}
              </Button>
            </Grid>
          </Grid>

          <br />
          {data && data.length > 0 ? (
            <SimpleDataTable
              dataKey="uuid"
              dataset={data && data}
              columns={MANAGE_TABLE_COLUMN}
            />
          ) : (
            <EmptyView icon={PeopleAltIcon} title="No members yet!" />
          )}
        </Box>
      )}
    </Paper>
  );
};

export default Member;
