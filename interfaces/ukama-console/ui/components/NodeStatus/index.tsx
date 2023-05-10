import { GetNodeStatusRes, NodeDto, Org_Node_State } from '@/generated';
import { Button, Grid } from '@mui/material';
import { styled } from '@mui/styles';
import { LoadingWrapper, NodeDropDown, SplitButton } from '..';

const StyledBtn = styled(Button)({
  whiteSpace: 'nowrap',
  minWidth: 'max-content',
});

interface INodeStatus {
  loading: boolean;
  nodes: NodeDto[] | [];
  onAddNode: Function;
  nodeActionOptions: any[];
  onNodeSelected: Function;
  nodeStatusLoading: boolean;
  onUpdateNodeClick: Function;
  handleNodeActionClick: Function;
  onNodeActionItemSelected: Function;
  selectedNode: NodeDto | undefined;
  nodeStatus: GetNodeStatusRes | undefined;
}

const NodeStatus = ({
  nodes,
  onAddNode,
  nodeStatus,
  loading = false,
  selectedNode = {
    id: '1',
    name: '',
    type: 'HOME',
    totalUser: 4,
    status: Org_Node_State.Undefined,
    description: 'Node 1 description',
    isUpdateAvailable: false,
    updateDescription: '',
    updateShortNote: '',
    updateVersion: '',
  },
  onNodeSelected,
  onUpdateNodeClick,
  nodeActionOptions,
  nodeStatusLoading,
  handleNodeActionClick,
  onNodeActionItemSelected,
}: INodeStatus) => {
  const handleUpdateNode = () =>
    onUpdateNodeClick(
      nodes.find((item: NodeDto) => item.id === selectedNode?.id),
    );

  return (
    <Grid container>
      <Grid item xs={12} md={8}>
        <NodeDropDown
          nodes={nodes}
          loading={loading}
          onAddNode={onAddNode}
          nodeStatus={nodeStatus}
          selectedNode={selectedNode}
          onNodeSelected={onNodeSelected}
          nodeStatusLoading={nodeStatusLoading}
        />
      </Grid>
      <Grid item md={4} xs={12} container spacing={2} justifyContent="flex-end">
        <Grid item>
          <LoadingWrapper isLoading={loading} height={40}>
            <StyledBtn variant="contained" onClick={handleUpdateNode}>
              UPDATE NODE
            </StyledBtn>
          </LoadingWrapper>
        </Grid>

        <Grid item>
          <LoadingWrapper isLoading={loading} height={40}>
            <SplitButton
              options={nodeActionOptions}
              handleSplitActionClick={handleNodeActionClick}
              handleSelectedItem={onNodeActionItemSelected}
            />
          </LoadingWrapper>
        </Grid>
      </Grid>
    </Grid>
  );
};

export default NodeStatus;
