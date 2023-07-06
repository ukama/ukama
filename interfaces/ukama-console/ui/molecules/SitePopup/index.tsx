import { SITE_PLANNING_AP_OPTIONS, SOLAR_UPTIME_OPTIONS } from '@/constants';
import { Site } from '@/generated/planning-tool';
import {
  Button,
  Divider,
  FormControl,
  FormControlLabel,
  FormLabel,
  InputAdornment,
  Paper,
  Radio,
  RadioGroup,
  Stack,
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import { useState } from 'react';

interface ISitePopup {
  site: Site;
  coverageLoading: boolean;
  handleAction: (a: Site) => void;
  handleDeleteSite: (i: string) => void;
  handleGenerateAction: (a: string, b: Site) => void;
}

const SitePopup = ({
  site,
  handleAction,
  coverageLoading,
  handleDeleteSite,
  handleGenerateAction,
}: ISitePopup) => {
  const [data, setData] = useState<Site>(site);
  return (
    <Paper elevation={0} sx={{ boxShadow: 'none', cursor: 'default' }}>
      <Stack spacing={1.2}>
        <Typography variant="h6">Place site</Typography>
        <TextField
          required
          label="NAME"
          value={data.name}
          variant="standard"
          sx={{
            '& .MuiInput-input': {
              fontSize: '16px',
            },
          }}
          InputLabelProps={{ shrink: true }}
          onChange={(e) => setData({ ...data, name: e.target.value })}
        />
        <TextField
          required
          value={
            data.location
              ? data.location.address
                ? data.location.address
                : `(${data.location?.lat}, ${data.location?.lng})`
              : '(0,0)'
          }
          label="LOCATION"
          variant="standard"
          disabled
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
          value={data.height}
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
          onChange={(e) =>
            setData({ ...data, height: parseFloat(e.target.value) })
          }
        />
        <FormControl>
          <FormLabel
            id="ap-selection"
            required
            sx={{ fontSize: '12px !important' }}
          >
            AP SELECTION
          </FormLabel>
          <RadioGroup
            value={data.apOption}
            name="ap-selection-group"
            aria-labelledby="ap-selection"
            defaultValue={SITE_PLANNING_AP_OPTIONS[0].value}
            onChange={(e) => setData({ ...data, apOption: e.target.value })}
          >
            {SITE_PLANNING_AP_OPTIONS.map(({ id, label, value }) => (
              <FormControlLabel
                key={id}
                value={value}
                label={label}
                control={<Radio />}
                sx={{ '.MuiTypography-root': { fontSize: '16px' } }}
              />
            ))}
          </RadioGroup>
        </FormControl>
        <FormControl>
          <FormLabel
            id="solar-uptime-selection"
            required
            sx={{ fontSize: '12px !important' }}
          >
            SOLAR SYSTEM TARGET UPTIME
          </FormLabel>
          <RadioGroup
            row
            value={data.solarUptime}
            name="solar-uptime-selection-group"
            aria-labelledby="solar-uptime-selection"
            defaultValue={SOLAR_UPTIME_OPTIONS[0].value}
            onChange={(e) =>
              setData({ ...data, solarUptime: parseInt(e.target.value) })
            }
          >
            {SOLAR_UPTIME_OPTIONS.map(({ id, label, value }) => (
              <FormControlLabel
                key={id}
                value={value}
                label={label}
                control={<Radio />}
                sx={{ '.MuiTypography-root': { fontSize: '16px' } }}
              />
            ))}
          </RadioGroup>
        </FormControl>
        <FormControl>
          <FormLabel
            required
            id="backhaul-option"
            sx={{ fontSize: '12px !important' }}
          >
            BACKHAUL OPTION
          </FormLabel>
          <FormControlLabel
            control={
              <Switch
                defaultChecked
                value={data.isSetlite}
                onChange={(e) =>
                  setData({ ...data, isSetlite: e.target.checked })
                }
              />
            }
            label="Add satellite"
            sx={{
              width: 'fit-content',
              '.MuiTypography-root': { fontSize: '16px' },
            }}
          />
        </FormControl>

        <Stack direction="row" justifyContent={'space-between'}>
          <Button
            variant="contained"
            color="error"
            sx={{ width: 'fit-content', fontSize: '12px' }}
            onClick={(e) => {
              e.preventDefault();
              handleDeleteSite(site.id);
            }}
          >
            Delete SITE
          </Button>
          <Button
            variant="contained"
            sx={{ width: 'fit-content', fontSize: '12px' }}
            onClick={(e) => {
              e.preventDefault();
              handleAction(data);
            }}
          >
            UPDATE SITE
          </Button>
        </Stack>
        {site?.id && (
          <Stack direction={'column'} spacing={1}>
            <Divider sx={{ width: '100%', height: '1px' }} />
            <Typography variant="caption">Generate Actions</Typography>
            <Stack
              direction="row"
              justifyContent={'space-between'}
              spacing={0.7}
            >
              <Button
                variant="contained"
                disabled={coverageLoading}
                sx={{
                  p: 1,
                  fontSize: '10px',
                  width: 'fit-content',
                  textTransform: 'capitalize',
                }}
                onClick={(e) => {
                  e.preventDefault();
                  handleGenerateAction('field_strength', site);
                }}
              >
                Field Strength
              </Button>
              <Button
                variant="contained"
                disabled={coverageLoading}
                sx={{
                  p: 1,
                  fontSize: '10px',
                  width: 'fit-content',
                  textTransform: 'capitalize',
                }}
                onClick={(e) => {
                  e.preventDefault();
                  handleGenerateAction('receive_power', site);
                }}
              >
                Receive Power
              </Button>
              <Button
                variant="contained"
                disabled={coverageLoading}
                sx={{
                  p: 1,
                  fontSize: '10px',
                  width: 'fit-content',
                  textTransform: 'capitalize',
                }}
                onClick={(e) => {
                  e.preventDefault();
                  handleGenerateAction('path_loss', site);
                }}
              >
                Path Loss
              </Button>
            </Stack>
          </Stack>
        )}
      </Stack>
    </Paper>
  );
};

export default SitePopup;
