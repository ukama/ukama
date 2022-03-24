import { LoadingWrapper } from "..";
import colors from "../../theme/colors";
import DeviceModalView from "../DeviceModalView";
import { Chip, Paper, Stack, Typography, Link, Grid } from "@mui/material";

interface INodeDetailsCard {
    loading: boolean;
    nodeTitle: string;
    isUpdateAvailable: boolean;
    handleUpdateNode: Function;
    getNodeUpdateInfos: Function;
}

const NodeDetailsCard = ({
    loading,
    nodeTitle,
    isUpdateAvailable,
    getNodeUpdateInfos,
}: INodeDetailsCard) => {
    return (
        <LoadingWrapper
            width="100%"
            height="100%"
            radius={"small"}
            isLoading={loading}
        >
            <Paper sx={{ p: 3, gap: 1 }}>
                <Stack spacing={3}>
                    <Grid container spacing={1}>
                        <Grid item xs={5}>
                            <Typography variant="h6">{nodeTitle}</Typography>
                        </Grid>
                        <Grid item xs={7}>
                            {isUpdateAvailable && (
                                <Chip
                                    variant="outlined"
                                    sx={{
                                        color: colors.primaryMain,
                                        border: `1px solid ${colors.primaryMain}`,
                                    }}
                                    label={
                                        <>
                                            <Stack
                                                direction="row"
                                                alignItems="center"
                                            >
                                                <Typography variant="body2">
                                                    Software update available —
                                                    view
                                                </Typography>
                                                <Link
                                                    onClick={() =>
                                                        getNodeUpdateInfos()
                                                    }
                                                    style={{
                                                        fontSize: "14px",
                                                        paddingLeft: "3px",
                                                        cursor: "pointer",
                                                        fontStyle: "norma",
                                                        fontWeight: 400,
                                                        lineHeight: "16px",
                                                        letterSpacing:
                                                            "-0.02em",
                                                        textDecoration:
                                                            "underline",
                                                        color: colors.primaryDark,
                                                    }}
                                                >
                                                    notes
                                                </Link>
                                            </Stack>
                                        </>
                                    }
                                />
                            )}
                        </Grid>
                    </Grid>

                    <DeviceModalView />
                </Stack>
            </Paper>
        </LoadingWrapper>
    );
};

export default NodeDetailsCard;
