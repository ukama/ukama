/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  NodeConnectivityEnum,
  NodeStateEnum,
  useGetNodesByStateQuery,
  useGetSitesLazyQuery,
} from '@/client/graphql/generated';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import { NODE_TABLE_COLUMNS, NODE_TABLE_MENU } from '@/constants';
import { useAppContext } from '@/context';
import { PageContainer } from '@/styles/global';
import { TNodePoolData } from '@/types';
import { NodeEnumToString } from '@/utils';
import RouterIcon from '@mui/icons-material/Router';
import { Stack } from '@mui/material';
import { useEffect, useState } from 'react';

export default function Page() {
  const [search, setSearch] = useState<string>('');
  const [pool, setPool] = useState<TNodePoolData[]>([]);
  const [nodes, setNodes] = useState<TNodePoolData[]>([]);
  const { setSnackbarMessage, network } = useAppContext();

  const [getSites] = useGetSitesLazyQuery({
    fetchPolicy: 'cache-first',
  });

  const { data: nodesData, loading: nodesLoading } = useGetNodesByStateQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {
        connectivity: NodeConnectivityEnum.Online,
        state: NodeStateEnum.Unknown,
      },
    },
    onCompleted: async (data) => {
      if (data?.getNodesByState.nodes.length > 0) {
        const sites = await getSites({
          variables: {
            networkId: network.id,
          },
        });

        const np: TNodePoolData[] = [];
        data.getNodesByState.nodes.filter((node) => {
          const s =
            node.site.siteId &&
            sites.data?.getSites.sites.find(
              (site) => site.id === node.site.siteId,
            )?.name;
          const net =
            node.site.networkId &&
            sites.data?.getSites.sites.find(
              (site) => site.id === node.site.networkId,
            )?.name;
          np.push({
            id: node.id,
            site: s ?? '-',
            network: net ?? '-',
            type: NodeEnumToString(node.type),
            createdAt: node.site.addedAt ?? '-',
            connectivity: node.status.connectivity,
          });
          if (
            sites.data?.getSites.sites.find(
              (site) => site.id === node.site.siteId,
            )
          )
            return node;
        });
        setNodes(np);
        setPool(np);
      }
    },
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
    if (search.length > 3) {
      const _nodes: TNodePoolData[] =
        pool.filter((node) => {
          const s = search.toLowerCase();
          if (node.id.toLowerCase().includes(s)) return node;
        }) ?? [];
      setNodes(_nodes);
    } else if (search.length === 0) {
      setNodes(pool);
    }
  }, [search, nodesData?.getNodesByState.nodes]);

  const handleSearchChange = (str: string) => {
    setSearch(str);
  };

  const handleActionMenuClick = (action: string, id: string) => {
    switch (id) {
      case 'edit-node':
        break;
      case 'node-off':
        break;
      case 'restart-node':
        break;
      case 'restart-rf':
        break;
    }
  };

  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={nodesLoading}
        height={'calc(100vh - 212px)'}
        cstyle={{ marginTop: nodesLoading ? '18px' : '0px' }}
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
              subtitle={`${nodes.length}`}
            />
            <DataTableWithOptions
              dataset={nodes}
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
    </>
  );
}
