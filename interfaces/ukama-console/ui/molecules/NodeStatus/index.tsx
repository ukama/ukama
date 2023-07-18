import { GetNodeStatusRes, NodeDto } from '@/generated';
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
  loading: boolean;
  onAddNode: Function;
  nodes: any;
  nodeActionOptions: any[];
  nodeStatusLoading: boolean;
  handleNodeSelected: Function;
  handleEditNodeClick: Function;
  handleNodeActionClick: Function;
  selectedNode: NodeDto | undefined;
  handleNodeActionItemSelected: Function;
  nodeStatus: GetNodeStatusRes | undefined;
}

const NodeStatus = ({
  nodes,
  onAddNode,
  nodeStatus,
  selectedNode,
  loading = false,
  nodeActionOptions,
  nodeStatusLoading,
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
          nodeStatus={nodeStatus}
          selectedNode={selectedNode}
          onNodeSelected={handleNodeSelected}
          nodeStatusLoading={nodeStatusLoading}
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
