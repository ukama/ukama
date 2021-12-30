import { useState } from "react";
import { Box, Grid } from "@mui/material";
import { NodesData } from "../../constants/stubData";
import { NodeDetails, NodeStatus } from "../../components";
import { NodePlaceholder, NodePlaceholderAlt } from "../../assets/images";

const Nodes = () => {
    const [selectedNodeIndex, setSelectedNodeIndex] = useState(0);

    return (
        <Box sx={{ height: "calc(100vh - 8vh)", p: "28px 0px" }}>
            <Box sx={{ flexGrow: 1, pb: "18px" }}>
                <Grid container spacing={3}>
                    <Grid xs={12} item>
                        <NodeStatus
                            nodes={NodesData}
                            selectedNodeIndex={selectedNodeIndex}
                            setSelectedNodeIndex={setSelectedNodeIndex}
                        />
                    </Grid>
                    <Grid xs={12} item container spacing={3}>
                        <NodeDetails
                            detailsList={[
                                {
                                    title: "Node Details",
                                    renderPropertyStats: false,
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
                                            value: "Home Node",
                                        },
                                        {
                                            name: "Serial #",
                                            value: "1111111111111111111",
                                        },
                                        {
                                            name: "MAC address",
                                            value: "1111111111111111111",
                                        },
                                        { name: "OS version", value: "1.0" },
                                        {
                                            name: "Manufacturing #",
                                            value: "1209391023209103",
                                        },
                                        { name: "Ukama OS", value: "1.0" },
                                        { name: "Hardware", value: "1.0" },
                                        {
                                            name: "Description",
                                            value: "Home node is a xyz.",
                                        },
                                    ],
                                },
                                {
                                    title: "Meta Data",
                                    properties: [
                                        { name: "Throughput", value: "10" },
                                        {
                                            name: "Users Attached",
                                            value: "5",
                                        },
                                    ],
                                },
                                {
                                    title: "Physical Health",
                                    image: {
                                        alt: "Node Image Alt",
                                        src: NodePlaceholderAlt,
                                    },
                                    properties: [
                                        { name: "Temperature", value: "10" },
                                        { name: "Memory", value: "5" },
                                        { name: "CPU", value: "10" },
                                        { name: "IO", value: "10" },
                                    ],
                                },
                                {
                                    title: "RF KPIs",
                                    image: {
                                        alt: "Node Image Alt",
                                        src: NodePlaceholderAlt,
                                    },
                                    properties: [
                                        { name: "QAM", value: "10" },
                                        { name: "RF Output", value: "5" },
                                        { name: "RSSI", value: "10" },
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
