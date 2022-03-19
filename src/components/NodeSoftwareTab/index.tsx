import { LoadingWrapper } from "..";
import {
    Paper,
    Typography,
    Card,
    Button,
    Stack,
    Grid,
    CardContent,
    CardActions,
} from "@mui/material";
import { SimpleDataTable } from "../../components";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { colors } from "../../theme";
import { NodeAppsColumns } from "../../constants/tableColumns";
interface INodeRadioTab {
    loading: boolean;
    nodeApps: any[];
    NodeLogs: any;
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
                <Grid container spacing={3}>
                    {nodeApps?.map(
                        ({ id, nodeAppName, cpu, memory, version }: any) => (
                            <Grid item xs={12} md={6} lg={3} key={id}>
                                <Card variant="outlined">
                                    <CardContent>
                                        <Stack
                                            direction="row"
                                            sx={{ alignItems: "center" }}
                                            spacing={1}
                                        >
                                            <CheckCircleIcon
                                                htmlColor={colors.green}
                                                sx={{
                                                    position: "relative",
                                                    left: -3,
                                                    fontSize: "23px",
                                                }}
                                            />
                                            <Typography variant="h5">
                                                {nodeAppName}
                                            </Typography>
                                        </Stack>
                                        <Typography
                                            variant="body2"
                                            color="text.secondary"
                                            gutterBottom
                                        >
                                            Version: {version}
                                        </Typography>
                                        <Stack direction="row" spacing={1 / 2}>
                                            <Typography variant="body2">
                                                CPU:
                                            </Typography>
                                            <Typography
                                                variant="body2"
                                                sx={{ color: colors.darkBlue }}
                                            >
                                                {cpu}%
                                            </Typography>
                                        </Stack>
                                        <Stack direction="row" spacing={1 / 2}>
                                            <Typography variant="body2">
                                                MEMORY:
                                            </Typography>
                                            <Typography
                                                variant="body2"
                                                sx={{ color: colors.darkBlue }}
                                            >
                                                {memory} KB
                                            </Typography>
                                        </Stack>
                                    </CardContent>
                                    <CardActions sx={{ ml: 1 }}>
                                        <Button
                                            onClick={() =>
                                                getNodeAppDetails(id)
                                            }
                                        >
                                            VIEW MORE
                                        </Button>
                                    </CardActions>
                                </Card>
                            </Grid>
                        )
                    )}
                </Grid>
            </Paper>
        </LoadingWrapper>
    );
};

export default NodeSoftwareTab;
