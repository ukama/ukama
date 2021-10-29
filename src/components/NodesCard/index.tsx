import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import {
    CpuIcon,
    NodeImg,
    UsersIcon,
    HistoryIcon,
    BatteryIcon,
    ThermometerIcon,
} from "../../assets/svg";
import OptionsPopover from "../OptionsPopover";
import { BASIC_MENU_ACTIONS } from "../../constants";
import { Paper, Typography, Stack, Grid, Divider } from "@mui/material";
const useStyles = makeStyles(() => ({
    container: {
        width: "214px",
        height: "206px",
        display: "flex",
        alignItems: "center",
        padding: "15px 18px",
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
    subTitle?: string;
    isConfigure?: boolean;
};

const NodeCard = ({
    title,
    subTitle,
    users,
    isConfigure = false,
}: NodeCardProps) => {
    const classes = useStyles();

    if (isConfigure)
        return (
            <Paper className={classes.container}>
                <ConfigureNode />
            </Paper>
        );

    return (
        <Paper className={classes.container}>
            <Grid container spacing={1.5}>
                <Grid item xs={10}>
                    <Grid>
                        <Typography variant="subtitle1">{title}</Typography>
                    </Grid>
                    <Grid>
                        <Typography variant="subtitle2">{subTitle}</Typography>
                    </Grid>
                </Grid>
                <Grid item xs={2} m="4px 0px">
                    <OptionsPopover
                        cid={"sadasda"}
                        options={BASIC_MENU_ACTIONS}
                        handleItemClick={() => {}}
                    />
                </Grid>
                <Grid item xs={12} sx={{ ...IconStyle }}>
                    <NodeImg />
                </Grid>
                <Grid item xs={12}>
                    <Divider sx={{ margin: "0px -18px" }} />
                </Grid>
                <Grid item xs={12} container spacing={1}>
                    <Grid xs={6} item container>
                        <UsersIcon
                            width="16px"
                            height="16px"
                            color={colors.black}
                        />
                        <Typography variant="body2" pl="8px">
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
    );
};

export default NodeCard;
