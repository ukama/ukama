import { SITE_PLANNING_AP_OPTIONS, SOLAR_UPTIME_OPTIONS } from '@/constants';
import { Site } from '@/generated/planning-tool';
import {
  Button,
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
import { Dispatch, SetStateAction } from 'react';

interface ISitePopup {
  data: Site;
  handleAction: () => void;
  setData: Dispatch<SetStateAction<any>>;
}

const SitePopup = ({ data, setData, handleAction }: ISitePopup) => {
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
            data.location.address
              ? data.location.address
              : `(${data.location.lat}, ${data.location.lng})`
          }
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
          onChange={(e) => setData({ ...data, height: e.target.value })}
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
            onChange={(e) => setData({ ...data, ap: e.target.value })}
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
            onChange={(e) => setData({ ...data, solarUptime: e.target.value })}
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
                  setData({ ...data, isBackhaul: e.target.checked })
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
        <Stack direction="row" justifyContent={'flex-end'}>
          <Button
            variant="contained"
            sx={{ width: 'fit-content', fontSize: '14px' }}
            onClick={(e) => {
              handleAction();
              e.preventDefault();
            }}
          >
            {data.name ? 'UPDATE SITE' : 'PLACE SITE'}
          </Button>
        </Stack>
      </Stack>
    </Paper>
  );
};

export default SitePopup;
