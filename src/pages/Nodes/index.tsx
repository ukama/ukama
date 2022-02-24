import {
    NodeDto,
    useGetMetricsCpuTrxQuery,
    useGetNodeDetailsQuery,
    useGetNodesByOrgQuery,
} from "../../generated";
import {
    TabPanel,
    NodeStatus,
    NodeRadioTab,
    LoadingWrapper,
    NodeNetworkTab,
    NodeSoftwareTab,
    NodeOverviewTab,
    PagePlaceholder,
    NodeResourcesTab,
    NodeAppDetailsDialog,
} from "../../components";
import React, { useState } from "react";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading } from "../../recoil";
import { getGraphTimespan, parseObjectInNameValue } from "../../utils";
import { Box, Grid, Paper, Tab, Tabs } from "@mui/material";
import {
    NodePageTabs,
    NodeApps,
    NodeAppLogs,
    NODE_ACTIONS,
} from "../../constants";

const Nodes = () => {
    const [selectedTab, setSelectedTab] = useState(0);
    const skeltonLoading = useRecoilValue(isSkeltonLoading);
    const [selectedNode, setSelectedNode] = useState<NodeDto>();
    const [showNodeAppDialog, setShowNodeAppDialog] = useState(false);
    const [cpuTrxLastTimespan, setCpuTrxLastTimespan] = useState<number | null>(
        null
    );

    const getMetricsPayload = (orgId: string, nodeId: string) => {
        return {
            orgId: orgId,
            nodeId: nodeId,
            ...getGraphTimespan(cpuTrxLastTimespan),
        };
    };

    const { data: nodesRes, loading: nodesLoading } = useGetNodesByOrgQuery({
        variables: { orgId: "1" || "" },
        onCompleted: res => {
            res.getNodesByOrg.nodes.length > 0 &&
                setSelectedNode(res.getNodesByOrg.nodes[0]);
        },
    });

    const { data: nodeDetailRes, loading: nodeDetailLoading } =
        useGetNodeDetailsQuery();

    const { data: nodeCpuTrxRes } = useGetMetricsCpuTrxQuery({
        variables: {
            data: getMetricsPayload(
                "a32485e4-d842-45da-bf3e-798889c68ad0",
                "uk-test36-hnode-a1-30df"
            ),
        },
        // pollInterval: 10000, //120000,
        onCompleted: res => {
            setCpuTrxLastTimespan(
                res.getMetricsCpuTRX[res.getMetricsCpuTRX.length - 1].timestamp
            );
        },
    });

    const onTabSelected = (event: React.SyntheticEvent, value: any) =>
        setSelectedTab(value);
    const onNodeSelected = (node: NodeDto) => setSelectedNode(node);

    const onUpdateNodeClick = () => {
        //TODO: Handle NODE RESTART ACTION
    };
    const handleNodeActionItemSelected = () => {
        //Todo :Handle nodeAction Itemselected
    };

    const handleNodeActioOptionClicked = () => {
        //Todo :Handle nodeAction selected and clicked
    };
    const onAddNode = () => {
        //TODO: Handle NODE ADD ACTION
    };
    const handleUpdateNode = () => {
        // TODO: Handle Update node Action
    };

    const getNodeDetails = () => {
        //TODO:Handle nodeDetails
        setShowNodeAppDialog(true);
    };
    const handleNodAppDetailsDialog = () => {
        setShowNodeAppDialog(false);
    };

    const isLoading = skeltonLoading || nodesLoading;

    return (
        <Box
            component="div"
            sx={{
                p: 0,
                mt: 3,
                pb: 2,
            }}
        >
            {nodesRes || isLoading ? (
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <NodeStatus
                            onAddNode={onAddNode}
                            loading={nodesLoading}
                            handleNodeActionClick={handleNodeActioOptionClicked}
                            selectedNode={selectedNode}
                            onNodeActionItemSelected={
                                handleNodeActionItemSelected
                            }
                            onNodeSelected={onNodeSelected}
                            nodeActionOptions={NODE_ACTIONS}
                            onUpdateNodeClick={onUpdateNodeClick}
                            nodes={nodesRes?.getNodesByOrg?.nodes || []}
                        />
                    </Grid>
                    <Grid item xs={12}>
                        <LoadingWrapper isLoading={isLoading} height={"40px"}>
                            <Tabs value={selectedTab} onChange={onTabSelected}>
                                {NodePageTabs.map(({ id, label, value }) => (
                                    <Tab
                                        key={id}
                                        label={label}
                                        id={`node-tab-${value}`}
                                    />
                                ))}
                            </Tabs>
                        </LoadingWrapper>
                    </Grid>

                    <Grid item xs={12}>
                        <TabPanel
                            id={"node-tab-0"}
                            value={selectedTab}
                            index={0}
                        >
                            <NodeOverviewTab
                                isUpdateAvailable={true}
                                selectedNode={selectedNode}
                                handleUpdateNode={handleUpdateNode}
                                loading={isLoading || nodeDetailLoading}
                                nodeDetails={parseObjectInNameValue(
                                    nodeDetailRes?.getNodeDetails
                                )}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-1"}
                            value={selectedTab}
                            index={1}
                        >
                            <NodeNetworkTab
                                loading={isLoading || nodeDetailLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-2"}
                            value={selectedTab}
                            index={2}
                        >
                            <NodeResourcesTab
                                loading={isLoading || nodeDetailLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-3"}
                            value={selectedTab}
                            index={3}
                        >
                            <NodeRadioTab
                                loading={isLoading || nodeDetailLoading}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-4"}
                            value={selectedTab}
                            index={4}
                        >
                            <NodeSoftwareTab
                                loading={isLoading || nodeDetailLoading}
                                nodeApps={NodeApps}
                                NodeLogs={NodeAppLogs}
                                getNodeAppDetails={getNodeDetails}
                            />
                        </TabPanel>
                        <TabPanel
                            id={"node-tab-5"}
                            value={selectedTab}
                            index={5}
                        >
                            <Paper>Schematic</Paper>
                        </TabPanel>
                    </Grid>
                </Grid>
            ) : (
                <PagePlaceholder
                    hyperlink="#"
                    linkText={"here"}
                    showActionButton={false}
                    buttonTitle="Install sims"
                    description="Your nodes have not arrived yet. View their status"
                />
            )}
            <NodeAppDetailsDialog
                closeBtnLabel="close"
                isOpen={showNodeAppDialog}
                handleClose={handleNodAppDetailsDialog}
            />
        </Box>
    );
};

export default Nodes;
