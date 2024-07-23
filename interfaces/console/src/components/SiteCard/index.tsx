/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import BatteryChargingFullIcon from '@mui/icons-material/BatteryChargingFull';
import CellTowerIcon from '@mui/icons-material/CellTower';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import RouterIcon from '@mui/icons-material/Router';
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
  users: number;
  siteStatus: boolean;
  status: {
    online: boolean;
    charging: boolean;
    signal: string;
  };
  onClickMenu?: (siteId: string) => void;
  loading?: boolean;
}

const SiteCard: React.FC<SiteCardProps> = ({
  siteId,
  name,
  address,
  users,
  status,
  siteStatus,
  onClickMenu,
  loading = false,
}) => {
  const isSmallScreen = useMediaQuery((theme: any) =>
    theme.breakpoints.down('sm'),
  );
  const router = useRouter();

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation(); // Stop propagation here
    setAnchorEl(event.currentTarget);
  };
  const handleMenuClick = (action: string) => {
    if (onClickMenu) {
      onClickMenu(siteId);
      if (action === 'edit') {
        /* empty */
      } else if (action === 'details') {
        onClickMenu(siteId);
      }
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  const navigateToDetails = () => {
    router.push(`/console/sites/${siteId}`);
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
                  href={`/console/sites/${siteId}`} // Replace with actual link URL
                  style={{ textDecoration: 'none', color: 'inherit' }}
                >
                  {name}
                </a>
              )}
            </Typography>
            <Typography color="textSecondary" variant="body2">
              {loading ? <Skeleton width={200} /> : address}
            </Typography>
          </Box>
          <IconButton onClick={handleClick}>
            <MoreVertIcon />
          </IconButton>

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleClose}
            onClick={(e) => e.stopPropagation()} // Stop propagation here
          >
            <MenuItem onClick={() => handleMenuClick('edit')}>
              Edit Name
            </MenuItem>
            <MenuItem onClick={() => handleMenuClick('details')}>
              View Site Details
            </MenuItem>
          </Menu>
        </Box>
        <Box
          display="flex"
          justifyContent="flex-start"
          alignItems="center"
          mt={2}
        >
          {/* TODO: Will be implemented in the future */}
          {/* <Box display="flex" alignItems="center" mr={10}>
            <PeopleIcon sx={{ fontSize: 30 }} />
            <Typography variant="body2" ml={0.5}>
              {loading ? <Skeleton width={30} /> : users}
            </Typography>
          </Box> */}
          <Box display="flex" alignItems="center" mr={2}>
            <RouterIcon
              sx={{
                color: loading ? 'gray' : status.online ? 'green' : 'red',
                fontSize: 30,
              }}
            />
            {!isSmallScreen && (
              <Typography variant="body2" ml={0.5}>
                {loading ? <Skeleton width={50} /> : 'Online'}
              </Typography>
            )}
          </Box>
          <Box display="flex" alignItems="center" mr={2}>
            <BatteryChargingFullIcon
              fontSize="large"
              sx={{
                color: loading ? 'gray' : status.charging ? 'green' : 'red',
                fontSize: 30,
              }}
            />
            {!isSmallScreen && (
              <Typography variant="body2" ml={0.5}>
                {loading ? <Skeleton width={70} /> : 'Charging'}
              </Typography>
            )}
          </Box>
          <Box display="flex" alignItems="center">
            <CellTowerIcon
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
