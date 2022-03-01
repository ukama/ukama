import { LoadingWrapper } from "..";
import { Paper, Typography, Stack, Button, Grid } from "@mui/material";
import { SimpleDataTable } from "../../components";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { colors } from "../../theme";
import { SimpleCardWithBorder } from "../../styles";
import { NodeAppsColumns } from "../../constants/tableColumns";
interface INodeRadioTab {
    loading: boolean;
    nodeApps: any[];
    NodeLogs: any[];
    getNodeAppDetails: Function;
}

const NodeSoftwareTab = ({
    getNodeAppDetails,
    loading,
    NodeLogs,
    nodeApps,
}: INodeRadioTab) => {
    return (
        <LoadingWrapper isLoading={loading} height={400}>
            <Paper
                sx={{
                    p: 3,
                    height: "100%",
                    borderRadius: "4px",
                    marginBottom: 2,
                }}
            >
                <Typography variant="h6" sx={{ marginBottom: 3 }}>
                    Change Logs
                </Typography>
                <SimpleDataTable columns={NodeAppsColumns} dataset={NodeLogs} />
            </Paper>
            <Paper sx={{ height: "100%", p: 3, borderRadius: "4px" }}>
                <Typography variant="h6" sx={{ mb: 4 }}>
                    Node Apps
                </Typography>
                <Grid container spacing={3} sx={{ p: 0 }}>
                    {nodeApps?.map(
                        ({ id, nodeAppName, cpu, memory, version }: any) => (
                            <Grid item xs={12} md={6} lg={3} key={id}>
                                <SimpleCardWithBorder>
                                    <Stack
                                        direction="column"
                                        justifyContent="flex-start"
                                        sx={{ p: 1 }}
                                    >
                                        <Stack
                                            direction="row"
                                            sx={{ alignItems: "center" }}
                                        >
                                            <CheckCircleIcon
                                                htmlColor={colors.green}
                                                sx={{
                                                    position: "relative",
                                                    left: -3,
                                                }}
                                            />
                                            <Typography variant="h5">
                                                {nodeAppName}
                                            </Typography>
                                        </Stack>
                                        <Typography
                                            variant="body2"
                                            sx={{
                                                color: colors.black70,
                                                mb: 1,
                                            }}
                                        >
                                            version: {version}
                                        </Typography>

                                        <Typography variant="body2">
                                            CPU: {cpu}%
                                        </Typography>
                                        <Typography variant="body2">
                                            memory: {memory} KB
                                        </Typography>
                                        <Button
                                            sx={{
                                                justifyContent: "flex-start",
                                                mt: 1,
                                            }}
                                            onClick={() =>
                                                getNodeAppDetails(id)
                                            }
                                        >
                                            VIEW MORE
                                        </Button>
                                    </Stack>
                                </SimpleCardWithBorder>
                            </Grid>
                        )
                    )}
                </Grid>
            </Paper>
        </LoadingWrapper>
    );
};

export default NodeSoftwareTab;
