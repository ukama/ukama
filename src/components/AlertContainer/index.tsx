import { Typography, Grid } from "@mui/material";
import { globalUseStyles } from "../../styles";
import { AlertIcon } from "../../assets/svg";
const AlertContainer = () => {
    const classes = globalUseStyles();

    return (
        <>
            <Grid item xs={4} className={classes.GridContainer}>
                <Grid item xs={12} container justifyContent="flex-start">
                    <Typography variant="h6">Alerts</Typography>
                </Grid>
                <Grid item xs={12} container justifyContent="center">
                    <AlertIcon />
                </Grid>
            </Grid>
        </>
    );
};

export default AlertContainer;
