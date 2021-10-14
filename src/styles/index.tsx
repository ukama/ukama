import { colors } from "../theme";
import { makeStyles } from "@mui/styles";
import { Box, styled, Link } from "@mui/material";

const globalUseStyles = makeStyles(() => ({
    inputFieldStyle: {
        height: "24px",
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
};
