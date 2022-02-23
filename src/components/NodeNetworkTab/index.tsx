import { TooltipsText } from "../../constants";
import { Stack, Paper, Typography } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, LineChart } from "..";

interface INodeOverviewTab {
    loading: boolean;
}
const NodeNetworkTab = ({ loading }: INodeOverviewTab) => {
    return (
        <Stack spacing={2} direction="row">
            <NodeStatsContainer
                index={0}
                selected={0}
                loading={loading}
                title={"Network"}
                isCollapsable={true}
            >
                <NodeStatItem
                    variant={"large"}
                    value={"100 Mbps"}
                    name={"Throughput (D/L)"}
                    nameInfo={TooltipsText.DL}
                />
                <NodeStatItem
                    variant={"large"}
                    value={"400 Mbps"}
                    name={"Throughput (U/L)"}
                    nameInfo={TooltipsText.UL}
                />
                <NodeStatItem
                    value={"80%"}
                    variant={"large"}
                    name={"RRC CNX Success"}
                    nameInfo={TooltipsText.RRCCNX}
                />
                <NodeStatItem
                    value={"72%"}
                    variant={"large"}
                    name={"ERAB Drop Rate"}
                    nameInfo={TooltipsText.ERAB}
                />
                <NodeStatItem
                    value={"60%"}
                    variant={"large"}
                    name={"RLS  Drop Rate"}
                    nameInfo={TooltipsText.RLS}
                />
            </NodeStatsContainer>

            <Paper sx={{ p: 3, width: "100%" }}>
                <Typography variant="h6">Network</Typography>

                <Stack spacing={6} pt={2}>
                    <LineChart hasData={true} title={"Throughput (U/L)"} />
                    <LineChart hasData={true} title={"Throughput (D/L)"} />
                    <LineChart hasData={true} title={"RRC CNX Success "} />
                    <LineChart hasData={true} title={"ERAB Drop Rate"} />
                    <LineChart hasData={true} title={"RLS  Drop Rate"} />
                </Stack>
            </Paper>
        </Stack>
    );
};

export default NodeNetworkTab;
