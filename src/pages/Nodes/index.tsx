import { useState } from "react";
import { Box, Grid } from "@mui/material";
import { NodeDetails, NodeStatus, LoadingWrapper } from "../../components";
import { useGetNodesByOrgQuery, useGetNodeDetailsQuery } from "../../generated";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading, organizationId } from "../../recoil";
import { NodePlaceholder } from "../../assets/images";

const Nodes = () => {
    const [selectedNodeIndex, setSelectedNodeIndex] = useState(0);
    const orgId = useRecoilValue(organizationId);
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const { data: nodesRes, loading: nodesLoading } = useGetNodesByOrgQuery({
        variables: { orgId: orgId || "" },
    });
    const { data: nodeDetailsRes, loading: nodeDetailsResLoading } =
        useGetNodeDetailsQuery();

    return (
        <Box sx={{ height: "calc(100vh - 8vh)" }}>
            <Box sx={{ flexGrow: 1 }}>
                <Grid container spacing={3}>
                    <Grid xs={12} item>
                        <LoadingWrapper
                            width="100%"
                            height="30px"
                            isLoading={isSkeltonLoad || nodesLoading}
                        >
                            <NodeStatus
                                nodes={nodesRes?.getNodesByOrg?.nodes}
                                selectedNodeIndex={selectedNodeIndex}
                                setSelectedNodeIndex={setSelectedNodeIndex}
                            />
                        </LoadingWrapper>
                    </Grid>
                    <Grid xs={12} item container spacing={3}>
                        <NodeDetails
                            detailsList={[
                                {
                                    loading: nodeDetailsResLoading,
                                    title: "Node Details",

                                    image: {
                                        alt: "Node Image",
                                        src: NodePlaceholder,
                                    },
                                    button: {
                                        label: "view diagnostics",
                                        onClick: () => {
                                            return;
                                        },
                                    },
                                    properties: [
                                        {
                                            name: "Model type",
                                            value: `${nodeDetailsRes?.getNodeDetails.modelType}`,
                                        },
                                        {
                                            name: "Serial#",
                                            value: `${nodeDetailsRes?.getNodeDetails.serial}`,
                                        },
                                        {
                                            name: "MAC address",
                                            value: `${nodeDetailsRes?.getNodeDetails.macAddress}`,
                                        },
                                        {
                                            name: "OS version",
                                            value: `${nodeDetailsRes?.getNodeDetails.osVersion}`,
                                        },
                                        {
                                            name: "Manufacturing#",
                                            value: `${nodeDetailsRes?.getNodeDetails.manufacturing}`,
                                        },
                                        {
                                            name: "Ukama OS",
                                            value: `${nodeDetailsRes?.getNodeDetails.ukamaOS}`,
                                        },
                                        {
                                            name: "Hardware",
                                            value: `${nodeDetailsRes?.getNodeDetails.hardware}`,
                                        },
                                        {
                                            name: "Description",
                                            value: `${nodeDetailsRes?.getNodeDetails.description}`,
                                        },
                                    ],
                                },
                                {
                                    title: "Meta Data",

                                    properties: [
                                        {
                                            name: "Throughput",
                                            value: `i3iopwiopiwpoe`,
                                        },
                                        {
                                            name: "Users Attached",
                                            value: 12,
                                        },
                                    ],
                                },
                                {
                                    title: "Physical Health",

                                    image: {
                                        alt: "Node Image Alt",
                                        src: "NodePlaceholderAlt",
                                    },
                                    properties: [
                                        {
                                            name: "Temperature",
                                            value: `98398`,
                                        },
                                        {
                                            name: "Memory",
                                            value: `92839`,
                                        },
                                        {
                                            name: "CPU",
                                            value: `392`,
                                        },
                                        {
                                            name: "IO",
                                            value: `28`,
                                        },
                                    ],
                                },
                                {
                                    title: "RF KPIs",
                                    image: {
                                        alt: "Node Image Alt",
                                        src: "NodePlaceholderAlt",
                                    },
                                    properties: [
                                        {
                                            name: "QAM",
                                            value: `343`,
                                        },
                                        {
                                            name: "RF Output",
                                            value: `32`,
                                        },
                                        {
                                            name: "RSSI",
                                            value: `23`,
                                        },
                                    ],
                                },
                            ]}
                        />
                    </Grid>
                </Grid>
            </Box>
        </Box>
    );
};

export default Nodes;
