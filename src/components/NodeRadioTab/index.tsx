import { TooltipsText } from "../../constants";
import { Stack, Paper, Typography } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, StackedAreaChart } from "..";

interface INodeRadioTab {
    loading: boolean;
}
const NodeRadioTab = ({ loading }: INodeRadioTab) => {
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
                    value={"NNN"}
                    variant={"large"}
                    name={"TX Power"}
                    nameInfo={TooltipsText.TXPOWER}
                />
                <NodeStatItem
                    value={"NNN"}
                    variant={"large"}
                    name={"RX Power"}
                    nameInfo={TooltipsText.RXPOWER}
                />
                <NodeStatItem
                    value={"NNN"}
                    name={"PA Power"}
                    variant={"large"}
                    nameInfo={TooltipsText.PAPOWER}
                />
            </NodeStatsContainer>
            <Paper sx={{ p: 3, width: "100%" }}>
                <Typography variant="h6">Radio</Typography>
                <Stack spacing={6} pt={2}>
                    <StackedAreaChart hasData={true} title={"TX Power"} />
                    <StackedAreaChart hasData={true} title={"RX Power "} />
                    <StackedAreaChart hasData={true} title={"PA Power "} />
                </Stack>
            </Paper>
        </Stack>
    );
};

export default NodeRadioTab;
