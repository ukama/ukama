import { styled } from "@mui/styles";
import { Button, Stack } from "@mui/material";
import { HorizontalContainerJustify } from "../../styles";
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
        title: "",
        type: "HOME",
        totalUser: 4,
        status: Org_Node_State.Undefined,
        description: "Node 1 description",
    },
    onNodeSelected,
    onUpdateNodeClick,
    onNodeActionItemSelected,
    nodeActionOptions,
    handleNodeActionClick,
}: INodeStatus) => {
    const handleUpdateNode = () =>
        onUpdateNodeClick(
            nodes.find((item: NodeDto) => item.id === selectedNode.id)
        );

    return (
        <HorizontalContainerJustify>
            <NodeDropDown
                nodes={nodes}
                loading={loading}
                onAddNode={onAddNode}
                selectedNode={selectedNode}
                onNodeSelected={onNodeSelected}
            />
            <Stack direction={"row"} spacing={2}>
                <LoadingWrapper isLoading={loading} height={40} width={100}>
                    <StyledBtn variant="contained" onClick={handleUpdateNode}>
                        UPDATE NODE
                    </StyledBtn>
                </LoadingWrapper>
                <LoadingWrapper isLoading={loading} height={40} width={100}>
                    <SplitButton
                        options={nodeActionOptions}
                        handleSplitActionClick={handleNodeActionClick}
                        handleSelectedItem={onNodeActionItemSelected}
                    />
                </LoadingWrapper>
            </Stack>
        </HorizontalContainerJustify>
    );
};

export default NodeStatus;
