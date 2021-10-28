import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import { UpgradeIcon } from "../../assets/svg";
import { Box, Button, Typography } from "@mui/material";

const useStyles = makeStyles(() => ({
    root: {
        display: "flex",
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
        margin: "18px",
        background: colors.aliceBlue,
    },
    buttonStyle: {
        width: "124px",
        fontWeight: 600,
        marginTop: "18px",
        letterSpacing: "0.4px",
        boxShadow:
            "0px 3px 1px -2px rgba(0, 0, 0, 0.2), 0px 2px 2px rgba(0, 0, 0, 0.14), 0px 1px 5px rgba(0, 0, 0, 0.12)",
    },
    imageWrapper: {
        zIndex: 10,
        alignSelf: "center",
        marginBottom: "-40px",
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
                <Typography variant="body1" sx={{ textAlign: "center" }}>
                    Enjoy more features with Ukama Pro!
                </Typography>
                <Button
                    size="medium"
                    type="submit"
                    color="primary"
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
