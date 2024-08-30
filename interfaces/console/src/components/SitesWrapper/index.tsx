/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { SiteDto } from '@/client/graphql/generated';
import CellTowerIcon from '@mui/icons-material/CellTower';
import { Grid, Skeleton, Stack, Typography } from '@mui/material';
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
      <Stack
        spacing={1}
        height="100%"
        direction={'column'}
        alignItems={'center'}
      >
        <CellTowerIcon htmlColor="#9E9E9E" />
        <Typography variant="body2">No sites available.</Typography>
      </Stack>
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
