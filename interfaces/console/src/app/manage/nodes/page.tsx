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
import { NodeEnumToString } from '@/utils';
import RouterIcon from '@mui/icons-material/Router';
import { Paper, Stack } from '@mui/material';
import { format } from 'date-fns';
import { useCallback, useMemo, useState } from 'react';

const Page = () => {
  const { network, setSnackbarMessage } = useAppContext();
  const [search, setSearch] = useState('');

  const { data: nodes, loading } = useGetNodesQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {},
    },
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
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {},
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'sites-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const getSiteName = useCallback(
    (siteId: string | undefined | null) => {
      if (!siteId || !sites?.getSites.sites) return '-';
      const site = sites.getSites.sites.find((site) => site.id === siteId);
      return site ? site.name : '-';
    },
    [sites],
  );

  const transformedData = useMemo(() => {
    if (!nodes?.getNodes.nodes || !sites?.getSites.sites) {
      return [];
    }

    return nodes.getNodes.nodes.map((node) => ({
      id: node.id,
      network: node.site.networkId
        ? (sites.getSites.sites.find((site) => site.id === node.site.networkId)
          ?.name ?? '-')
        : '-',
      state: node.status.state,
      site: getSiteName(node.site.siteId),
      type: NodeEnumToString(node.type),
      connectivity: node.status.connectivity,
      createdAt: node.site.addedAt
        ? format(new Date(node.site.addedAt), 'MM/dd/yyyy hh:mm a')
        : '-',
    }));
  }, [nodes?.getNodes.nodes, sites?.getSites.sites, getSiteName]);

  const filteredData = useMemo(() => {
    if (search.length > 3) {
      return transformedData.filter((node) => node.id.includes(search));
    }
    return transformedData;
  }, [transformedData, search]);

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
        <Stack sx={{ width: '100%', height: '100%' }} spacing={4}>
          <PageContainerHeader
            search={search}
            title={'My node pool'}
            handleButtonAction={() => {}}
            subtitle={filteredData.length.toString()}
            onSearchChange={(e: string) => setSearch(e)}
          />
          {filteredData.length === 0 ? (
            <EmptyView icon={RouterIcon} title="No node in nodes pool!" />
          ) : (
            <SimpleDataTable
              dataset={filteredData}
              isIdHyperlink={true}
              columns={MANAGE_NODE_POOL_COLUMN}
              hyperlinkPrefix={'/console/nodes/'}
            />
          )}
        </Stack>
      </Paper>
    </LoadingWrapper>
  );
};

export default Page;
