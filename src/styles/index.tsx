import { colors } from "../theme";
import { makeStyles } from "@mui/styles";
import { Box, styled, Link, Paper, Skeleton } from "@mui/material";

const globalUseStyles = makeStyles(() => ({
    inputFieldStyle: {
        height: "24px",
        padding: "12px 14px",
    },
    disableInputFieldStyle: {
        padding: "4px 0px",
        "-webkit-text-fill-color": `${colors.black} !important`,
    },
    GridContainer: {
        padding: "1em",
    },
}));

const HorizontalContainerJustify = styled(Box)({
    width: "100%",
    height: "auto",
    display: "flex",
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
});

const HorizontalContainer = styled(Box)({
    width: "100%",
    height: "auto",
    display: "flex",
    alignItems: "center",
    flexDirection: "row",
});

const VerticalContainer = styled(Box)({
    width: "100%",
    height: "auto",
    display: "flex",
    alignItems: "center",
    flexDirection: "column",
});

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

const RoundedCard = styled(Paper)(props => ({
    width: "100%",
    padding: "18px 28px",
    height: "100%",
    borderRadius: "10px",
    display: "inline-block",
    background:
        props.theme.palette.mode === "dark" ? colors.darkGreen05 : colors.white,
    boxShadow: "2px 2px 6px rgba(0, 0, 0, 0.05)",
    [props.theme.breakpoints.down("sm")]: {
        padding: "18px",
    },
}));

const SkeletonRoundedCard = styled(Skeleton)(() => ({
    width: "100%",
    height: "100%",
    borderRadius: "10px",
    display: "inline-block",
}));

const FullscreenContainer = styled(Box)(props => ({
    width: "100%",
    height: "100%",
    background:
        props.theme.palette.mode === "dark"
            ? colors.darkGreen
            : colors.solitude,
}));

export {
    LinkStyle,
    RoundedCard,
    globalUseStyles,
    CenterContainer,
    MessageContainer,
    VerticalContainer,
    SkeletonRoundedCard,
    HorizontalContainer,
    FullscreenContainer,
    ContainerJustifySpaceBtw,
    HorizontalContainerJustify,
};
