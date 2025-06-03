/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  NodeStateEnum,
  useGetNodesLazyQuery,
  useGetSitesQuery,
} from '@/client/graphql/generated';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import {
  NODE_ACTIONS_ENUM,
  NODE_TABLE_COLUMNS,
  NODE_TABLE_MENU,
} from '@/constants';
import { useAppContext } from '@/context';
import { PageContainer } from '@/styles/global';
import { NodeEnumToString } from '@/utils';
import RouterIcon from '@mui/icons-material/Router';
import { Stack } from '@mui/material';
import { format } from 'date-fns';
import { useEffect, useMemo, useState } from 'react';

export default function Page() {
  const [search, setSearch] = useState<string>('');
  const { setSnackbarMessage, network } = useAppContext();

  const { data: sitesData, loading: sitesLoading } = useGetSitesQuery({
    skip: !network.id,
    fetchPolicy: 'cache-first',
    variables: {
      data: {
        networkId: network.id,
      },
    },
  });

  const [getNodes, { data: nodesData, loading: nodesLoading }] =
    useGetNodesLazyQuery({
      fetchPolicy: 'cache-and-network',
      onError: (err) => {
        setSnackbarMessage({
          id: 'nodes-msg',
          message: err.message,
          type: 'error',
          show: true,
        });
      },
    });

  useEffect(() => {
    if (sitesData?.getSites?.sites) {
      getNodes({
        variables: {
          data: {
            state: NodeStateEnum.Configured,
          },
        },
      });
    }
  }, [sitesData, getNodes]);

  const getSiteName = (siteId: string | undefined | null) => {
    if (siteId === undefined || siteId === null) return '-';
    const site = sitesData?.getSites.sites.find((site) => site.id === siteId);
    return site ? site.name : '-';
  };

  const nodes = useMemo(() => {
    if (!nodesData?.getNodes.nodes) return [];
    return nodesData.getNodes.nodes
      .filter((node) => node.site.networkId === network.id)
      .map((node) => ({
        id: node.id,
        site: getSiteName(node?.site?.siteId),
        network: network.id ?? '-',
        state: node.status.state,
        type: NodeEnumToString(node.type),
        connectivity: node.status.connectivity,
        createdAt: node.site.addedAt
          ? format(new Date(node.site.addedAt), 'MM/dd/yyyy hha')
          : '-',
      }));
  }, [nodesData?.getNodes.nodes, network.id, sitesData?.getSites.sites]);

  const filteredNodes = useMemo(() => {
    if (search.length <= 3) return nodes;
    return nodes.filter((node) =>
      node.id.toLowerCase().includes(search.toLowerCase()),
    );
  }, [nodes, search]);

  const handleSearchChange = (str: string) => {
    setSearch(str);
  };

  const handleActionMenuClick = (action: string, _: string) => {
    switch (action) {
      case 'edit-node':
        break;
      case NODE_ACTIONS_ENUM.NODE_OFF:
        break;
      case NODE_ACTIONS_ENUM.NODE_RESTART:
        break;
      case NODE_ACTIONS_ENUM.NODE_RF_OFF:
        break;
    }
  };

  return (
    <LoadingWrapper
      radius="small"
      width={'100%'}
      height={'calc(100vh - 212px)'}
      isLoading={nodesLoading || sitesLoading}
      cstyle={{ marginTop: nodesLoading || sitesLoading ? '18px' : '0px' }}
    >
      <PageContainer>
        <Stack
          spacing={2}
          height={'100%'}
          direction={'column'}
          alignItems={'center'}
          justifyContent={'flex-start'}
        >
          <PageContainerHeader
            search={search}
            title={'My Nodes'}
            showSearch={true}
            onSearchChange={handleSearchChange}
            subtitle={`${filteredNodes.length}`}
          />
          <DataTableWithOptions
            dataset={filteredNodes}
            icon={RouterIcon}
            columns={NODE_TABLE_COLUMNS}
            menuOptions={NODE_TABLE_MENU}
            emptyViewLabel={'No nodes yet!'}
            onMenuItemClick={handleActionMenuClick}
            emptyViewDescription={
              'A node is the hardware piece (tower + amplifier) that connects your device to the network. Install your node, and other site components, to get started.'
            }
          />
        </Stack>
      </PageContainer>
    </LoadingWrapper>
  );
}
