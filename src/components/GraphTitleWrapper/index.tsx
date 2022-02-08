import { Stack, Typography } from "@mui/material";
import { Variant } from "@mui/material/styles/createTypography";
import React from "react";

interface IGraphTitleWrapper {
    title: string;
    variant?: Variant;
    children: React.ReactNode;
}

const GraphTitleWrapper = ({
    children,
    title = "",
    variant = "subtitle1",
}: IGraphTitleWrapper) => {
    return (
        <Stack spacing={2}>
            <Typography variant={variant} pl={2}>
                {title}
            </Typography>
            {children}
        </Stack>
    );
};

export default GraphTitleWrapper;
