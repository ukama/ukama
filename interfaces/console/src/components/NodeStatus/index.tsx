/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, NodeConnectivityEnum } from '@/client/graphql/generated';
import { Button } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { styled } from '@mui/styles';
import LoadingWrapper from '../LoadingWrapper';
import NodeDropDown from '../NodeDropDown';
import SplitButton from '../SplitButton';

const StyledBtn = styled(Button)({
  whiteSpace: 'nowrap',
  minWidth: 'max-content',
});

interface INodeStatus {
  nodes: any;
  uptime: number;
  loading: boolean;
  onAddNode: Function;
  nodeActionOptions: any[];
  handleNodeSelected: Function;
  handleEditNodeClick: Function;
  selectedNode: Node | undefined;
  handleNodeActionClick: Function;
}

const NodeStatus = ({
  nodes,
  uptime,
  onAddNode,
  selectedNode,
  loading = false,
  nodeActionOptions,
  handleNodeSelected,
  handleEditNodeClick,
  handleNodeActionClick,
}: INodeStatus) => {
  const handleUpdateNode = () =>
    handleEditNodeClick(nodes.find((item: any) => item.id === selectedNode));

  return (
    <Grid container alignItems={'center'}>
      <Grid size={{ xs: 12, md: 9 }}>
        <NodeDropDown
          nodes={nodes}
          uptime={uptime}
          loading={loading}
          onAddNode={onAddNode}
          selectedNode={selectedNode}
          onNodeSelected={handleNodeSelected}
        />
      </Grid>
      <Grid
        container
        columnSpacing={2}
        size={{ xs: 12, md: 3 }}
        justifyContent="flex-end"
        visibility={
          selectedNode?.status.connectivity === NodeConnectivityEnum.Offline
            ? 'hidden'
            : 'visible'
        }
      >
        <Grid>
          <LoadingWrapper isLoading={loading} height={40}>
            <StyledBtn variant="contained" onClick={handleUpdateNode}>
              Edit NODE
            </StyledBtn>
          </LoadingWrapper>
        </Grid>

        <Grid>
          <LoadingWrapper isLoading={loading} height={40} width={'100%'}>
            <SplitButton
              options={nodeActionOptions}
              handleSplitActionClick={handleNodeActionClick}
            />
          </LoadingWrapper>
        </Grid>
      </Grid>
    </Grid>
  );
};

export default NodeStatus;
