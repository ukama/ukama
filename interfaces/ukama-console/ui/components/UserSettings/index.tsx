import { isDarkmode } from '@/app-recoil';
import { LANGUAGE_OPTIONS, TimeZones } from '@/constants';
import { SettingsArrowIcon } from '@/public/svg';
import { ExportOptionsType } from '@/types';
import {
  Button,
  Divider,
  Grid,
  MenuItem,
  TextField,
  Typography,
} from '@mui/material';
import { useState } from 'react';
import { useRecoilValue } from 'recoil';

const LineDivider = () => (
  <Grid item xs={12}>
    <Divider />
  </Grid>
);

const UserSettings = () => {
  const defaultTimeZone = localStorage['timeZone']
    ? localStorage['timeZone']
    : 'Pacific Standard Time';
  const _isDarkMod = useRecoilValue(isDarkmode);
  const [language, setLanguage] = useState(localStorage['i18nextLng']);
  const [timezone, setTimezone] = useState(defaultTimeZone);

  const handleLanguageChange = (event: any) => {
    setLanguage(event.target.value);
    localStorage.setItem('i18nextLng', event.target.value);
  };

  const handleTimezoneChange = (event: any) => {
    setTimezone(event.target.value);
    localStorage.setItem('timeZone', event.target.value);
  };

  const handleAccountSettings = () => {
    typeof window !== 'undefined' &&
      window.location.replace(
        `${process.env.REACT_APP_AUTH_URL}/userAccountSettings?mode=${
          _isDarkMod ? 1 : 0
        }`,
      );
  };

  return (
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
          <LineDivider />
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
              {LANGUAGE_OPTIONS.map(({ value, label }: ExportOptionsType) => (
                <MenuItem key={value} value={value}>
                  <Typography variant="body1">{label}</Typography>
                </MenuItem>
              ))}
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
  );
};

export default UserSettings;
