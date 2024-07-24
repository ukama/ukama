import React, { useState } from 'react';
import {
  Box,
  Grid,
  IconButton,
  Menu,
  MenuItem,
  Typography,
  Skeleton,
} from '@mui/material';
import { CheckCircle } from '@mui/icons-material';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { SiteDto } from '@/client/graphql/generated';

interface SiteDetailsHeaderProps {
  addSite: () => void;
  siteList: SiteDto[];
  selectedSiteId: string | null;
  onSiteChange: (siteId: string) => void;
  isLoading: boolean;
}

const SiteDetailsHeader: React.FC<SiteDetailsHeaderProps> = ({
  addSite,
  siteList,
  selectedSiteId,
  onSiteChange,
  isLoading,
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
                <CheckCircle color="success" />
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
        <Menu
          anchorEl={anchorEl}
          open={Boolean(anchorEl)}
          onClose={handleMenuClose}
        >
          {isLoading ? (
            <>
              <MenuItem>
                <Skeleton variant="rectangular" height={30} width={200} />
              </MenuItem>
              <MenuItem>
                <Skeleton variant="rectangular" height={30} width={200} />
              </MenuItem>
              <MenuItem>
                <Skeleton variant="rectangular" height={30} width={200} />
              </MenuItem>
            </>
          ) : (
            <>
              {siteList.map((site) => (
                <MenuItem
                  key={site.id}
                  onClick={() => handleSiteSelect(site.id)}
                  selected={selectedSiteId === site.id}
                >
                  {site.name}
                </MenuItem>
              ))}
              <MenuItem onClick={addSite}>Add site</MenuItem>
            </>
          )}
        </Menu>
      </Box>
    </Grid>
  );
};

export default SiteDetailsHeader;
