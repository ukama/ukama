import {
    CpuIcon,
    NodeImg,
    UsersIcon,
    HistoryIcon,
    BatteryIcon,
    ThermometerIcon,
} from "../../assets/svg";
import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import OptionsPopover from "../OptionsPopover";
import { BASIC_MENU_ACTIONS } from "../../constants";
import { Paper, Typography, Stack, Grid, Divider } from "@mui/material";
import { SkeletonRoundedCard } from "../../styles";
const useStyles = makeStyles(() => ({
    container: {
        width: "214px",
        height: "206px",
        display: "flex",
        alignItems: "center",
        padding: "15px 18px 8px 18px",
        borderRadius: "10px",
        justifyContent: "center",
        background: colors.white,
    },
}));

const IconStyle = {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
};

const ConfigureNode = () => {
    return (
        <Stack spacing={2} alignItems="center">
            <Typography variant="body2">Configuring your node</Typography>
            <HistoryIcon />
        </Stack>
    );
};

type NodeCardProps = {
    title?: string;
    users?: string;
    loading?: boolean;
    subTitle?: string;
    isConfigure?: boolean;
    handleOptionItemClick?: Function;
};

const NodeCard = ({
    title,
    users,
    subTitle,
    loading,
    isConfigure = false,
    handleOptionItemClick = () => {},
}: NodeCardProps) => {
    const classes = useStyles();

    if (isConfigure)
        return (
            <Paper className={classes.container}>
                <ConfigureNode />
            </Paper>
        );

    return (
        <>
            {loading ? (
                <SkeletonRoundedCard
                    variant="rectangular"
                    width={204}
                    height={206}
                />
            ) : (
                <Paper className={classes.container}>
                    <Grid container spacing={0.8}>
                        <Grid item xs={10}>
                            <Grid textAlign="initial">
                                <Typography
                                    variant="subtitle1"
                                    sx={{
                                        fontWeight: 500,
                                        letterSpacing: "-0.02em",
                                        lineHeight: "19px",
                                    }}
                                >
                                    {title}
                                </Typography>
                            </Grid>
                            <Grid textAlign="initial">
                                <Typography variant="caption">
                                    {subTitle}
                                </Typography>
                            </Grid>
                        </Grid>
                        <Grid item xs={2} m="4px 0px">
                            <OptionsPopover
                                cid={"node-card-options"}
                                menuOptions={BASIC_MENU_ACTIONS}
                                handleItemClick={(type: string) =>
                                    handleOptionItemClick(type)
                                }
                            />
                        </Grid>
                        <Grid item xs={12} sx={{ ...IconStyle }}>
                            <NodeImg />
                        </Grid>
                        <Grid item xs={12}>
                            <Divider sx={{ m: "18px -18px 4px -18px" }} />
                        </Grid>
                        <Grid
                            item
                            xs={12}
                            container
                            spacing={1}
                            sx={{ alignItems: "center" }}
                        >
                            <Grid
                                item
                                xs={6}
                                container
                                display="flex"
                                alignSelf="end"
                                pt="0px !important"
                            >
                                <UsersIcon
                                    width="16px"
                                    height="16px"
                                    color={colors.black}
                                />
                                <Typography variant="caption" pl="8px">
                                    {users}
                                </Typography>
                            </Grid>
                            <Grid xs={2} item sx={{ ...IconStyle }}>
                                <ThermometerIcon />
                            </Grid>
                            <Grid xs={2} item sx={{ ...IconStyle }}>
                                <BatteryIcon />
                            </Grid>
                            <Grid xs={2} item sx={{ ...IconStyle }}>
                                <CpuIcon />
                            </Grid>
                        </Grid>
                    </Grid>
                </Paper>
            )}
        </>
    );
};

export default NodeCard;
