import { DarkTooltip } from '@/styles/global';
import { colors } from '@/styles/theme';
import SearchBar from '@/ui/molecules/SearchBar';
import AddLocationIcon from '@mui/icons-material/AddLocation';
import LocationOnOutlinedIcon from '@mui/icons-material/LocationOnOutlined';
import PowerIcon from '@mui/icons-material/PowerSettingsNewOutlined';
import RouteOutlinedIcon from '@mui/icons-material/RouteOutlined';
import { Box, IconButton, Stack, Tooltip, Typography } from '@mui/material';

import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import DotIcon from '@mui/icons-material/FiberManualRecord';

const SITES_MOCK = [
  {
    id: 1,
    name: 'Site Name 1',
    status: 'up',
  },
  {
    id: 2,
    name: 'Site Name 2',
    status: 'down',
  },
  {
    id: 3,
    name: 'Site Name 3',
    status: 'unknown',
  },
];

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

const LeftOverlayUI = () => (
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
    <Stack spacing={1.5} width={'400px'} alignItems={'flex-start'}>
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
);

const RightOverlayUI = ({ id, handleClick }: any) => (
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
);

const SiteSummary = () => (
  <Stack spacing={1.2}>
    <Stack direction={'row'} alignItems={'center'} spacing={0.5}>
      <Typography variant="body2" sx={{ fontSize: 14, fontWeight: 600 }}>
        Site summary
      </Typography>
      <Typography variant="caption">{`(${SITES_MOCK.length})`}</Typography>
    </Stack>
    {SITES_MOCK.map(({ id, name, status }) => (
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
    ))}
  </Stack>
);

export { LeftOverlayUI, RightOverlayUI, SiteSummary };
