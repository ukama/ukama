import { commonData } from '@/app-recoil';
import { useGetSitesQuery } from '@/generated';
import { TCommonData } from '@/types';
import SiteOverallHealth from '@/ui/molecules/SiteHealth';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import { colors } from '@/styles/theme';
import SiteHeader from '@/ui/molecules/SiteHeader';
import { Site } from '@/types';
import { PageContainer } from '@/styles/global';
import dynamic from 'next/dynamic';
import { useRecoilValue } from 'recoil';
import { Grid, Typography, Paper, Stack } from '@mui/material';
const DynamicMap = dynamic(
  () => import('../ui/molecules/NetworkMap/DynamicMap'),
  {
    ssr: false,
  },
);

const sites: Site[] = [
  { name: 'site1', health: 'online', duration: '3 days' },
  { name: 'site2', health: 'offline', duration: '1 week' },
  { name: 'site3', health: 'online', duration: '2 days' },
];

export default function Page() {
  const _commonData = useRecoilValue<TCommonData>(commonData);

  const { data: networkRes, loading: networkLoading } = useGetSitesQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      networkId: _commonData?.networkId,
    },
  });

  const handleSiteSelect = (site: any): void => {
    console.log(site);
  };
  const handleAddSite = () => {
    // Logic to add a new site
  };
  const handleSiteRestart = () => {
    // Logic to restart a site
  };

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
          <Grid item xs={4}>
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
          <Grid item xs={8}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                Map
              </Typography>
            </Paper>
          </Grid>

          <Grid item xs={12}>
            <Paper sx={{ p: 2 }}>
              <Typography variant="h6" gutterBottom>
                Site components
              </Typography>

              <SiteOverallHealth />
            </Paper>
          </Grid>
        </Grid>
      </LoadingWrapper>
    </>
  );
}
