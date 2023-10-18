import { colors } from '@/styles/theme';
import { Site } from '@/types';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
// import Map from '@/ui/molecules/MapComponent';
import SiteHeader from '@/ui/molecules/SiteHeader';
import SiteOverallHealth from '@/ui/molecules/SiteHealth';
import { Grid, Paper, Stack, Typography } from '@mui/material';

const sites: Site[] = [
  { name: 'site1', health: 'online', duration: '3 days' },
  { name: 'site2', health: 'offline', duration: '1 week' },
  { name: 'site3', health: 'online', duration: '2 days' },
];

export default function Page() {
  const handleSiteSelect = (site: any): void => {};
  const handleAddSite = () => {
    // Logic to add a new site
  };
  const handleSiteRestart = () => {
    // Logic to restart a site
  };

  const batteryInfo = [
    { label: 'Model number', value: 'V1234' },
    { label: 'Current', value: '10 A' },
    { label: 'Charge', value: '80 %' },
    { label: 'Power', value: '100 W' },
    { label: 'Voltage', value: '12 V' },
  ];
  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <SiteHeader
          sites={sites}
          sitesAction={handleSiteSelect}
          addSiteAction={handleAddSite}
          restartSiteAction={handleSiteRestart}
        />

        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={12} sm={4}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                Site details
              </Typography>
              <Stack direction="column" spacing={2}>
                <Stack direction="row" spacing={2} alignItems={'center'}>
                  <Typography variant="subtitle1">Date created:</Typography>
                  <Typography variant="subtitle1"> July 13 2023</Typography>
                </Stack>
                <Stack direction="row" spacing={2} alignItems={'center'}>
                  <Typography variant="subtitle1"> Address:</Typography>
                  <Typography variant="subtitle1"> 1000 Nelson Way</Typography>
                </Stack>
              </Stack>
            </Paper>
          </Grid>
          <Grid item xs={12} sm={8}>
            {/* <Map site={''} users={0} /> */}
          </Grid>

          <Grid item xs={12}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                Site components
              </Typography>

              <SiteOverallHealth
                solarHealth={'warning'}
                nodeHealth={'good'}
                switchHealth={'good'}
                controllerHealth={'good'}
                batteryHealth={'good'}
                backhaulHealth={'good'}
                batteryInfo={batteryInfo}
              />
            </Paper>
          </Grid>
        </Grid>
      </LoadingWrapper>
    </>
  );
}
