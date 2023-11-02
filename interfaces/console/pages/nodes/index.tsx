/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { snackbarMessage } from '@/app-recoil';
import { NODE_TABLE_COLUMNS, NODE_TABLE_MENU } from '@/constants';
import { Node, useGetNodesLazyQuery, useGetNodesQuery } from '@/generated';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { TSnackMessage } from '@/types';
import AddNodeDialog from '@/ui/molecules/AddNode';
import DataTableWithOptions from '@/ui/molecules/DataTableWithOptions';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import RouterIcon from '@mui/icons-material/Router';
import { Stack } from '@mui/material';
import { useEffect, useState } from 'react';
import { useSetRecoilState } from 'recoil';

export default function Page() {
  const [search, setSearch] = useState<string>('');
  const [nodes, setNodes] = useState<Node[] | undefined>(undefined);
  const [availableNodes, setAvailableNodes] = useState<
    Record<string, string | boolean>[] | undefined
  >(undefined);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [isShowAddNodeDialog, setIsShowAddNodeDialog] =
    useState<boolean>(false);

  const { data: nodesData, loading: nodesLoading } = useGetNodesQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {
        isFree: false,
      },
    },
    onCompleted: (data) => {
      setNodes(data.getNodes.nodes);
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

  const [getAvailableNodes, { loading: availableNodeLoading }] =
    useGetNodesLazyQuery({
      fetchPolicy: 'cache-and-network',
      onCompleted: (data) => {
        setAvailableNodes(
          data.getNodes?.nodes?.map((node) => ({
            id: node.id,
            name: node.name,
            isChecked: false,
          })),
        );
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

  useEffect(() => {
    if (search.length > 3) {
      const _nodes: Node[] =
        nodesData?.getNodes.nodes.filter((node) => {
          const s = search.toLowerCase();
          if (
            node.name.toLowerCase().includes(s) ||
            node.name.toLowerCase().includes(s)
          )
            return node;
        }) || [];
      setNodes(_nodes);
    } else if (search.length === 0) {
      setNodes(nodesData?.getNodes.nodes || []);
    }
  }, [search]);

  const handleSearchChange = (str: string) => {
    setSearch(str);
  };

  const handleAddNode = () => {};

  const handleNodeCheck = (id: string, isChecked: boolean) => {
    setAvailableNodes((prev) => {
      const nodes = prev?.map((node) => {
        if (node.id === id) {
          return { ...node, isChecked };
        }
        return node;
      });
      return nodes;
    });
  };

  const handleClaimNodeAction = () => {
    getAvailableNodes({
      variables: {
        data: {
          isFree: true,
        },
      },
    });
    setIsShowAddNodeDialog(true);
  };

  const handleCloseAddNodeDialog = () => setIsShowAddNodeDialog(false);

  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={nodesLoading}
        cstyle={{
          backgroundColor: nodesLoading ? colors.white : 'transparent',
        }}
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
              subtitle={nodes?.length ? `${nodes?.length}` : '0'}
              search={search}
              title={'My Nodes'}
              showSearch={true}
              buttonTitle="Add Nodes"
              onSearchChange={handleSearchChange}
              handleButtonAction={handleClaimNodeAction}
            />
            <DataTableWithOptions
              dataset={nodes || []}
              icon={RouterIcon}
              onMenuItemClick={() => {}}
              columns={NODE_TABLE_COLUMNS}
              menuOptions={NODE_TABLE_MENU}
              emptyViewLabel={'No node yet!'}
            />
          </Stack>
        </PageContainer>
      </LoadingWrapper>
      <AddNodeDialog
        data={availableNodes}
        labelNegativeBtn="Close"
        labelSuccessBtn="Add Nodes"
        isOpen={isShowAddNodeDialog}
        handleNodeCheck={handleNodeCheck}
        title="Add nodes to your network"
        handleSuccessAction={handleAddNode}
        handleCloseAction={handleCloseAddNodeDialog}
        description="Add nodes to your network to start managing them here. If you cannot find a desired node, check to make sure it is not already added to another network."
      />
    </>
  );
}
