import { DarkTooltip } from '@/styles/global';
import { colors } from '@/styles/theme';
import SearchBar from '@/ui/molecules/SearchBar';
import AddLocationIcon from '@mui/icons-material/AddLocation';
import LocationOnOutlinedIcon from '@mui/icons-material/LocationOnOutlined';
import PowerIcon from '@mui/icons-material/PowerSettingsNewOutlined';
import RouteOutlinedIcon from '@mui/icons-material/RouteOutlined';
import {
  Box,
  IconButton,
  Stack,
  Tab,
  Tabs,
  Tooltip,
  Typography,
} from '@mui/material';

import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import DotIcon from '@mui/icons-material/FiberManualRecord';
import { LatLngLiteral } from 'leaflet';
import { useState } from 'react';
import TabPanel from '../TabPanel';

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
  handleAddSite: Function;
  handleAddLink: Function;
  isCurrentDraft: boolean;
  handleLocationSelected: (loc: LatLngLiteral) => void;
}

const LeftOverlayUI = ({
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
            sx={LeftIconButtonStyle}
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
            sx={LeftIconButtonStyle}
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
  id: string | undefined;
  handleClick: Function;
  isCurrentDraft: boolean;
  handleTogglePower: Function;
}

const RightOverlayUI = ({
  id,
  handleClick,
  isCurrentDraft,
  handleTogglePower,
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
      <Tooltip title="Turn Site On/Off">
        <IconButton
          sx={RightIconButtonStyle}
          onClick={(e) => {
            e.stopPropagation();
            handleTogglePower();
          }}
        >
          <PowerIcon htmlColor={colors.vulcan} />
        </IconButton>
      </Tooltip>
      <Tooltip title="Site Info">
        <IconButton
          aria-describedby={id}
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

interface IPlanningSummary {
  subtitleOne: string;
  subtitleTwo: string;
  siteSummary: any;
  powerSummary: any;
}

const tabProps = (index: number) => ({
  id: `sites-summary-tab-${index}`,
  'aria-controls': `sites-summary-tabpanel-${index}`,
});

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

const SiteSummary = ({ siteSummary }: any) => (
  <>
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
  </>
);
const PowerSummary = ({ powerSummary }: any) => (
  <Stack direction={'column'} spacing={1}>
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

const PlanningSummary = ({
  subtitleOne,
  subtitleTwo,
  siteSummary,
  powerSummary,
}: IPlanningSummary) => {
  const [value, setValue] = useState(0);
  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };
  return (
    <Stack spacing={1.2}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider', marginBottom: 0.8 }}>
        <Tabs value={value} onChange={handleChange} aria-label="summary tabs">
          <Tab
            label={subtitleOne}
            {...tabProps(0)}
            sx={{ fontSize: '12px', fontWeight: 400, p: 0 }}
          />
          <Tab
            label={subtitleTwo}
            {...tabProps(1)}
            sx={{ fontSize: '12px', fontWeight: 400, p: 0 }}
          />
        </Tabs>
      </Box>
      <TabPanel id={'sites-summary-tab'} value={value} index={0}>
        <SiteSummary siteSummary={siteSummary} />
      </TabPanel>
      <TabPanel id={'sites-power-tab'} value={value} index={1}>
        <PowerSummary powerSummary={powerSummary} />
      </TabPanel>
    </Stack>
  );
};

export { LeftOverlayUI, PlanningSummary, RightOverlayUI };
