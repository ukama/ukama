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
import { Dispatch, SetStateAction } from 'react';

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
  zIndex: 400,
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
  search: string;
  handleAddSite: Function;
  handleAddLink: Function;
  setSearch: Dispatch<SetStateAction<string>>;
}

const LeftOverlayUI = ({
  search,
  setSearch,
  handleAddLink,
  handleAddSite,
}: ILeftOverlayUI) => (
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
        value={search}
        key={'searchbox'}
        handleOnChange={(v: string) => setSearch(v)}
        placeholderText="Search for a location, address, or coordinates"
      />
      <DarkTooltip title="Place site" placement="right-end">
        <IconButton sx={LeftIconButtonStyle} onClick={() => handleAddSite()}>
          <AddLocationIcon htmlColor="white" />
        </IconButton>
      </DarkTooltip>
      <DarkTooltip title="Add Link" placement="right-end">
        <IconButton sx={LeftIconButtonStyle} onClick={() => handleAddLink()}>
          <RouteOutlinedIcon htmlColor="white" />
        </IconButton>
      </DarkTooltip>
    </Stack>
  </Box>
);

interface IRightOverlayUI {
  id: string | undefined;
  handleClick: Function;
  handleTogglePower: Function;
}

const RightOverlayUI = ({
  id,
  handleClick,
  handleTogglePower,
}: IRightOverlayUI) => (
  <Box
    sx={{
      top: 24,
      right: 24,
      display: 'flex',
      position: 'absolute',
    }}
  >
    <Stack direction={'row'} spacing={1} alignItems={'flex-end'}>
      <Tooltip title="Turn Site On/Off">
        <IconButton
          sx={RightIconButtonStyle}
          onClick={() => handleTogglePower()}
        >
          <PowerIcon htmlColor={colors.vulcan} />
        </IconButton>
      </Tooltip>
      <Tooltip title="Site Info">
        <IconButton
          aria-describedby={id}
          sx={RightIconButtonStyle}
          onClick={(e) => handleClick(e)}
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
