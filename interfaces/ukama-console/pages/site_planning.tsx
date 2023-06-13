import styles from '@/styles/Site_Planning.module.css';
import { DarkTooltip, PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import Map from '@/ui/molecules/Map';
import SearchBar from '@/ui/molecules/SearchBar';
import AddLocationIcon from '@mui/icons-material/AddLocation';
import LocationOnOutlinedIcon from '@mui/icons-material/LocationOnOutlined';
import PowerIcon from '@mui/icons-material/PowerSettingsNewOutlined';
import RouteOutlinedIcon from '@mui/icons-material/RouteOutlined';
import {
  Box,
  IconButton,
  Popover,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material';
import { useState } from 'react';
const DEFAULT_CENTER = [38.907132, -77.036546];

const LeftIconButtonStyle = {
  borderRadius: '4px',
  backgroundColor: colors.primaryDark,
  boxShadow:
    '0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)',
  ':hover': {
    backgroundColor: colors.primaryDark,
    svg: {
      path: {
        fill: colors.white,
      },
    },
  },
};
const RightIconButtonStyle = {
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

const SiteSummary = () => (
  <Stack>
    <Stack flexDirection={'row'} alignItems={'center'} spacing={1}>
      <Typography variant="body2" sx={{ fontSize: 14, fontWeight: 600 }}>
        Site summary
      </Typography>
      <Typography variant="caption">{`(${5})`}</Typography>
    </Stack>
    <Typography variant="h4" sx={{ fontWeight: 600, marginBottom: '8px' }}>
      Site 1
    </Typography>
  </Stack>
);

export default function Page() {
  const [anchorSiteInfo, setAnchorSiteInfo] =
    useState<HTMLButtonElement | null>(null);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorSiteInfo(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorSiteInfo(null);
  };

  const open = Boolean(anchorSiteInfo);
  const id = open ? 'site-info-popover' : undefined;

  return (
    <>
      <Popover
        id={id}
        open={open}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        onClose={handleClose}
        anchorEl={anchorSiteInfo}
        sx={{ top: 4, left: -40 }}
        PaperProps={{
          sx: {
            width: '204px',
            padding: '16px 24px',
          },
        }}
      >
        <SiteSummary />
      </Popover>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <PageContainer sx={{ padding: 0 }}>
          <Map
            zoom={12}
            width={800}
            height={418}
            center={DEFAULT_CENTER}
            className={styles.homeMap}
          >
            {({ TileLayer, Marker, Popup, Pane }: any) => (
              <>
                <Box
                  sx={{
                    top: 24,
                    left: 24,
                    width: '100%',
                    zIndex: 1000,
                    display: 'flex',
                    position: 'absolute',
                  }}
                >
                  <Stack
                    spacing={1.5}
                    width={'400px'}
                    alignItems={'flex-start'}
                  >
                    <SearchBar
                      value=""
                      key={'searchbox'}
                      handleOnChange={() => {}}
                      placeholderText="Search for a location, address, or coordinates"
                    />
                    <DarkTooltip title="Place site" placement="right-end">
                      <IconButton sx={LeftIconButtonStyle}>
                        <AddLocationIcon htmlColor="white" />
                      </IconButton>
                    </DarkTooltip>
                    <DarkTooltip title="Add Link" placement="right-end">
                      <IconButton sx={LeftIconButtonStyle}>
                        <RouteOutlinedIcon htmlColor="white" />
                      </IconButton>
                    </DarkTooltip>
                  </Stack>
                </Box>
                <Box
                  sx={{
                    top: 24,
                    right: 24,
                    zIndex: 1000,
                    display: 'flex',
                    position: 'absolute',
                  }}
                >
                  <Stack direction={'row'} spacing={1} alignItems={'flex-end'}>
                    <Tooltip title="Turn Site On/Off">
                      <IconButton sx={RightIconButtonStyle}>
                        <PowerIcon htmlColor={colors.vulcan} />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Site Info">
                      <IconButton
                        aria-describedby={id}
                        onClick={handleClick}
                        sx={RightIconButtonStyle}
                      >
                        <LocationOnOutlinedIcon htmlColor={colors.vulcan} />
                      </IconButton>
                    </Tooltip>
                  </Stack>
                </Box>

                <TileLayer url="https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png" />
                <Marker position={DEFAULT_CENTER}>
                  <Popup>Site Info</Popup>
                </Marker>
              </>
            )}
          </Map>
        </PageContainer>
      </LoadingWrapper>
    </>
  );
}
