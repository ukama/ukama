import { useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { Box, Grid, Stack } from "@mui/material";
import { isSkeltonLoading, organizationId } from "../../recoil";
import {
    NodeDetailsCard,
    NodeStatus,
    NodeInfoCard,
    LoadingWrapper,
    PagePlaceholder,
} from "../../components";
import {
    NODES,
    NODE_PROPERTIES8,
    NODE_PROPERTIES2,
    NODE_PROPERTIES4,
    NODE_PROPERTIES3,
} from "../../constants/stubData";
import {
    NodeDto,
    useGetNodesByOrgQuery,
    Org_Node_State,
} from "../../generated";

const Nodes = () => {
    const [selectedNode, setSelectedNode] = useState<NodeDto>({
        id: "temp",
        totalUser: 3,
        title: "Node 1",
        description: "testing",
        status: Org_Node_State.Undefined,
    });
    const [selectedTab, setSelectedTab] = useState(1);
    const orgId = useRecoilValue(organizationId);
    const skeltonLoading = useRecoilValue(isSkeltonLoading);
    const { data: nodesRes, loading: nodesLoading } = useGetNodesByOrgQuery({
        variables: { orgId: orgId || "" },
    });
    // const { data: nodeDetailsRes, loading: nodeDetailsResLoading } =
    //     useGetNodeDetailsQuery();

    const onTabSelected = (value: number) => setSelectedTab(value);
    const onNodeSelected = (node: NodeDto) => setSelectedNode(node);
    const onNodeRFClick = () => {
        //TODO: Handle NODE RF ACTIONS
    };
    const onNodeSwitchClick = () => {
        //TODO: Handle NODE ON/OFF ACTIONS
    };
    const onRestartNodeClick = () => {
        //TODO: Handle NODE RESTART ACTION
    };

    const isLoading = skeltonLoading || nodesLoading;

    if (nodesRes && nodesRes?.getNodesByOrg?.nodes?.length === 0)
        return (
            <RoundedCard
                sx={{
                    p: 0,
                    mt: 3,
                    mb: 2,
                    borderRadius: "4px",
                    height: "calc(100% - 15%)",
                }}
            >
                <PagePlaceholder description="Order your node now." />
            </RoundedCard>
        );

    return (
        <Box
            sx={{
                p: 0,
                mt: 3,
                pb: 2,
            }}
        >
            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <NodeStatus
                        selectedNode={selectedNode}
                        onNodeRFClick={onNodeRFClick}
                        onNodeSelected={onNodeSelected}
                        onNodeSwitchClick={onNodeSwitchClick}
                        onRestartNodeClick={onRestartNodeClick}
                        nodes={nodesRes?.getNodesByOrg?.nodes || NODES}
                    />
                </Grid>
                <Grid item container xs={4}>
                    <Stack spacing={2} sx={{ width: "100%" }}>
                        <NodeInfoCard
                            index={1}
                            loading={isLoading}
                            title={"Node Detail"}
                            properties={NODE_PROPERTIES8}
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 1}
                        />
                        <NodeInfoCard
                            index={2}
                            loading={isLoading}
                            title={"Meta Data"}
                            properties={NODE_PROPERTIES2}
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 2}
                        />
                        <NodeInfoCard
                            index={3}
                            loading={isLoading}
                            title={"Physical Health"}
                            properties={NODE_PROPERTIES4}
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 3}
                        />
                        <NodeInfoCard
                            index={4}
                            loading={isLoading}
                            title={"RF KPIs"}
                            properties={NODE_PROPERTIES3}
                            onSelected={onTabSelected}
                            isSelected={selectedTab === 4}
                        />
                    </Stack>
                </Grid>
                <Grid item container xs={8}>
                    <RoundedCard sx={{ borderRadius: "4px" }}>
                        <LoadingWrapper
                            height={"70%"}
                            radius={"small"}
                            isLoading={isLoading}
                        >
                            {selectedTab === 1 && <NodeDetailsCard />}
                        </LoadingWrapper>
                    </RoundedCard>
                </Grid>
            </Grid>
        </Box>
    );
};

export default Nodes;
