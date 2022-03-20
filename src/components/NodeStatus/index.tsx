import { styled } from "@mui/styles";
import { NodeDto } from "../../generated";
import { Button, Stack } from "@mui/material";
import { HorizontalContainerJustify } from "../../styles";
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
    selectedNode,
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
