import { SiteDto } from '@/client/graphql/generated';
import { Box, Grid, Skeleton, Typography } from '@mui/material';
import SiteCard from '../SiteCard';

interface ISitesWrapper {
  sites: SiteDto[];
  loading: boolean;
  handleSiteNameUpdate: any;
}

const SiteCardSkelton = (
  <Skeleton
    variant="rectangular"
    height={158}
    width={'100%'}
    sx={{ borderRadius: '4px' }}
  />
);

const SitesWrapper = ({
  loading,
  sites,
  handleSiteNameUpdate,
}: ISitesWrapper) => {
  if (loading)
    return (
      <Grid container columnSpacing={2}>
        {[1, 2, 3].map((item) => (
          <Grid item xs={12} md={4} key={item}>
            {SiteCardSkelton}
          </Grid>
        ))}
      </Grid>
    );

  if (sites.length === 0)
    return (
      <Box
        sx={{
          width: '100%',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Typography variant="body1">No sites available.</Typography>
      </Box>
    );

  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      {sites.map((site) => (
        <Grid item xs={12} md={4} lg={4} key={site.id}>
          <SiteCard
            siteId={site.id}
            name={site.name}
            loading={loading}
            address={site.location}
            siteStatus={site.isDeactivated}
            handleSiteNameUpdate={handleSiteNameUpdate}
          />
        </Grid>
      ))}
    </Grid>
  );
};

export default SitesWrapper;
