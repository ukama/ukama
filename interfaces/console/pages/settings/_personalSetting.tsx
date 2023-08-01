import { LANGUAGE_OPTIONS, TimeZones } from '@/constants';
import { colors } from '@/styles/theme';
import { ExportOptionsType } from '@/types';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import SettingsArrowIcon from '@mui/icons-material/CallMade';
import {
  Button,
  Divider,
  Grid,
  MenuItem,
  Paper,
  TextField,
  Typography,
} from '@mui/material';
import { useState } from 'react';
const defaultTimeZone = 'Pacific Standard Time';
// localStorage['timeZone']? localStorage['timeZone']:'Pacific Standard Time'

export default function PersonalSettings() {
  const [language, setLanguage] = useState('en');
  const [timezone, setTimezone] = useState(defaultTimeZone);

  const handleLanguageChange = (event: any) => {
    setLanguage(event.target.value);
    // localStorage.setItem('i18nextLng', event.target.value);
  };

  const handleTimezoneChange = (event: any) => {
    setTimezone(event.target.value);
    // localStorage.setItem('timeZone', event.target.value);
  };

  const handleAccountSettings = () => {};
  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        borderRadius: '5px',
        height: 'calc(100vh - 200px)',
      }}
    >
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          overflow: 'auto',
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <Grid container spacing={2}>
          <Grid item container xs={12} spacing={2}>
            <Grid item xs={12} sm={4}>
              <Typography variant="h6">My Account Details</Typography>
            </Grid>
            <Grid item xs={12} sm={8}>
              <Button
                size="large"
                variant="outlined"
                endIcon={<SettingsArrowIcon />}
                onClick={handleAccountSettings}
              >
                UKAMA ACCOUNT SETTINGS
              </Button>
            </Grid>
            <Grid item xs={12}>
              <Divider />
            </Grid>
          </Grid>

          <Grid item container xs={12} spacing={2}>
            <Grid item xs={12} md={4}>
              <Typography variant="h6">Language & Region</Typography>
            </Grid>
            <Grid item md={8} xs={12} spacing={3} container>
              <Grid item xs={12} sm={12} md={8}>
                <TextField
                  select
                  id="language"
                  label="LANGUAGE"
                  value={language}
                  sx={{ width: '100%' }}
                  onChange={handleLanguageChange}
                >
                  {LANGUAGE_OPTIONS.map(
                    ({ value, label }: ExportOptionsType) => (
                      <MenuItem key={value} value={value}>
                        <Typography variant="body1">{label}</Typography>
                      </MenuItem>
                    ),
                  )}
                </TextField>
              </Grid>
              <Grid item xs={12} sm={12} md={8}>
                <TextField
                  select
                  id="timezone"
                  label="TIME ZONE"
                  value={timezone}
                  onChange={handleTimezoneChange}
                  sx={{ width: '100%' }}
                >
                  {TimeZones.map(({ value, text }: any) => (
                    <MenuItem key={value} value={value}>
                      <Typography variant="body1">{text}</Typography>
                    </MenuItem>
                  ))}
                </TextField>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </LoadingWrapper>
    </Paper>
  );
}
