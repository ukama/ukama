import React, { useState } from 'react';
import {
  Grid,
  Typography,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Stack,
} from '@mui/material';
import { RoundedCard } from '@/styles/global';
import { PersonIcon, TowerIcon, NodeIcon, BatteryIcon } from '../SvgIcons';
import ErrorIcon from '@mui/icons-material/Error';
import { colors } from '@/styles/theme';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import Link from 'next/link';

interface Site {
  id: string;
  name: string;
  details: string;
  batteryStatus: 'charging' | 'notCharging';
  nodeStatus: 'online' | 'offline';
  towerStatus: 'online' | 'offline';
  numberOfPersonsConnected: number;
}

interface SiteCardProps {
  sites: Site[];
  handleDeleteSite: (siteId?: string) => void;
}

const SiteCard: React.FC<SiteCardProps> = ({ sites, handleDeleteSite }) => {
  const [anchorEl, setAnchorEl] = useState(null);

  const handleMenuOpen = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleDelete = (siteId?: string) => {
    handleDeleteSite(siteId);
    handleMenuClose();
  };

  return (
    <RoundedCard>
      {sites.map((site, index) => (
        <Grid container spacing={1} alignItems={'center'} key={index}>
          <Link
            href={`sites/${site.id}`}
            passHref
            legacyBehavior
            replace
            unselectable="on"
          >
            <Grid item xs={12} sm={6} sx={{ cursor: 'pointer' }}>
              <Stack direction="row" spacing={1} alignItems={'center'}>
                <Typography variant="h5">{site.name}</Typography>
                {(site.towerStatus === 'offline' ||
                  site.nodeStatus === 'offline' ||
                  site.batteryStatus === 'notCharging') && (
                  <ErrorIcon sx={{ color: colors.red, fontSize: 18 }} />
                )}
              </Stack>
              <Typography variant="body2">{site.details}</Typography>
            </Grid>
          </Link>

          <Grid item xs={12} sm={6} container justifyContent={'flex-end'}>
            <IconButton onClick={handleMenuOpen}>
              <MoreVertIcon />
            </IconButton>
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleMenuClose}
            >
              <MenuItem onClick={() => handleDelete(site.id)}>Delete</MenuItem>
            </Menu>
          </Grid>

          <Grid item xs={6} sm={2}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <PersonIcon />
              <Typography variant="caption">
                {`${site.numberOfPersonsConnected}`}
              </Typography>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <NodeIcon status={site.nodeStatus} />
              <Typography variant="caption">{`${site.nodeStatus}`}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={4}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <BatteryIcon status={site.batteryStatus} />
              <Typography variant="caption">{`${site.batteryStatus}`}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <TowerIcon status={site.towerStatus} />
              <Typography variant="caption">{`${site.towerStatus}`}</Typography>
            </Stack>
          </Grid>
        </Grid>
      ))}
    </RoundedCard>
  );
};

export default SiteCard;
