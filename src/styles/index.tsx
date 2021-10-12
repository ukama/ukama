import { Box, styled, Link } from "@mui/material";
import { colors } from "../theme";

const CenterContainer = styled(Box)({
    width: "100%",
    height: "100%",
    display: "flex",
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

export { CenterContainer, LinkStyle };
