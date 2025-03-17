/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  Component_Type,
  SiteDto,
  useAddSiteMutation,
  useGetComponentsByUserIdLazyQuery,
  useGetNetworksQuery,
  useGetNodesByNetworkLazyQuery,
  useGetSiteLazyQuery,
  useGetSitesQuery,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  useGetMetricByTabLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import ConfigureSiteDialog from '@/components/ConfigureSiteDialog';
import SiteComponents, {
  KPI_TO_SECTION_MAP,
  SectionData,
} from '@/components/SiteComponents';
import SiteDetailsHeader from '@/components/SiteDetailsHeader';
import SiteInfo from '@/components/SiteInfos';
import SiteOverview from '@/components/SiteOverView';
import { SITE_KPIS } from '@/constants';
import { METRIC_RANGE_10800 } from '@/constants';
import { useAppContext } from '@/context';
import { TSiteForm } from '@/types';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { getUnixTime } from '@/utils';
import { AlertColor, Box, Grid, Paper, Skeleton } from '@mui/material';
import { formatISO } from 'date-fns';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/navigation';
import React, { useEffect, useState } from 'react';

const SiteMapComponent = dynamic(
  () => import('@/components/SiteMapComponent'),
  {
    ssr: false,
  },
);

const SITE_INIT = {
  switch: '',
  power: '',
  access: '',
  backhaul: '',
  address: '',
  spectrum: '',
  siteName: '',
  latitude: NaN,
  longitude: NaN,
  network: '',
};

const defaultSite: SiteDto = {
  id: '',
  accessId: '',
  backhaulId: '',
  createdAt: '',
  installDate: '',
  isDeactivated: false,
  latitude: 0,
  location: '',
  longitude: 0,
  name: '',
  networkId: '',
  powerId: '',
  spectrumId: '',
  switchId: '',
};

interface SiteDetailsProps {
  params: {
    id: string;
  };
}

const Page: React.FC<SiteDetailsProps> = ({ params }) => {
  const { id } = params;
  const [site, setSite] = useState<TSiteForm>(SITE_INIT);
  const [activeSite, setActiveSite] = useState<SiteDto>(defaultSite);
  const [openSiteConfig, setOpenSiteConfig] = useState(false);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const [selectedSiteId, setSelectedSiteId] = useState<string | null>(null);
  const [sitesList, setSitesList] = useState<SiteDto[]>([]);
  const [nodeIds, setNodeIds] = useState<string[]>([]);
  const {
    setSnackbarMessage,
    network,
    setSelectedDefaultSite,
    user,
    subscriptionClient,
  } = useAppContext();
  const {
    address: CurrentSiteaddress,
    isLoading: CurrentSiteAddressLoading,
    error,
    fetchAddress,
  } = useFetchAddress();
  const [isDataReady, setIsDataReady] = useState(false);

  // Site Components state
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [graphType, setGraphType] = useState<Graphs_Type>(Graphs_Type.Solar);
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [activeSection, setActiveSection] = useState<string>('SOLAR');
  const [activeKPI, setActiveKPI] = useState<string>('solar');

  const sections: SectionData = {
    SOLAR: SITE_KPIS.SOLAR.metrics,
    BATTERY: SITE_KPIS.BATTERY.metrics,
    CONTROLLER: SITE_KPIS.CONTROLLER.metrics,
    MAIN_BACKHAUL: SITE_KPIS.MAIN_BACKHAUL.metrics,
    SWITCH: SITE_KPIS.SWITCH.metrics,
  };

  const router = useRouter();
  const [selectedComponent, setSelectedComponent] = useState<any>('backhaul');

  const handleComponentSelect = (component: any) => {
    setSelectedComponent(component);
    console.log(`Selected component: ${component}`);
  };

  const handleSiteConfigOpen = () => {
    setOpenSiteConfig(true);
  };

  const handleCloseSiteConfig = () => {
    setOpenSiteConfig(false);
  };

  // Metrics query
  const [
    getMetricByTab,
    { loading: metricsLoading, variables: metricsVariables },
  ] = useGetMetricByTabLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setMetrics(data.getMetricByTab);
      console.log('METRICS :', data.getMetricByTab);
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'site-metrics-err-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: (res) => {
      setSnackbarMessage({
        id: 'add-site-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-subscriber-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
    variables: {
      networkId: network.id,
    },
    onCompleted: (res) => {
      setSitesList(res.getSites.sites);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'fetching-sites-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const handleSiteConfiguration = async (data: any) => {
    const variables = {
      access_id: data.access,
      backhaul_id: data.backhaul,
      install_date: formatISO(new Date()),
      latitude: data.coordinates.lat,
      location: data.location,
      longitude: data.coordinates.lng,
      name: data.siteName,
      network_id: data.selectedNetwork,
      power_id: data.power,
      spectrum_id: data.spectrumId || '',
      switch_id: data.switch,
      is_deactivated: data.is_deactivated || false,
    };

    try {
      await addSite({ variables: { data: variables } });
    } catch (error) {
      console.error('Error submitting site configuration:', error);
    }
  };

  const [getComponents] = useGetComponentsByUserIdLazyQuery({
    onCompleted: (res) => {
      if (res.getComponentsByUserId) {
        setComponentsList(res.getComponentsByUserId.components);
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'components-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: networks } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [getSite, { loading: getSiteLoading }] = useGetSiteLazyQuery({
    onCompleted: (res) => {
      setSite({
        power: res.getSite.powerId,
        siteName: res.getSite.name,
        switch: res.getSite.switchId,
        access: res.getSite.accessId,
        address: res.getSite.location,
        latitude: res.getSite.latitude,
        network: res.getSite.networkId,
        spectrum: res.getSite.spectrumId,
        backhaul: res.getSite.backhaulId,
        longitude: res.getSite.longitude,
      });
      setActiveSite({
        id: res.getSite.id,
        accessId: res.getSite.accessId,
        backhaulId: res.getSite.backhaulId,
        createdAt: res.getSite.createdAt,
        installDate: res.getSite.installDate,
        isDeactivated: res.getSite.isDeactivated,
        latitude: res.getSite.latitude,
        location: res.getSite.location,
        longitude: res.getSite.longitude,
        name: res.getSite.name,
        networkId: res.getSite.networkId,
        powerId: res.getSite.powerId,
        spectrumId: res.getSite.spectrumId,
        switchId: res.getSite.switchId,
      });
      checkDataReadiness();
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'sites-msg',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [fetchNode] = useGetNodesByNetworkLazyQuery({
    onCompleted: (res) => {
      if (res.getNodesByNetwork?.nodes) {
        const ids = res.getNodesByNetwork.nodes.map((node) => node.id);
        setNodeIds(ids);
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'nodes-msg',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  // Metrics-related effects
  useEffect(() => {
    if (selectedSiteId) {
      setMetricFrom(getUnixTime() - METRIC_RANGE_10800);
    }
  }, [selectedSiteId]);

  useEffect(() => {
    if (
      metricFrom > 0 &&
      selectedSiteId &&
      metricsVariables?.data?.from !== metricFrom
    ) {
      getMetricByTab({
        variables: {
          data: {
            step: 30,
            siteId: selectedSiteId,
            userId: user.id,
            nodeId: '',
            type: graphType,
            from: metricFrom,
            orgName: user.orgName,
            withSubscription: false,
            to: metricFrom + METRIC_RANGE_10800,
          },
        },
      });
    }
  }, [
    metricFrom,
    graphType,
    metricsVariables?.data?.from,
    selectedSiteId,
    user.id,
    user.orgName,
    getMetricByTab,
  ]);

  const handleSectionChange = (section: string): void => {
    setActiveSection(section);

    if (section !== 'NODE') {
      let newGraphType: Graphs_Type;
      switch (section) {
        case 'SOLAR':
          newGraphType = Graphs_Type.Solar;
          break;
        case 'BATTERY':
          newGraphType = Graphs_Type.Battery;
          break;
        case 'CONTROLLER':
          newGraphType = Graphs_Type.Controller;
          break;
        case 'MAIN_BACKHAUL':
          newGraphType = Graphs_Type.MainBackhaul;
          break;
        default:
          newGraphType = Graphs_Type.Solar;
      }

      setGraphType(newGraphType);
      setMetricFrom(getUnixTime() - METRIC_RANGE_10800);
    }
  };

  const handleNodeClick = (kpiType: string) => {
    setActiveKPI(kpiType);
    const sectionName = KPI_TO_SECTION_MAP[kpiType] || 'SOLAR';
    handleSectionChange(sectionName);
  };

  useEffect(() => {
    getComponents({
      variables: {
        data: {
          category: Component_Type.All,
        },
      },
    });
  }, []);

  useEffect(() => {
    getSite({ variables: { siteId: id } });
  }, []);

  useEffect(() => {
    if (id) {
      setSelectedSiteId(id);
    } else if (sitesList.length > 0) {
      setSelectedSiteId(sitesList[0].id);
    }
  }, [id, sitesList]);

  const handleSiteChange = (newSiteId: string) => {
    setSelectedSiteId(newSiteId);
    router.push(`/console/sites/${newSiteId}`);
  };

  useEffect(() => {
    if (selectedSiteId) {
      getSite({ variables: { siteId: selectedSiteId } });
    }
  }, [selectedSiteId]);

  useEffect(() => {
    const handleFetchAddress = async () => {
      setSnackbarMessage({
        id: 'fetching-address',
        type: 'info',
        show: true,
        message: 'Fetching address with coordinates',
      });
      await fetchAddress(activeSite.latitude, activeSite.longitude);
    };

    setSelectedDefaultSite(activeSite.name);

    if (activeSite && activeSite.latitude && activeSite.longitude) {
      handleFetchAddress();
    }
  }, [activeSite, setSnackbarMessage, fetchAddress, setSelectedDefaultSite]);

  useEffect(() => {
    if (error) {
      setSnackbarMessage({
        id: 'error-fetching-address',
        type: 'error',
        show: true,
        message: 'Error fetching address from coordinates',
      });
    }
  }, [error, setSnackbarMessage]);

  const checkDataReadiness = () => {
    if (
      activeSite.id &&
      CurrentSiteaddress &&
      !getSiteLoading &&
      !CurrentSiteAddressLoading
    ) {
      setIsDataReady(true);
    }
  };

  useEffect(() => {
    checkDataReadiness();
  }, [
    activeSite,
    CurrentSiteaddress,
    getSiteLoading,
    CurrentSiteAddressLoading,
  ]);

  useEffect(() => {
    if (activeSite?.networkId) {
      fetchNode({ variables: { networkId: activeSite.networkId } });
    }
  }, [activeSite, fetchNode]);

  if (!isDataReady) {
    return (
      <Grid container columnSpacing={2}>
        {[1, 2].map((item) => (
          <Grid item xs={6} key={item}>
            <Skeleton
              variant="rectangular"
              height={158}
              width={'100%'}
              sx={{ borderRadius: '5px' }}
            />
          </Grid>
        ))}
      </Grid>
    );
  }

  return (
    <Box
      sx={{
        overflow: 'auto',
        borderRadius: '10px',
      }}
    >
      <SiteDetailsHeader
        addSite={handleSiteConfigOpen}
        siteList={sitesList || []}
        selectedSiteId={selectedSiteId}
        onSiteChange={handleSiteChange}
        isLoading={sitesLoading}
      />

      <Grid
        container
        spacing={2}
        sx={{
          mt: 1,
          height: 'calc(50vh - 50px)',
        }}
      >
        <Grid item xs={4} sx={{ height: '100%' }}>
          <SiteInfo
            selectedSite={activeSite}
            address={CurrentSiteaddress}
            nodeIds={nodeIds}
          />
        </Grid>
        <Grid item xs={5} sx={{ height: '100%' }}>
          <SiteOverview uptimePercentage={98.5} daysRange={60} />
        </Grid>
        <Grid item xs={3} sx={{ height: '100%' }}>
          <SiteMapComponent
            posix={[activeSite.latitude, activeSite.longitude]}
            address={CurrentSiteaddress}
          />
        </Grid>
      </Grid>

      <Box sx={{ mt: 4, mb: 4 }}>
        <SiteComponents
          siteId={selectedSiteId || ''}
          metrics={metrics}
          sections={sections}
          nodeIds={nodeIds}
          activeKPI={activeKPI}
          activeSection={activeSection}
          metricFrom={metricFrom}
          metricsLoading={metricsLoading}
          onNodeClick={handleNodeClick}
        />
      </Box>

      <ConfigureSiteDialog
        site={site}
        open={openSiteConfig}
        addSiteLoading={addSiteLoading}
        onClose={handleCloseSiteConfig}
        components={componentsList || []}
        networks={networks?.getNetworks?.networks || []}
        handleSiteConfiguration={handleSiteConfiguration}
      />
    </Box>
  );
};

export default Page;
