/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { SiteDto } from '@/client/graphql/generated';
import CellTowerIcon from '@mui/icons-material/CellTower';
import BatteryChargingFullIcon from '@mui/icons-material/BatteryChargingFull';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import GroupIcon from '@mui/icons-material/Group';
import RouterIcon from '@mui/icons-material/Router';
import SignalCellularAltIcon from '@mui/icons-material/SignalCellularAlt';
import {
  Grid,
  Skeleton,
  Stack,
  Typography,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Paper,
  useMediaQuery,
  useTheme,
  Button,
} from '@mui/material';
import { useState } from 'react';
import colors from '@/theme/colors';
import { useRouter } from 'next/navigation';

interface ISitesWrapper {
  sites: SiteDto[];
  loading: boolean;
  handleSiteNameUpdate: any;
  subscriberCount: number | undefined;
  unnamedNodes: any[];
  handleConfigureSite: (nodeId: string) => void;
}

const SiteCardSkeleton = (
  <Skeleton
    variant="rectangular"
    height={158}
    width="100%"
    sx={{ borderRadius: '4px' }}
  />
);

const SitesWrapper = ({
  loading,
  sites,
  subscriberCount,
  handleSiteNameUpdate,
  unnamedNodes,
  handleConfigureSite,
}: ISitesWrapper) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedSiteId, setSelectedSiteId] = useState<string | null>(null);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const router = useRouter();

  const handleMenuOpen = (
    event: React.MouseEvent<HTMLElement>,
    siteId: string,
  ) => {
    setAnchorEl(event.currentTarget);
    setSelectedSiteId(siteId);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedSiteId(null);
  };

  const handleSiteNameClick = (siteId: string) => {
    router.push(`/console/sites/${siteId}`);
  };

  if (loading)
    return (
      <Grid container columnSpacing={2}>
        {[1, 2, 3].map((item) => (
          <Grid item xs={12} md={4} key={item}>
            {SiteCardSkeleton}
          </Grid>
        ))}
      </Grid>
    );

  if (sites.length === 0 && unnamedNodes.length === 0)
    return (
      <Box
        sx={{
          padding: 3,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          textAlign: 'center',
          maxWidth: '600px',
          margin: '0 auto',
          minHeight: '200px',
          justifyContent: 'center',
        }}
      >
        <CellTowerIcon
          sx={{
            fontSize: 48,
            color: colors.black54,
            mb: 2,
          }}
        />
        <Typography
          variant="h6"
          sx={{
            mb: 1,
            fontWeight: 600,
          }}
        >
          No sites yet!
        </Typography>
        <Typography variant="body2" color="text.secondary">
          A site is a complete connection point to the network, made up of your
          Ukama node, and the power and backhaul components. Install a site to
          get started - read more here.
        </Typography>
      </Box>
    );

  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      {sites.map((site) => (
        <Grid item xs={12} md={4} lg={4} key={site.id}>
          <Paper
            sx={{
              padding: '16px',
              display: 'flex',
              flexDirection: 'column',
              border: `1px solid ${colors.black10}`,
              borderRadius: '5px',
            }}
            elevation={3}
          >
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                mb: 1,
              }}
            >
              <Typography
                variant="h6"
                sx={{
                  fontWeight: 500,
                  textDecoration: 'underline',
                  cursor: 'pointer',
                }}
                onClick={() => handleSiteNameClick(site.id)}
              >
                {site.name}
              </Typography>
              <IconButton
                onClick={(event) => handleMenuOpen(event, site.id)}
                size="small"
              >
                <MoreVertIcon />
              </IconButton>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl) && selectedSiteId === site.id}
                onClose={handleMenuClose}
              >
                <MenuItem
                  onClick={() => {
                    handleMenuClose();
                    handleSiteNameUpdate(site.id, site.name);
                  }}
                >
                  Edit Name
                </MenuItem>
              </Menu>
            </Box>
            <Typography
              variant="body2"
              color="text.secondary"
              sx={{
                mb: 3,
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
              }}
            >
              {site.location}
            </Typography>
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                width: '100%',
                gap: 1,
              }}
            >
              <Stack direction="row" spacing={0.5} alignItems="center">
                <GroupIcon />
                <Typography variant="body2">{subscriberCount || 0}</Typography>
              </Stack>

              <Box
                sx={{
                  display: 'flex',
                  gap: 1,
                  alignItems: 'center',
                }}
              >
                <RouterIcon color="success" />
                {!isMobile && <Typography variant="body2">Online</Typography>}

                <BatteryChargingFullIcon color="success" />
                {!isMobile && <Typography variant="body2">Charged</Typography>}

                <SignalCellularAltIcon color="success" />
                {!isMobile && <Typography variant="body2">Strong</Typography>}
              </Box>
            </Box>
          </Paper>
        </Grid>
      ))}

      {unnamedNodes.map((node) => (
        <Grid item xs={12} md={4} lg={4} key={node.id}>
          <Paper
            sx={{
              padding: '16px',
              display: 'flex',
              flexDirection: 'column',
              border: `1px solid ${colors.black10}`,
              borderRadius: '5px',
            }}
            elevation={3}
          >
            <Stack direction="column" spacing={1}>
              <Typography
                variant="h6"
                sx={{
                  fontWeight: 500,
                  textDecoration: 'underline',
                  cursor: 'pointer',
                }}
              >
                Undefined site
              </Typography>
              <Typography
                variant="body2"
                color="text.secondary"
                sx={{
                  mb: 3,
                  whiteSpace: 'nowrap',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                }}
              >
                Node ID: {node.id}
              </Typography>
              <Button
                variant="contained"
                sx={{
                  width: { xs: '100%', md: '50%', xl: '50%' },
                }}
                onClick={() => handleConfigureSite(node.id)}
              >
                Configure Site
              </Button>
            </Stack>
          </Paper>
        </Grid>
      ))}
    </Grid>
  );
};

export default SitesWrapper;