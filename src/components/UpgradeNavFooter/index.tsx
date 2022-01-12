import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import { UpgradeIcon } from "../../assets/svg";
import { Box, Button, Typography } from "@mui/material";

const useStyles = makeStyles(() => ({
    root: {
        display: "flex",
        position: "relative",
        flexDirection: "column",
    },
    container: {
        display: "flex",
        height: "148px",
        padding: "0px 18px",
        borderRadius: "10px",
        alignItems: "center",
        flexDirection: "column",
        justifyContent: "center",
        background: colors.aliceBlue,
    },
    buttonStyle: {
        width: "124px",
        marginTop: "18px",
    },
    imageWrapper: {
        top: -65,
        zIndex: 10,
        alignSelf: "center",
        position: "absolute",
    },
}));

const UpgradeNavFooter = () => {
    const classes = useStyles();
    return (
        <Box className={classes.root}>
            <Box className={classes.imageWrapper}>
                <UpgradeIcon />
            </Box>
            <Box className={classes.container}>
                <Typography variant="body2" sx={{ textAlign: "center" }}>
                    Enjoy more features with Ukama Pro!
                </Typography>
                <Button
                    type="submit"
                    size="medium"
                    variant="contained"
                    className={classes.buttonStyle}
                >
                    UPGRADE
                </Button>
            </Box>
        </Box>
    );
};

export default UpgradeNavFooter;
