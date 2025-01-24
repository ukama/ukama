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
  IconButton,
  Menu,
  MenuItem,
  Typography,
  Skeleton,
  Button,
  Grid,
  useMediaQuery,
} from '@mui/material';
import { CheckCircle, Add, ArrowDropDown } from '@mui/icons-material';
import { SiteDto } from '@/client/graphql/generated';

interface SiteDetailsHeaderProps {
  siteList: SiteDto[];
  selectedSiteId: string | null;
  onSiteChange: (siteId: string) => void;
  isLoading: boolean;
  onRestartSite: () => void;
}

const SiteDetailsHeader: React.FC<SiteDetailsHeaderProps> = ({
  siteList,
  selectedSiteId,
  onSiteChange,
  isLoading,
  onRestartSite,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const isMobile = useMediaQuery('(max-width:600px)');

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
    <Grid container spacing={2} justifyItems={'center'} sx={{ mb: 1 }}>
      <Grid item xs={12} md={6}>
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
                  <CheckCircle color="success" />
                  <Typography variant="body1" fontWeight="medium">
                    {selectedSite.name}
                  </Typography>
                </>
              ) : (
                <Typography variant="body1" fontWeight="medium">
                  Select a site
                </Typography>
              )}
              <IconButton onClick={handleMenuClick} size="small">
                <ArrowDropDown />
              </IconButton>
            </>
          )}
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
          >
            {isLoading
              ? [1, 2, 3].map((key) => (
                  <MenuItem key={key}>
                    <Skeleton variant="rectangular" height={30} width={200} />
                  </MenuItem>
                ))
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
        </Box>
      </Grid>

      {!isMobile && (
        <Grid item xs={6} justifyContent="flex-end" container sx={{ mt: 1 }}>
          <Button
            variant="outlined"
            onClick={onRestartSite}
            disabled={!selectedSiteId}
          >
            Restart Site
          </Button>
        </Grid>
      )}
    </Grid>
  );
};

export default SiteDetailsHeader;
