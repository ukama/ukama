import styled from "@emotion/styled";
import { colors } from "../../theme";
import { Container, Box } from "@mui/material";

const RootContainer = styled(Container)({
    height: "auto",
    padding: "0px !important",
    background: colors.white,
    boxShadow:
        "-4px 0px 4px 4px rgba(0, 0, 0, 0.05), 4px 4px 4px 4px rgba(0, 0, 0, 0.05)",
    borderRadius: "5px",
});

const GradiantBar = styled(Box)({
    width: "100%",
    height: "12px",
    background:
        "linear-gradient(90deg, #00D3EB 0%, #2190F6 14.06%, #6974F8 44.27%, #6974F8 58.85%, #271452 100%)",
    borderRadius: "4px 4px 0px 0px",
});

const ComponentContainer = {
    width: "auto",
    height: "auto",
    overflow: "hidden",
    margin: {
        xs: "24px 28px",
        sm: "44px 64px",
        md: "54px 80px",
    },
};

export { RootContainer, GradiantBar, ComponentContainer };
