import { styled } from "@mui/styles";
import { Button, Grid } from "@mui/material";
import { NodeDto, Org_Node_State } from "../../generated";
import { LoadingWrapper, NodeDropDown, SplitButton } from "..";

const StyledBtn = styled(Button)({
    whiteSpace: "nowrap",
    minWidth: "max-content",
});

interface INodeStatus {
    loading: boolean;
    nodes: NodeDto[] | [];
    onAddNode: Function;
    nodeActionOptions: any[];
    onNodeSelected: Function;
    handleNodeActionClick: Function;
    onNodeActionItemSelected: Function;
    onUpdateNodeClick: Function;
    selectedNode: NodeDto | undefined;
}

const NodeStatus = ({
    nodes,
    onAddNode,
    loading = false,
    selectedNode = {
        id: "1",
        name: "",
        type: "HOME",
        totalUser: 4,
        status: Org_Node_State.Undefined,
        description: "Node 1 description",
        isUpdateAvailable: false,
        updateDescription: "",
        updateShortNote: "",
        updateVersion: "",
    },
    onNodeSelected,
    onUpdateNodeClick,
    onNodeActionItemSelected,
    nodeActionOptions,
    handleNodeActionClick,
}: INodeStatus) => {
    const handleUpdateNode = () =>
        onUpdateNodeClick(
            nodes.find((item: NodeDto) => item.id === selectedNode?.id)
        );

    return (
        <Grid container>
            <Grid item xs={12} md={9}>
                <NodeDropDown
                    nodes={nodes}
                    loading={loading}
                    onAddNode={onAddNode}
                    selectedNode={selectedNode}
                    onNodeSelected={onNodeSelected}
                />
            </Grid>
            <Grid
                item
                md={3}
                xs={12}
                container
                spacing={2}
                justifyContent="flex-end"
            >
                <Grid item>
                    <LoadingWrapper isLoading={loading} height={40}>
                        <StyledBtn
                            variant="contained"
                            onClick={handleUpdateNode}
                        >
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
