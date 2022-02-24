import { Stack, Typography } from "@mui/material";
import React from "react";
import { EmptyView } from "..";
import BarChartIcon from "@mui/icons-material/BarChart";
import { Variant } from "@mui/material/styles/createTypography";

interface IGraphTitleWrapper {
    title?: string;
    hasData?: boolean;
    variant?: Variant;
    children: React.ReactNode;
}

const GraphTitleWrapper = ({
    children,
    title = "",
    hasData = false,
    variant = "subtitle1",
}: IGraphTitleWrapper) => {
    return (
        <Stack spacing={2}>
            {title && (
                <Typography variant={variant} fontWeight={500}>
                    {title}
                </Typography>
            )}
            {hasData ? (
                children
            ) : (
                <EmptyView
                    size="large"
                    title="No activity yet!"
                    icon={BarChartIcon}
                />
            )}
        </Stack>
    );
};

export default GraphTitleWrapper;
