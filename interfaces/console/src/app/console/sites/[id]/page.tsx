/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import SiteComponents from '@/app/console/sites/[id]/_components/SiteComponents';
import SiteInfo from '@/app/console/sites/[id]/_components/SiteInfos';
import SiteOverview from '@/app/console/sites/[id]/_components/SiteOverView';
import { useSiteDetailPage } from '@/app/console/sites/[id]/_hooks/useSiteDetailPage';
import StatusBar from '@/app/console/_components/StatusBar';
import { Box, Grid2, Skeleton } from '@mui/material';
import dynamic from 'next/dynamic';
import React from 'react';

const SiteMapComponent = dynamic(
  () => import('@/app/console/sites/[id]/_components/SiteMapComponent'),
  { ssr: false },
);

interface SiteDetailsProps {
  params: { id: string };
}

const Page: React.FC<SiteDetailsProps> = ({ params }) => {
  const vm = useSiteDetailPage(params.id);

  if (!vm.isDataReady) {
    return (
      <Grid2 container columnSpacing={2} rowSpacing={2}>
        <Grid2 size={12}>
          <Skeleton height={64} width={'100%'} variant="rectangular" sx={{ borderRadius: '5px' }} />
        </Grid2>
        {[1, 2, 3].map((item) => (
          <Grid2 size={4} key={item}>
            <Skeleton height={164} width={'100%'} variant="rectangular" sx={{ borderRadius: '5px' }} />
          </Grid2>
        ))}
        <Grid2 size={12}>
          <Skeleton height={300} width={'100%'} variant="rectangular" sx={{ borderRadius: '5px' }} />
        </Grid2>
      </Grid2>
    );
  }

  return (
    <Box
      sx={{
        overflowY: 'auto',
        overflowX: 'hidden',
        borderRadius: '10px',
        width: '100%',
        height: 'calc(100vh - 164px)',
      }}
    >
      <Grid2 container spacing={2} alignItems="stretch" sx={{ mt: 1, height: 'max-content' }}>
        <Grid2 size={12}>
          <StatusBar
            type="toggle"
            selected={vm.activeSite}
            uptime={vm.siteUptime}
            actionLoading={vm.healthLoading || vm.toggleRFStatusLoading || vm.toggleServiceLoading}
            objs={vm.siteData?.getSites.sites ?? []}
            handleActionClick={vm.handleActionClick}
            actionOptions={vm.actionOptions}
            handleSelected={vm.handleSelected}
            actionOptionValues={vm.siteActionData}
          />
        </Grid2>

        <Grid2 size={{ xs: 12, sm: 6, md: 4 }} sx={{ height: 'auto', display: 'flex' }}>
          <SiteInfo
            selectedSite={vm.activeSite}
            address={vm.currentSiteAddress}
            nodeIds={vm.nodes.map((node) => node.id)}
          />
        </Grid2>

        <Grid2 size={{ xs: 12, sm: 6, md: 5 }} sx={{ height: '100%', display: 'flex' }}>
          <SiteOverview
            installationDate={new Date(vm.activeSite.installDate)}
            isLoading={vm.statLoading}
            siteId={vm.activeSite.id}
            siteStatMetrics={vm.statData?.getSiteStat ?? { metrics: [] }}
          />
        </Grid2>

        <Grid2 size={{ xs: 12, sm: 6, md: 3 }} sx={{ height: 'auto', display: 'flex', minHeight: 200 }}>
          <SiteMapComponent
            id="site-map"
            zoom={15}
            posix={[vm.activeSite.latitude ?? '0', vm.activeSite.longitude ?? '0']}
            address={vm.currentSiteAddress}
            height={'100%'}
            mapStyle="satellite"
            showUserCount={true}
            userCount={vm.activeSubscribers}
          />
        </Grid2>

        <Grid2 size={12}>
          <SiteComponents
            key={`${vm.activeView.kpi}-${vm.metricFrom}`}
            siteId={vm.activeSite.id}
            metrics={vm.metrics}
            sections={vm.sections}
            activeKPI={vm.activeView.kpi}
            activeSection={vm.getSectionName(vm.activeView.graphType)}
            metricFrom={vm.metricFrom}
            metricsLoading={vm.metricsLoading}
            onComponentClick={vm.handleViewChange}
            onSwitchChange={vm.handleSwitchChange}
            nodes={vm.nodes}
            initialNodeUptimes={vm.initialNodeUptimes}
          />
        </Grid2>
      </Grid2>
    </Box>
  );
};

export default Page;
