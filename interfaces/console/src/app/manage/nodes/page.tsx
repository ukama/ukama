/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { NetworkDto, useGetNetworksQuery } from '@/client/graphql/generated';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import SimpleDataTable from '@/components/SimpleDataTable';
import { MANAGE_NODE_POOL_COLUMN } from '@/constants';
import { useAppContext } from '@/context';
import RouterIcon from '@mui/icons-material/Router';
import { AlertColor, Box, Paper } from '@mui/material';
import { useState } from 'react';

const Page = () => {
  // const [data, setData] = useState([]);
  const [search, setSearch] = useState('');
  const { setSnackbarMessage } = useAppContext();
  const [networkList, setNetworkList] = useState<NetworkDto[]>([]);

  // const [getNodes, { loading: getNodesLoading }] = useGetNodesLazyQuery({
  //   fetchPolicy: 'cache-and-network',
  //   onCompleted: (data) => {
  //     const filteredNodes = data?.getNodes.nodes;
  // .filter((node) => node.created_at)
  // .map((node) => ({
  //   ...node,
  //   created_at: format(parseISO(node.created_at), 'dd MMM yyyy'),
  // }));

  //     setData((prev: any) => ({
  //       ...prev,
  //       node: filteredNodes ?? [],
  //     }));
  //   },
  //   onError: (error) => {
  //     setSnackbarMessage({
  //       id: 'node',
  //       message: error.message,
  //       type: 'error' as AlertColor,
  //       show: true,
  //     });
  //   },
  // });

  const { loading: networkLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setNetworkList(data?.getNetworks.networks ?? []);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'network',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  // useEffect(() => {
  //   if (nodeSearch.length > 3) {
  //     const nodes = data.node.filter((node: any) => {
  //       if (node.id.includes(nodeSearch)) return node;
  //     });
  //     setData((prev: any) => ({ ...prev, node: nodes ?? [] }));
  //   } else if (nodeSearch.length === 0) {
  //     setData((prev: any) => ({ ...prev, node: data.node }));
  //   }
  // }, [nodeSearch]);

  const handleCreateNetwork = () => {};

  return (
    <LoadingWrapper
      width={'100%'}
      radius="medium"
      isLoading={networkLoading}
      height={'calc(100vh - 400px)'}
    >
      <Paper
        sx={{
          py: 3,
          px: 4,
          overflow: 'scroll',
          borderRadius: '10px',
          height: 'calc(100vh - 400px)',
        }}
      >
        <Box sx={{ width: '100%', height: '100%' }}>
          <PageContainerHeader
            search={search}
            title={'My node pool'}
            buttonTitle={'CLAIM NODE'}
            handleButtonAction={() => {}}
            subtitle={'0'}
            onSearchChange={(e: string) => setSearch(e)}
          />
          <br />
          {[].length === 0 ? (
            <EmptyView icon={RouterIcon} title="No node in nodes pool!" />
          ) : (
            <SimpleDataTable
              dataset={[]}
              networkList={networkList}
              columns={MANAGE_NODE_POOL_COLUMN}
              handleCreateNetwork={handleCreateNetwork}
            />
          )}
        </Box>
      </Paper>
    </LoadingWrapper>
  );
};

export default Page;
