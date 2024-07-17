/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  IconButton,
  useMediaQuery,
  Menu,
  MenuItem,
} from '@mui/material';
import PeopleIcon from '@mui/icons-material/People';
import RouterIcon from '@mui/icons-material/Router';
import BatteryChargingFullIcon from '@mui/icons-material/BatteryChargingFull';
import CellTowerIcon from '@mui/icons-material/CellTower';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import colors from '@/theme/colors';
import { useRouter } from 'next/navigation';

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
}

const SiteCard: React.FC<SiteCardProps> = ({
  siteId,
  name,
  address,
  users,
  status,
  siteStatus,
  onClickMenu,
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
              <a
                href={`/console/sites/${siteId}`} // Replace with actual link URL
                style={{ textDecoration: 'none', color: 'inherit' }}
              >
                {name}
              </a>
            </Typography>
            <Typography color="textSecondary" variant="body2">
              {address}
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
          <Box display="flex" alignItems="center" mr={10}>
            <PeopleIcon sx={{ fontSize: 30 }} />
            <Typography variant="body2" ml={0.5}>
              {users}
            </Typography>
          </Box>
          <Box display="flex" alignItems="center" mr={2}>
            <RouterIcon
              sx={{ color: status.online ? 'green' : 'red', fontSize: 30 }}
            />
            {!isSmallScreen && (
              <Typography variant="body2" ml={0.5}>
                Online
              </Typography>
            )}
          </Box>
          <Box display="flex" alignItems="center" mr={2}>
            <BatteryChargingFullIcon
              fontSize="large"
              sx={{ color: status.charging ? 'green' : 'red', fontSize: 30 }}
            />
            {!isSmallScreen && (
              <Typography variant="body2" ml={0.5}>
                Charging
              </Typography>
            )}
          </Box>
          <Box display="flex" alignItems="center">
            <CellTowerIcon
              fontSize="large"
              sx={{
                color: siteStatus ? 'green' : 'red',
                fontSize: 30,
              }}
            />
            {!isSmallScreen && (
              <Typography variant="body2" ml={0.5}>
                {siteStatus ? 'Site is set up' : 'Site is deactivated'}
              </Typography>
            )}
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteCard;
