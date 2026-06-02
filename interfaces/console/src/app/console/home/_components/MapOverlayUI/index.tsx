/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import React, { memo, useCallback, useMemo } from 'react';
import SearchBar from '@/components/ui/SearchBar';
import { DarkTooltip } from '@/styles/global';
import { colors } from '@/theme';
import AddLocationIcon from '@mui/icons-material/AddLocation';
import BatteryOutlinedIcon from '@mui/icons-material/BatteryChargingFullOutlined';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import DotIcon from '@mui/icons-material/FiberManualRecord';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import LocationOnOutlinedIcon from '@mui/icons-material/LocationOnOutlined';
import RouteOutlinedIcon from '@mui/icons-material/RouteOutlined';
import SatelliteIcon from '@mui/icons-material/Satellite';
import SignalIcon from '@mui/icons-material/SignalCellularAlt';
import SpeedIcon from '@mui/icons-material/Speed';
import TerrainIcon from '@mui/icons-material/Terrain';
import {
  Box,
  Button,
  Card,
  Grid,
  IconButton,
  InputAdornment,
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Tooltip,
  Typography,
} from '@mui/material';
import { LatLngLiteral } from 'leaflet';
import Image from 'next/image';

const LeftIconButtonStyle = {
  zIndex: 400,
  borderRadius: '4px',
  boxShadow:
    '0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)',
  ':hover': {
    backgroundColor: colors.primaryMain,
    svg: {
      path: {
        fill: colors.white,
      },
    },
  },
};
const RightIconButtonStyle = {
  zIndex: 400,
  borderRadius: '4px',
  backgroundColor: colors.whiteLilac,
  boxShadow:
    '0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)',
  ':hover': {
    backgroundColor: colors.white,
    svg: {
      path: {
        fill: colors.vulcan,
      },
      circle: {
        fill: colors.vulcan,
      },
    },
  },
};

interface ILeftOverlayUI {
  isAddSite: boolean;
  isAddLink: boolean;
  handleAddSite: () => void;
  handleAddLink: () => void;
  isCurrentDraft: boolean;
  handleLocationSelected: (loc: LatLngLiteral) => void;
}

export const LeftOverlayUI = memo(({
  isAddSite,
  isAddLink,
  handleAddLink,
  handleAddSite,
  isCurrentDraft,
  handleLocationSelected,
}: ILeftOverlayUI) => {
  const handleAddSiteClick = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      e.bubbles = false;
      handleAddSite();
    },
    [handleAddSite],
  );

  const handleAddLinkClick = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      e.stopPropagation();
      handleAddLink();
    },
    [handleAddLink],
  );

  const addSiteButtonSx = useMemo(
    () => ({
      ...LeftIconButtonStyle,
      backgroundColor: isAddSite ? colors.primaryMain : colors.primaryDark,
    }),
    [isAddSite],
  );

  const addLinkButtonSx = useMemo(
    () => ({
      ...LeftIconButtonStyle,
      backgroundColor: isAddLink ? colors.primaryMain : colors.primaryDark,
    }),
    [isAddLink],
  );

  return (
    <Box
      sx={{
        top: 24,
        left: 24,
        width: '100%',
        display: 'flex',
        position: 'absolute',
      }}
    >
      <Stack spacing={1.5} width={'400px'} alignItems={'flex-start'}>
        <SearchBar
          key={'searchbox'}
          handleLocationSelected={handleLocationSelected}
          placeholderText="Search for a location, address, or coordinates"
        />
        {isCurrentDraft && (
          <DarkTooltip title="Place site" placement="right-end">
            <IconButton sx={addSiteButtonSx} onClick={handleAddSiteClick}>
              <AddLocationIcon htmlColor="white" />
            </IconButton>
          </DarkTooltip>
        )}
        {isCurrentDraft && (
          <DarkTooltip title="Add Link" placement="right-end">
            <IconButton sx={addLinkButtonSx} onClick={handleAddLinkClick}>
              <RouteOutlinedIcon htmlColor="white" />
            </IconButton>
          </DarkTooltip>
        )}
      </Stack>
    </Box>
  );
});
LeftOverlayUI.displayName = 'LeftOverlayUI';

interface IRightOverlayUI {
  handleClick: (event: React.MouseEvent<HTMLElement>) => void;
  isCurrentDraft: boolean;
  handlePowerInfo: (event: React.MouseEvent<HTMLElement>) => void;
  siteInfoId: string | undefined;
  powerInfoId: string | undefined;
}

export const RightOverlayUI = memo(({
  siteInfoId,
  powerInfoId,
  handleClick,
  isCurrentDraft,
  handlePowerInfo,
}: IRightOverlayUI) => {
  const handlePowerInfoClick = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      e.stopPropagation();
      handlePowerInfo(e);
    },
    [handlePowerInfo],
  );

  const handleSiteInfoClick = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      e.stopPropagation();
      handleClick(e);
    },
    [handleClick],
  );

  const containerSx = useMemo(
    () => ({
      top: 24,
      right: 24,
      position: 'absolute',
      display: isCurrentDraft ? 'flex' : 'none',
    }),
    [isCurrentDraft],
  );

  return (
    <Box sx={containerSx}>
      <Stack direction={'row'} spacing={1} alignItems={'flex-end'}>
        <Tooltip title="Power Info">
          <IconButton
            aria-describedby={powerInfoId}
            sx={RightIconButtonStyle}
            onClick={handlePowerInfoClick}
          >
            <BatteryOutlinedIcon
              htmlColor={colors.vulcan}
              sx={{ transform: 'rotate(90deg)' }}
            />
          </IconButton>
        </Tooltip>
        <Tooltip title="Site Info">
          <IconButton
            aria-describedby={siteInfoId}
            sx={RightIconButtonStyle}
            onClick={handleSiteInfoClick}
          >
            <LocationOnOutlinedIcon htmlColor={colors.vulcan} />
          </IconButton>
        </Tooltip>
      </Stack>
    </Box>
  );
});
RightOverlayUI.displayName = 'RightOverlayUI';

const PowerSummarySections = [
  {
    id: 'power-usage',
    unit: 'W',
    title: 'Power Usage',
  },
  {
    id: 'solar-panels',
    title: 'Solar Panels',
    unit: '',
  },
  {
    id: 'batteries',
    title: 'Batteries',
    unit: '',
  },
];

interface SiteSummaryItem {
  id: string;
  name: string;
  status: string;
}

export const SiteSummary = memo(({ siteSummary }: { siteSummary: SiteSummaryItem[] }) => (
  <Stack spacing={1}>
    <Typography variant="body2" sx={{ fontWeight: 500 }}>
      {`Site Summary (${siteSummary.length})`}
    </Typography>
    {siteSummary.length > 0 ? (
      siteSummary.map(({ id, name, status }: SiteSummaryItem) => (
        <Stack key={id} direction="row" spacing={1} alignItems={'center'}>
          {status === 'unknown' ? (
            <DotIcon color={'disabled'} sx={{ fontSize: '18px' }} />
          ) : (
            <CheckCircleOutlineIcon
              color={status === 'up' ? 'success' : 'error'}
              sx={{ fontSize: '18px' }}
            />
          )}
          <Typography variant="body2" sx={{ fontSize: 14, fontWeight: 500 }}>
            {name}
          </Typography>
        </Stack>
      ))
    ) : (
      <Typography variant="caption" sx={{ fontSize: 14 }}>
        No site added yet!
      </Typography>
    )}
  </Stack>
));
SiteSummary.displayName = 'SiteSummary';

interface PowerSiteSummaryItem {
  id: string;
  name: string;
  usage: number;
  panels: number;
  battries: number;
}

interface PowerSummaryData {
  sites: PowerSiteSummaryItem[];
}

interface PowerSectionData {
  total: number;
  info: string;
}

export const PowerSummary = memo(({ powerSummary }: { powerSummary: PowerSummaryData }) => {
  const aggregated = useMemo<Record<string, PowerSectionData>>(() => {
    const d: Record<string, PowerSectionData> = {
      'power-usage': { total: 0, info: '' },
      'solar-panels': { total: 0, info: '' },
      batteries: { total: 0, info: '' },
    };

    powerSummary.sites.forEach(
      ({ name, usage, panels, battries }: PowerSiteSummaryItem, i: number) => {
        const isLastItem = powerSummary.sites.length === i + 1;
        const separator = isLastItem ? '' : ' + ';
        d['power-usage'].total += usage;
        d['solar-panels'].total += panels;
        d['batteries'].total += battries;
        d['power-usage'].info += `${name} (${usage})${separator}`;
        d['solar-panels'].info += `${name} (${panels})${separator}`;
        d['batteries'].info += `${name} (${battries})${separator}`;
      },
    );

    return d;
  }, [powerSummary.sites]);

  return (
    <Stack spacing={1}>
      <Typography variant="body2" sx={{ fontWeight: 500 }}>
        Power Summary
      </Typography>
      {PowerSummarySections.map(({ id, title, unit }) => (
        <Stack
          key={id}
          direction="column"
          spacing={0.6}
          alignItems={'flex-start'}
        >
          <Stack
            width={'100%'}
            direction="row"
            spacing={1}
            justifyContent={'space-between'}
          >
            <Typography variant="body2" sx={{ fontSize: 14, fontWeight: 500 }}>
              {title}
            </Typography>
            <Typography variant="body2" sx={{ fontSize: 14, fontWeight: 500 }}>
              {aggregated[id].total} {unit}
            </Typography>
          </Stack>
          <Typography variant="body2" sx={{ fontSize: 14, fontWeight: 300 }}>
            {aggregated[id].info ? aggregated[id].info : 'No data available!'}
          </Typography>
        </Stack>
      ))}
    </Stack>
  );
});
PowerSummary.displayName = 'PowerSummary';

interface PlanSiteForDetails {
  name?: string;
  height?: number;
  location: { lat: string; lng: string };
}

interface ISiteDetails {
  site: PlanSiteForDetails;
}

const SiteDetails = memo(({ site }: ISiteDetails) => (
  <Stack spacing={2} py={3}>
    <Stack direction={'row'} spacing={1} alignItems={'center'}>
      <LocationOnIcon fontSize="small" />
      <Typography variant="body2">{site?.name}</Typography>
    </Stack>
    <TextField
      required
      value={`${site?.location.lat}, ${site?.location.lng}`}
      label="LOCATION"
      variant="standard"
      InputLabelProps={{ shrink: true }}
      placeholder="Location, address, or coordinates"
      sx={{
        '& .MuiInput-input': {
          fontSize: '16px',
        },
      }}
      onChange={() => {}}
    />
    <TextField
      required
      value={site?.height}
      type="number"
      label="HEIGHT"
      variant="standard"
      InputLabelProps={{ shrink: true }}
      InputProps={{
        endAdornment: <InputAdornment position="end">m</InputAdornment>,
      }}
      sx={{
        width: { xs: '100%', sm: '100px' },
        '& .MuiInput-input': {
          fontSize: '16px',
        },
      }}
      onChange={() => {}}
    />
  </Stack>
));
SiteDetails.displayName = 'SiteDetails';

interface ISites {
  sites: PlanSiteForDetails[];
  handleDeleteLink: () => void;
}

export const SiteLink = memo(({ sites, handleDeleteLink }: ISites) => {
  const onDeleteLink = useCallback(() => {
    handleDeleteLink();
  }, [handleDeleteLink]);

  return (
    <Card sx={{ boxShadow: 'none' }}>
      <Grid container height="100%" columnSpacing={3}>
        <Grid item xs={3.5}>
          <SiteDetails site={sites[0]} />
        </Grid>
        <Grid item xs={5}>
          <Stack height={'100%'} sx={{ border: '0.5px solid grey' }}>
            <Stack
              height={'56px'}
              direction={'row'}
              alignItems={'center'}
              justifyContent={'space-around'}
            >
              <Stack direction={'row'} spacing={1} alignItems={'center'}>
                <SignalIcon fontSize="small" color="success" />
                <Typography variant="body2">-45 dBm</Typography>
              </Stack>
              <Stack direction={'row'} spacing={1} alignItems={'center'}>
                <SpeedIcon fontSize="small" color="success" />
                <Typography variant="body2">100 Mbps</Typography>
              </Stack>
              <Button
                size="small"
                color="error"
                variant="outlined"
                sx={{ height: 'fit-content', fontSize: '12px' }}
                onClick={onDeleteLink}
              >
                Delete Link
              </Button>
            </Stack>
            <Image
              width={502}
              height={170}
              src="/temp_link.png"
              alt="ukama-sites-link"
            />
          </Stack>
        </Grid>
        <Grid item xs={3.5}>
          <SiteDetails site={sites[1]} />
        </Grid>
      </Grid>
    </Card>
  );
});
SiteLink.displayName = 'SiteLink';

interface ILayerSwitch {
  value: string;
  handleLayerSwitch: (event: React.MouseEvent<HTMLElement>, value: string) => void;
}

export const LayerSwitch = memo(({ handleLayerSwitch, value }: ILayerSwitch) => {
  const satelliteButtonSx = useMemo(
    () => ({
      border: 0,
      borderRadius: '4px',
      backgroundColor: `${value === 'satellite' ? colors.primaryMain02 : 'transparent'} !important`,
    }),
    [value],
  );

  const terrainButtonSx = useMemo(
    () => ({
      border: 0,
      borderRadius: '4px',
      backgroundColor: `${value === 'terrain' ? colors.primaryMain02 : 'transparent'} !important`,
    }),
    [value],
  );

  const satelliteLabelSx = useMemo(
    () => ({
      fontSize: 12,
      textTransform: 'capitalize' as const,
      fontWeight: value === 'satellite' ? 600 : 300,
    }),
    [value],
  );

  const terrainLabelSx = useMemo(
    () => ({
      fontSize: 12,
      textTransform: 'capitalize' as const,
      fontWeight: value === 'terrain' ? 600 : 300,
    }),
    [value],
  );

  return (
    <Box
      sx={{
        left: 24,
        bottom: 24,
        zIndex: 400,
        position: 'absolute',
      }}
    >
      <Card variant="elevation" sx={{ p: 0.8, width: 'fit-content' }}>
        <ToggleButtonGroup
          value={value}
          exclusive
          onChange={handleLayerSwitch}
          aria-label="text alignment"
          sx={{ height: '36px' }}
        >
          <ToggleButton
            value="satellite"
            aria-label="left aligned"
            sx={satelliteButtonSx}
          >
            <Stack direction={'row'} alignItems={'center'} spacing={0.8}>
              <SatelliteIcon fontSize="small" />
              <Typography variant="caption" sx={satelliteLabelSx}>
                Satellite
              </Typography>
            </Stack>
          </ToggleButton>
          <ToggleButton
            value="terrain"
            aria-label="centered"
            sx={terrainButtonSx}
          >
            <Stack direction={'row'} alignItems={'center'} spacing={0.8}>
              <TerrainIcon fontSize="small" />
              <Typography variant="caption" sx={terrainLabelSx}>
                Terrain
              </Typography>
            </Stack>
          </ToggleButton>
        </ToggleButtonGroup>
      </Card>
    </Box>
  );
});
LayerSwitch.displayName = 'LayerSwitch';
