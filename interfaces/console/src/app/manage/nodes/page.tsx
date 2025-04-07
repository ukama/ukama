/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { useGetNodesQuery, useGetSitesQuery } from '@/client/graphql/generated';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import SimpleDataTable from '@/components/SimpleDataTable';
import { MANAGE_NODE_POOL_COLUMN } from '@/constants';
import { useAppContext } from '@/context';
import { TNodePoolData } from '@/types';
import { NodeEnumToString } from '@/utils';
import RouterIcon from '@mui/icons-material/Router';
import { Box, Paper } from '@mui/material';
import { format } from 'date-fns';
import { useEffect, useState } from 'react';

const Page = () => {
  const { network } = useAppContext();
  const [pool, setPool] = useState<TNodePoolData[]>([]);
  const [data, setData] = useState<TNodePoolData[]>([]);
  const [search, setSearch] = useState('');
  const { setSnackbarMessage } = useAppContext();

  const { data: nodes, loading } = useGetNodesQuery({
    fetchPolicy: 'cache-and-network',
    onError: (err) => {
      setSnackbarMessage({
        id: 'available-nodes-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const { data: sites } = useGetSitesQuery({
    skip: !network?.id,
    variables: {
      networkId: network?.id,
    },
    fetchPolicy: 'cache-and-network',
    onError: (err) => {
      setSnackbarMessage({
        id: 'sites-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  useEffect(() => {
    if (nodes && nodes?.getNodes.nodes.length > 0 && sites) {
      const np: TNodePoolData[] = [];
      nodes.getNodes.nodes.filter((node) => {
        const s =
          node.site.siteId &&
          sites.getSites.sites.find((site) => site.id === node.site.siteId)
            ?.name;
        const net =
          node.site.networkId &&
          sites.getSites.sites.find((site) => site.id === node.site.networkId)
            ?.name;
        np.push({
          id: node.id,
          site: s ?? '-',
          network: net ?? '-',
          state: node.status.state,
          type: NodeEnumToString(node.type),
          connectivity: node.status.connectivity,
          createdAt: node.site.addedAt
            ? format(new Date(node.site.addedAt), 'MM/dd/yyyy hha')
            : '-',
        });
        if (sites.getSites.sites.find((site) => site.id === node.site.siteId))
          return node;
      });
      setData(np);
      setPool(np);
    }
  }, [sites, nodes]);

  useEffect(() => {
    if (search.length > 3) {
      const nodes = pool.filter((node: any) => {
        if (node.id.includes(search)) return node;
      });
      setData(nodes ?? []);
    } else if (search.length === 0) {
      setData(pool);
    }
  }, [search]);

  return (
    <LoadingWrapper
      width={'100%'}
      radius="medium"
      isLoading={loading}
      height={'calc(100vh - 244px)'}
    >
      <Paper
        sx={{
          py: { xs: 1.5, md: 3 },
          px: { xs: 2, md: 4 },
          overflow: 'scroll',
          borderRadius: '10px',
          height: '100%',
        }}
      >
        <Box sx={{ width: '100%', height: '100%' }}>
          <PageContainerHeader
            search={search}
            title={'My node pool'}
            handleButtonAction={() => {}}
            subtitle={'0'}
            onSearchChange={(e: string) => setSearch(e)}
          />
          <br />
          {data.length === 0 ? (
            <EmptyView icon={RouterIcon} title="No node in nodes pool!" />
          ) : (
            <SimpleDataTable
              dataset={data}
              isIdHyperlink={true}
              columns={MANAGE_NODE_POOL_COLUMN}
              hyperlinkPrefix={'/console/nodes'}
            />
          )}
        </Box>
      </Paper>
    </LoadingWrapper>
  );
};

export default Page;
