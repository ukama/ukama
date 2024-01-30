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
import { SiteDto } from '@/generated';

interface SiteCardProps {
  sites: SiteDto[];
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
          <Link href={`/sites/${site.id}`} unselectable="on" legacyBehavior>
            <Grid item xs={12} sm={6} sx={{ cursor: 'pointer' }}>
              <Stack direction="row" spacing={1} alignItems={'center'}>
                <Typography variant="h5">{site.name}</Typography>
                <ErrorIcon sx={{ color: colors.red, fontSize: 18 }} />
              </Stack>
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
              <MenuItem onClick={() => handleDelete(site.id)}>
                <Typography variant="body1" sx={{ color: colors.red }}>
                  {' '}
                  Delete site
                </Typography>
              </MenuItem>
            </Menu>
          </Grid>

          <Grid item xs={6} sm={2}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <PersonIcon />
              <Typography variant="caption">{`${23}`}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <NodeIcon status={`online`} />
              <Typography variant="caption">{`online`}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={4}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <BatteryIcon status={`charging`} />
              <Typography variant="caption">{`charging`}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <TowerIcon status={`offline`} />
              <Typography variant="caption">{`online`}</Typography>
            </Stack>
          </Grid>
        </Grid>
      ))}
    </RoundedCard>
  );
};

export default SiteCard;
