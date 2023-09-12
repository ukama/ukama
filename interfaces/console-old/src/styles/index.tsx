import { colors } from "../theme";
import { hexToRGB } from "../utils";
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
    backToNodeGroupButtonStyle: {
        position: "fixed",
        left: "50%",
        bottom: "20px",
        transform: "translate(-50%, -50%)",
        margin: "0 auto",
        pointer: "cursor",
    },
    GridContainer: {
        padding: "1em",
    },
}));

const HorizontalContainerJustify = styled(Box)(() => ({
    width: "100%",
    height: "auto",
    display: "flex",
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
}));

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
    fontSize: "14px",
    width: "fit-content",
    alignSelf: "flex-end",
    color: colors.primaryMain,
    letterSpacing: "0.4px",
    textDecoration: "none",
    "&:hover": {
        textDecoration: "underline",
    },
});

const MessageContainer = styled(Box)({
    paddingBottom: "5%",
});

const ContainerJustifySpaceBtw = styled(Box)({
    width: "100%",
    display: "flex",
    paddingBottom: 10,
    flexDirection: "row",
    justifyContent: "space-between",
    textAlign: "center",
});

const RoundedCard = styled(Paper)(props => ({
    width: "100%",
    padding: "18px 28px",
    height: "100%",
    borderRadius: "10px",
    display: "inline-block",
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

const FullscreenContainer = styled(Box)(() => ({
    width: "100%",
    height: "100%",
}));

const SimpleCardWithBorder = styled(Box)(props => ({
    borderRadius: "4px",
    border: `1px solid ${hexToRGB(props.theme.palette.text.primary, 0.1)}`,
}));
export {
    LinkStyle,
    RoundedCard,
    globalUseStyles,
    CenterContainer,
    MessageContainer,
    VerticalContainer,
    SimpleCardWithBorder,
    SkeletonRoundedCard,
    HorizontalContainer,
    FullscreenContainer,
    ContainerJustifySpaceBtw,
    HorizontalContainerJustify,
};
