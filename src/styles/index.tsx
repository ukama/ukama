import { colors } from "../theme";
import { makeStyles } from "@mui/styles";
import { Box, styled, Link } from "@mui/material";

const globalUseStyles = makeStyles(() => ({
    inputFieldStyle: {
        height: "24px",
        padding: "12px 14px",
    },
    GridContainer: {
        borderRadius: "10px",
        backgroundColor: "#FFFFFF",
        boxShadow:
            "0px 2px 1px -1px rgba(0,0,0,0.2), 0px 1px 1px 0px rgba(0,0,0,0.14), 0px 1px 3px 0px rgba(0,0,0,0.12);",
    },
}));

const CenterContainer = styled(Box)({
    width: "100%",
    height: "100%",
    display: "flex",
    padding: "18px",
    alignItems: "center",
    flexDirection: "column",
    justifyContent: "center",
});

const LinkStyle = styled(Link)({
    fontSize: "0.75rem",
    width: "fit-content",
    alignSelf: "flex-end",
    color: colors.primary,
    letterSpacing: "0.4px",
    textDecoration: "none",
    "&:hover": {
        textDecoration: "underline",
    },
});
const MessageContainer = styled(Box)({
    paddingBottom: "5%",
});
const ContainerJustifySpaceBtw = styled(Box)(props => ({
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    [props.theme.breakpoints.down("sm")]: {
        flexDirection: "column",
    },
}));

export {
    LinkStyle,
    globalUseStyles,
    CenterContainer,
    ContainerJustifySpaceBtw,
    MessageContainer,
};
