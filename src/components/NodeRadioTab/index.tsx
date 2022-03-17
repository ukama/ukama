import { useState } from "react";
import { TooltipsText } from "../../constants";
import { Paper, Grid, Typography, Stack } from "@mui/material";
import { NodeStatsContainer, NodeStatItem, StackedAreaChart } from "..";
interface INodeRadioTab {
    loading: boolean;
}
const NodeRadioTab = ({ loading }: INodeRadioTab) => {
    const [isCollapse, setIsCollapse] = useState<boolean>(false);
    const handleCollapse = () => setIsCollapse(prev => !prev);
    return (
        <Grid container spacing={3}>
            <Grid item lg={!isCollapse ? 3 : 1} md xs>
                <NodeStatsContainer
                    index={0}
                    selected={0}
                    title={"Radio"}
                    loading={loading}
                    isCollapsable={true}
                    isCollapse={isCollapse}
                    onCollapse={handleCollapse}
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
            </Grid>
            <Grid item lg={isCollapse ? 11 : 9} md xs>
                <Paper sx={{ padding: "22px 18px 0px 30px", width: "100%" }}>
                    <Stack spacing={1}>
                        <Typography variant="h6">Radio</Typography>
                        <StackedAreaChart hasData={true} title={"TX Power"} />
                        <StackedAreaChart hasData={true} title={"RX Power "} />
                        <StackedAreaChart hasData={true} title={"PA Power "} />
                    </Stack>
                </Paper>
            </Grid>
        </Grid>
    );
};

export default NodeRadioTab;
