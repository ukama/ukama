/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { CurrentBillColumns } from '@/constants/tableColumns';
import { NoBillYet } from '@/public/svg';
import { RoundedCard } from '@/styles/global';
import {
  Grid,
  Stack,
  Button,
  MenuItem,
  Menu,
  Box,
  Typography,
} from '@mui/material';
import { useState } from 'react';
import CurrentBill from '../CurrentBill';
import NotificationContainer from '../NotificationContainer';
import SimpleDataTable from '../SimpleDataTable';
import TableHeader from '../TableHeader';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';

interface ICurrentBillTab {
  data: any;
  loading: boolean;
  planName: string;
  totalAmount: string;
  currentBill: string;
}

const CurrentPlanTab = () => {
  const [anchorEl, setAnchorEl] = useState(null);

  const handleClick = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <Grid container item spacing={2}>
      <Grid xs={12} item>
        <RoundedCard radius="4px">
          <Typography variant="body1" color="initial"></Typography>
          <TableHeader title={'Current plan'} showSecondaryButton={false} />
          <Grid container spacing={2} sx={{ py: 2 }}>
            <Grid item xs={6}>
              <Typography variant="body2">
                Community bundle - free Console plan for basic network
                management needs.
              </Typography>
            </Grid>
            <Grid item xs={6} container justifyContent={'flex-end'}>
              <Typography variant="h6" color="initial">
                $ 20.30
              </Typography>
            </Grid>
            <Box>
              <Button
                aria-controls="dropdown-menu"
                aria-haspopup="true"
                onClick={handleClick}
                endIcon={<ExpandMoreIcon />}
              >
                View Bundle Features
              </Button>
              <Menu
                id="dropdown-menu"
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleClose}
              >
                <MenuItem onClick={handleClose}>Feature 1</MenuItem>
                <MenuItem onClick={handleClose}>Feature 2</MenuItem>
                <MenuItem onClick={handleClose}>Feature 3</MenuItem>
              </Menu>
            </Box>
          </Grid>
        </RoundedCard>
      </Grid>
    </Grid>
  );
};

export default CurrentPlanTab;
