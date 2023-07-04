import { Site } from '@/generated/planning-tool';
import { DarkTooltip } from '@/styles/global';
import { colors } from '@/styles/theme';
import SearchBar from '@/ui/molecules/SearchBar';
import AddLocationIcon from '@mui/icons-material/AddLocation';
import BatteryOutlinedIcon from '@mui/icons-material/BatteryChargingFullOutlined';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import DotIcon from '@mui/icons-material/FiberManualRecord';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import LocationOnOutlinedIcon from '@mui/icons-material/LocationOnOutlined';
import RouteOutlinedIcon from '@mui/icons-material/RouteOutlined';
import SignalIcon from '@mui/icons-material/SignalCellularAlt';
import SpeedIcon from '@mui/icons-material/Speed';
import {
  Box,
  Button,
  Grid,
  IconButton,
  InputAdornment,
  Stack,
  TextField,
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
  handleAddSite: Function;
  handleAddLink: Function;
  isCurrentDraft: boolean;
  handleLocationSelected: (loc: LatLngLiteral) => void;
}

export const LeftOverlayUI = ({
  isAddSite,
  isAddLink,
  handleAddLink,
  handleAddSite,
  isCurrentDraft,
  handleLocationSelected,
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
        key={'searchbox'}
        handleLocationSelected={handleLocationSelected}
        placeholderText="Search for a location, address, or coordinates"
      />
      {isCurrentDraft && (
        <DarkTooltip title="Place site" placement="right-end">
          <IconButton
            sx={{
              ...LeftIconButtonStyle,
              backgroundColor: isAddSite
                ? colors.primaryMain
                : colors.primaryDark,
            }}
            onClick={(e) => {
              e.bubbles = false;
              handleAddSite();
            }}
          >
            <AddLocationIcon htmlColor="white" />
          </IconButton>
        </DarkTooltip>
      )}
      {isCurrentDraft && (
        <DarkTooltip title="Add Link" placement="right-end">
          <IconButton
            sx={{
              ...LeftIconButtonStyle,
              backgroundColor: isAddLink
                ? colors.primaryMain
                : colors.primaryDark,
            }}
            onClick={(e) => {
              e.stopPropagation();
              handleAddLink();
            }}
          >
            <RouteOutlinedIcon htmlColor="white" />
          </IconButton>
        </DarkTooltip>
      )}
    </Stack>
  </Box>
);

interface IRightOverlayUI {
  handleClick: Function;
  isCurrentDraft: boolean;
  handlePowerInfo: Function;
  siteInfoId: string | undefined;
  powerInfoId: string | undefined;
}

export const RightOverlayUI = ({
  siteInfoId,
  powerInfoId,
  handleClick,
  isCurrentDraft,
  handlePowerInfo,
}: IRightOverlayUI) => (
  <Box
    sx={{
      top: 24,
      right: 24,
      position: 'absolute',
      display: isCurrentDraft ? 'flex' : 'none',
    }}
  >
    <Stack direction={'row'} spacing={1} alignItems={'flex-end'}>
      <Tooltip title="Power Info">
        <IconButton
          aria-describedby={powerInfoId}
          sx={RightIconButtonStyle}
          onClick={(e) => {
            e.stopPropagation();
            handlePowerInfo(e);
          }}
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
          onClick={(e) => {
            e.stopPropagation();
            handleClick(e);
          }}
        >
          <LocationOnOutlinedIcon htmlColor={colors.vulcan} />
        </IconButton>
      </Tooltip>
    </Stack>
  </Box>
);

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

export const SiteSummary = ({ siteSummary }: any) => (
  <Stack spacing={1}>
    <Typography variant="body2" sx={{ fontWeight: 500 }}>
      {`Site Summary (${siteSummary.length})`}
    </Typography>
    {siteSummary.length > 0 ? (
      siteSummary.map(({ id, name, status }: any) => (
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
);

export const PowerSummary = ({ powerSummary }: any) => (
  <Stack spacing={1}>
    <Typography variant="body2" sx={{ fontWeight: 500 }}>
      Power Summary
    </Typography>
    {PowerSummarySections.map(({ id, title, unit }) => {
      const d: any = {
        'power-usage': {
          total: 0,
          info: '',
        },
        'solar-panels': {
          total: 0,
          info: '',
        },
        batteries: {
          total: 0,
          info: '',
        },
      };

      powerSummary.sites.forEach(
        ({ id: _id, name, usage, panels, battries }: any, i: number) => {
          const isLastItem = powerSummary.sites.length === i + 1;
          d['power-usage'].total = d['power-usage'].total + usage;
          d['solar-panels'].total = d['solar-panels'].total + panels;
          d['batteries'].total = d['batteries'].total + battries;
          d['power-usage'].info =
            d['power-usage'].info +
            `${name} (${usage})` +
            (isLastItem ? '' : ' + ');
          d['solar-panels'].info =
            d['solar-panels'].info +
            `${name} (${panels})` +
            (isLastItem ? '' : ' + ');
          d['batteries'].info =
            d['batteries'].info +
            `${name} (${battries})` +
            (isLastItem ? '' : ' + ');
        },
      );

      return (
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
              {d[id].total} {unit}
            </Typography>
          </Stack>
          <Typography variant="body2" sx={{ fontSize: 14, fontWeight: 300 }}>
            {d[id].info ? d[id].info : 'No data available!'}
          </Typography>
        </Stack>
      );
    })}
  </Stack>
);

interface ISiteDetails {
  site: Site;
}
const SiteDetails = ({ site }: ISiteDetails) => (
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
      onChange={(e) => {}}
    />
  </Stack>
);

interface ISites {
  sites: Site[];
  handleDeleteLink: () => void;
}

export const SiteLink = ({ sites, handleDeleteLink }: ISites) => (
  <Grid container height="100%" columnSpacing={3}>
    <Grid item xs={3.5}>
      <SiteDetails site={sites[0]} />
    </Grid>
    <Grid item direction={'row'} xs={5}>
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
            sx={{ height: 'fit-content' }}
            onClick={() => handleDeleteLink()}
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
);
