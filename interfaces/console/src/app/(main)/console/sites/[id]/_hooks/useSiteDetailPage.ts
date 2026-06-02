/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { Node, SiteDto } from '@/client/graphql/generated';
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { SITE_ACTIONS_BUTTONS, SITE_KPIS } from '@/constants';
import { SectionData } from '@/constants/index';
import { useNetworkContext } from '@/context';
import { ActiveView, TSiteActionToggle, TStatusBarObj } from '@/types';
import { graphTypeToSection } from '@/utils';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useSiteActions } from './useSiteActions';
import { useSiteData } from './useSiteData';
import { useSiteMetrics } from './useSiteMetrics';

const defaultSite: SiteDto = {
  id: '',
  accessId: '',
  backhaulId: '',
  createdAt: '',
  installDate: '',
  isDeactivated: false,
  latitude: '',
  location: '',
  longitude: '',
  name: '',
  networkId: '',
  powerId: '',
  spectrumId: '',
  switchId: '',
};

export function useSiteDetailPage(id: string) {
  const router = useRouter();

  // Shared state that is used across multiple hooks
  const [activeSite, setActiveSite] = useState<SiteDto>(defaultSite);
  const [nodes, setNodes] = useState<Node[]>([]);
  const [nodesFetched, setNodesFetched] = useState(false);
  const [isDataReady, setIsDataReady] = useState(false);
  const [siteActionData, setSiteActionData] = useState<TSiteActionToggle[]>([]);
  const [activeView, setActiveView] = useState<ActiveView>({
    graphType: Graphs_Type.Solar,
    kpi: 'node',
  });

  const { setSelectedDefaultSite } = useNetworkContext();

  // Address lookup
  const {
    address: currentSiteAddress,
    isLoading: currentSiteAddressLoading,
    error: addressError,
    fetchAddress,
  } = useFetchAddress();

  // GraphQL queries / mutations
  const {
    siteData,
    fetchNodesForSite,
    toggleRFStatus,
    toggleRFStatusLoading,
    toggleService,
    toggleServiceLoading,
    updateSwitchPort,
    healthLoading,
    notify,
  } = useSiteData(
    id,
    activeSite,
    nodes,
    setNodes,
    setNodesFetched,
    setSiteActionData,
  );

  // Metric subscriptions and derived values
  const {
    metrics,
    metricFrom,
    metricsLoading,
    statData,
    statLoading,
    resetMetrics,
    activeSubscribers,
    initialNodeUptimes,
    siteUptime,
  } = useSiteMetrics(id, activeSite, nodes, nodesFetched, activeView);

  // Event handlers
  const {
    handleViewChange,
    handleSwitchChange,
    handleSiteChange,
    handleActionClick,
    handleSelected,
  } = useSiteActions({
    id,
    nodes,
    setSiteActionData,
    setActiveView,
    resetMetrics,
    toggleRFStatus,
    toggleService,
    updateSwitchPort,
  });

  // Static section map
  const sections: SectionData = useMemo(
    () => ({
      SOLAR: SITE_KPIS.SOLAR.metrics,
      BATTERY: SITE_KPIS.BATTERY.metrics,
      CONTROLLER: SITE_KPIS.CONTROLLER.metrics,
      MAIN_BACKHAUL: SITE_KPIS.MAIN_BACKHAUL.metrics,
      SWITCH: SITE_KPIS.SWITCH.metrics,
    }),
    [],
  );

  const getSectionName = useCallback(
    (graphType: Graphs_Type): string =>
      graphTypeToSection[graphType] || 'SOLAR',
    [],
  );

  // Data-readiness check
  const checkDataReadiness = useCallback(() => {
    if (activeSite.id && currentSiteAddress && !currentSiteAddressLoading) {
      setIsDataReady(true);
    }
  }, [activeSite.id, currentSiteAddress, currentSiteAddressLoading]);

  const filterActiveSite = useCallback(
    (siteId: string) => {
      const found = siteData?.getSites.sites.find((s) => s.id === siteId);
      if (found) {
        setActiveSite(found);
        checkDataReadiness();
      }
    },
    [siteData, checkDataReadiness],
  );

  useEffect(() => {
    if (id) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      filterActiveSite(id);
      setActiveView({ graphType: Graphs_Type.NodeHealth, kpi: 'node' });
    }
  }, [id, filterActiveSite]);

  useEffect(() => {
    if (siteData?.getSites?.sites) {
      const found = siteData.getSites.sites.find((s) => s.id === id);
      if (found) {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setActiveSite(found);
      } else if (siteData.getSites.sites.length > 0 && id === '') {
        router.push('/console/sites/' + siteData.getSites.sites[0].id);
      }
    }
  }, [id, siteData, router]);

  useEffect(() => {
    const handleFetchAddress = async () => {
      if (activeSite.latitude && activeSite.longitude) {
        await fetchAddress(
          activeSite.latitude.toString(),
          activeSite.longitude.toString(),
        );
      }
    };
     
    setSelectedDefaultSite(activeSite.name);
    if (activeSite.id && activeSite.latitude && activeSite.longitude) {
      handleFetchAddress();
    }
  }, [activeSite, fetchAddress, setSelectedDefaultSite]);

  useEffect(() => {
    if (addressError) {
      notify(
        'error-fetching-address',
        'Error fetching address from coordinates',
        'error',
      );
    }
  }, [addressError, notify]);

  useEffect(() => {
    if (activeSite.id) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setNodesFetched(false);
      fetchNodesForSite({ variables: { siteId: activeSite.id } });
    }
  }, [activeSite.id, fetchNodesForSite]);

  return {
    id,
    activeSite,
    nodes,
    isDataReady,
    activeSubscribers,
    siteActionData,
    activeView,
    sections,
    metrics,
    metricFrom,
    metricsLoading,
    statData,
    statLoading,
    siteData,
    currentSiteAddress,
    healthLoading,
    toggleRFStatusLoading,
    toggleServiceLoading,
    initialNodeUptimes,
    siteUptime,
    actionOptions: SITE_ACTIONS_BUTTONS,
    getSectionName,
    handleViewChange,
    handleSwitchChange,
    handleSiteChange,
    handleActionClick,
    handleSelected: (obj: TStatusBarObj) => handleSelected(obj),
  };
}
