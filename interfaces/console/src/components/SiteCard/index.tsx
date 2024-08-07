/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import {
  Box,
  Card,
  CardContent,
  IconButton,
  Menu,
  MenuItem,
  Skeleton,
  Typography,
  useMediaQuery,
} from '@mui/material';
import { useRouter } from 'next/navigation';
import React, { useState } from 'react';

interface SiteCardProps {
  siteId: string;
  name: string;
  address: string;
  siteStatus: boolean;
  onClickMenu?: (siteId: string) => void;
  loading?: boolean;
  handleSiteNameUpdate: (siteId: string, newSiteName: string) => void;
}

const SiteCard: React.FC<SiteCardProps> = ({
  siteId,
  name,
  address,
  siteStatus,
  onClickMenu,
  handleSiteNameUpdate,
  loading = false,
}) => {
  const isSmallScreen = useMediaQuery((theme: any) =>
    theme.breakpoints.down('sm'),
  );
  const router = useRouter();

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };
  const handleMenuClick = (action: string) => {
    if (handleSiteNameUpdate) {
      handleSiteNameUpdate(siteId, name);
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  const navigateToDetails = () => {
    router.push(`/console/sites/${siteId}`);
  };

  const truncateAddress = (address: string, limit: number) => {
    const segments = address.split(',');
    return (
      segments.slice(0, limit).join(', ') +
      (segments.length > limit ? '...' : '')
    );
  };

  return (
    <Card
      sx={{
        border: `1px solid ${colors.darkGradient}`,
        borderRadius: 2,
        marginBottom: 2,
        '&:hover': {
          border: `1px solid ${colors.primaryDark}`,
          cursor: 'pointer',
        },
      }}
      onClick={navigateToDetails}
    >
      <CardContent>
        <Box
          display="flex"
          justifyContent="space-between"
          flexDirection={isSmallScreen ? 'column' : 'row'}
        >
          <Box mb={isSmallScreen ? 2 : 0}>
            <Typography variant="h6">
              {loading ? (
                <Skeleton width={150} />
              ) : (
                <a
                  href={`/console/sites/${siteId}`}
                  style={{ textDecoration: 'none', color: 'inherit' }}
                >
                  {name.charAt(0).toUpperCase() + name.slice(1)}
                </a>
              )}
            </Typography>
            <Typography color="textSecondary" variant="body2">
              {loading ? <Skeleton width={200} /> : truncateAddress(address, 3)}
            </Typography>
          </Box>
          <IconButton onClick={handleClick}>
            <MoreVertIcon />
          </IconButton>

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleClose}
            onClick={(e) => e.stopPropagation()}
          >
            <MenuItem onClick={() => handleMenuClick('edit')}>
              Edit Name
            </MenuItem>
          </Menu>
        </Box>
        <Box
          display="flex"
          justifyContent="flex-start"
          alignItems="center"
          mt={2}
        >
          <Box display="flex" alignItems="center">
            <CheckCircleIcon
              fontSize="large"
              sx={{
                color: loading ? 'gray' : siteStatus ? 'green' : 'red',
                fontSize: 30,
              }}
            />
            {!isSmallScreen && (
              <Typography variant="body2" ml={0.5}>
                {loading ? (
                  <Skeleton width={100} />
                ) : siteStatus ? (
                  'Site is set up'
                ) : (
                  'Site is deactivated'
                )}
              </Typography>
            )}
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteCard;
