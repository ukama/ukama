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
  Grid,
  IconButton,
  Menu,
  MenuItem,
  Typography,
  Skeleton,
  Stack,
} from '@mui/material';
import { CheckCircle } from '@mui/icons-material';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { SiteDto } from '@/client/graphql/generated';
import { duration } from '@/utils';
import CancelIcon from '@mui/icons-material/Cancel';
interface SiteDetailsHeaderProps {
  siteList: SiteDto[];
  selectedSiteId: string | null;
  onSiteChange: (siteId: string) => void;
  isLoading: boolean;
  siteUpTime: number;
}

const SiteDetailsHeader: React.FC<SiteDetailsHeaderProps> = ({
  siteList,
  selectedSiteId,
  onSiteChange,
  isLoading,
  siteUpTime,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenuClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleSiteSelect = (siteId: string) => {
    onSiteChange(siteId);
    handleMenuClose();
  };

  const selectedSite =
    siteList.find((site) => site.id === selectedSiteId) || null;

  return (
    <Grid item xs={6}>
      <Box display="flex" alignItems="center" gap={1}>
        {isLoading ? (
          <>
            <Skeleton variant="circular" width={24} height={24} />
            <Skeleton variant="text" width={100} height={24} />
          </>
        ) : (
          <>
            {selectedSite ? (
              <>
                {siteUpTime > 0 ? (
                  <CheckCircle color="success" />
                ) : (
                  <CancelIcon color="error" />
                )}
                <Typography variant="body1">{selectedSite.name}</Typography>
              </>
            ) : (
              <Typography variant="body1">Select a site</Typography>
            )}
            <IconButton onClick={handleMenuClick}>
              <ArrowDropDownIcon />
            </IconButton>
          </>
        )}
        <Stack direction="row" spacing={1}>
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
          >
            {isLoading
              ? [
                  <MenuItem key="loading-1">
                    <Skeleton variant="rectangular" height={30} width={200} />
                  </MenuItem>,
                  <MenuItem key="loading-2">
                    <Skeleton variant="rectangular" height={30} width={200} />
                  </MenuItem>,
                  <MenuItem key="loading-3">
                    <Skeleton variant="rectangular" height={30} width={200} />
                  </MenuItem>,
                ]
              : [
                  ...siteList.map((site) => (
                    <MenuItem
                      key={site.id}
                      onClick={() => handleSiteSelect(site.id)}
                      selected={selectedSiteId === site.id}
                    >
                      {site.name}
                    </MenuItem>
                  )),
                ]}
          </Menu>
        </Stack>
        {isLoading || siteUpTime == null ? (
          <Skeleton variant="text" width="60" height={20} sx={{ mb: 2 }} />
        ) : siteUpTime === 0 ? (
          <Typography variant="body1">Site is currently down</Typography>
        ) : (
          <Typography variant="body1">
            Site is up for <b>{duration(siteUpTime)}</b>
          </Typography>
        )}
      </Box>
    </Grid>
  );
};

export default SiteDetailsHeader;
