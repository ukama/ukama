/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node } from '@/generated';
import { Button, Grid } from '@mui/material';
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
  loading: boolean;
  onAddNode: Function;
  nodeActionOptions: any[];
  handleNodeSelected: Function;
  handleEditNodeClick: Function;
  selectedNode: Node | undefined;
  handleNodeActionClick: Function;
  handleNodeActionItemSelected: Function;
}

const NodeStatus = ({
  nodes,
  onAddNode,
  selectedNode,
  loading = false,
  nodeActionOptions,
  handleNodeSelected,
  handleEditNodeClick,
  handleNodeActionClick,
  handleNodeActionItemSelected,
}: INodeStatus) => {
  const handleUpdateNode = () =>
    handleEditNodeClick(nodes.find((item: any) => item.id === selectedNode));

  return (
    <Grid container>
      <Grid item xs={12} md={8}>
        <NodeDropDown
          nodes={nodes}
          loading={loading}
          onAddNode={onAddNode}
          selectedNode={selectedNode}
          onNodeSelected={handleNodeSelected}
        />
      </Grid>
      <Grid item md={4} xs={12} container spacing={2} justifyContent="flex-end">
        <Grid item>
          <LoadingWrapper isLoading={loading} height={40}>
            <StyledBtn variant="contained" onClick={handleUpdateNode}>
              Edit NODE
            </StyledBtn>
          </LoadingWrapper>
        </Grid>

        <Grid item>
          <LoadingWrapper isLoading={loading} height={40}>
            <SplitButton
              options={nodeActionOptions}
              handleSplitActionClick={handleNodeActionClick}
              handleSelectedItem={handleNodeActionItemSelected}
            />
          </LoadingWrapper>
        </Grid>
      </Grid>
    </Grid>
  );
};

export default NodeStatus;
