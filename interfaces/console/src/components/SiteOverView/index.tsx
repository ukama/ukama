/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useState } from 'react';
import {
  Box,
  Typography,
  Button,
  Menu,
  MenuItem,
  IconButton,
  Grid,
  Paper,
  Stack,
} from '@mui/material';
import { CheckCircle } from '@mui/icons-material';
import RestartAltIcon from '@mui/icons-material/RestartAlt';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import dynamic from 'next/dynamic';

const SiteStatusIcon = React.lazy(() =>
  import('@/../public/svg').then((module) => ({
    default: module.SiteStatusIcon,
  })),
);

// Dummy data for sites
const sites = [
  { id: 'siteX', name: 'Site X' },
  { id: 'siteY', name: 'Site Y' },
  { id: 'siteZ', name: 'Site Z' },
];

interface SiteOverViewProps {
  restartSite: (site: { id: string; name: string }) => void;
  addSite: () => void;
}

const SiteOverView: React.FC<SiteOverViewProps> = ({
  restartSite,
  addSite,
}) => {
  const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
    loading: () => <p>Site map is loading</p>,
    ssr: false,
  });
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedSite, setSelectedSite] = useState(sites[0]);
  const [address, setAddress] = useState('');

  const handleMenuClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleSiteSelect = (site: { id: string; name: string }) => {
    setSelectedSite(site);
    handleMenuClose();
  };

  const handleAddressChange = (address: string) => {
    setAddress(address);
  };

  const handleRestartClick = () => {
    restartSite(selectedSite);
  };

  const handleAddSiteClick = () => {
    addSite();
  };

  return (
    <>
      <Grid item xs={6}>
        <Box display="flex" alignItems="center" gap={1}>
          <CheckCircle color="success" />
          <Typography variant="body1">{selectedSite.name}</Typography>
          <IconButton onClick={handleMenuClick}>
            <ArrowDropDownIcon />
          </IconButton>
          <Typography variant="body1" color="initial">
            is online for 3 days
          </Typography>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
          >
            {sites.map((site) => (
              <MenuItem key={site.id} onClick={() => handleSiteSelect(site)}>
                {site.name}
              </MenuItem>
            ))}
            <MenuItem onClick={handleAddSiteClick}>Add site</MenuItem>
          </Menu>
        </Box>
      </Grid>
      <Grid
        item
        xs={6}
        container
        justifyItems={'center'}
        justifyContent={'flex-end'}
      >
        <Button
          variant="contained"
          startIcon={<RestartAltIcon />}
          color="primary"
          onClick={handleRestartClick}
        >
          Restart Site
        </Button>
      </Grid>
      <Grid item xs={3}>
        <Paper sx={{ p: 2 }}>
          <Stack direction="column" spacing={1}>
            <Typography variant="body1" sx={{ fontWeight: 'bold' }}>
              Site information
            </Typography>
            <Stack direction="row" spacing={2}>
              <Typography variant="subtitle2">
                Date created: July 13 2023
              </Typography>
            </Stack>
            <Typography variant="subtitle2">Location : {address}</Typography>
            <Typography variant="subtitle2">Lat : -23783</Typography>
            <Typography variant="subtitle2">Long : 783</Typography>
            <Typography variant="subtitle2">Nodes : 783</Typography>
          </Stack>
        </Paper>
      </Grid>
      <Grid item xs={6}>
        <Paper sx={{ p: 2, height: '230px' }}>
          <Typography variant="body1" sx={{ fontWeight: 'bold', mb: 2 }}>
            Site overview
          </Typography>
          <Stack direction="row" spacing={2} alignItems={'center'}>
            <Box display="flex" alignItems="center">
              <SiteStatusIcon />
              <Typography variant="caption" component="span" marginLeft={1}>
                Input power
              </Typography>
            </Box>
            <Box display="flex" alignItems="center">
              <SiteStatusIcon />
              <Typography variant="caption" component="span" marginLeft={1}>
                Storage
              </Typography>
            </Box>
            <Box display="flex" alignItems="center">
              <SiteStatusIcon />
              <Typography variant="caption" component="span" marginLeft={1}>
                Consumption
              </Typography>
            </Box>
          </Stack>
        </Paper>
      </Grid>
      <Grid
        item
        xs={3}
        container
        justifyItems={'center'}
        justifyContent={'flex-end'}
      >
        <SiteMapComponent
          posix={[-2.4906, 28.8428]}
          onAddressChange={handleAddressChange}
        />
      </Grid>
    </>
  );
};

export default SiteOverView;
