import React from "react";
import { EmptyView, TimeFilter } from "..";
import { Grid, Typography } from "@mui/material";
import BarChartIcon from "@mui/icons-material/BarChart";
import { Variant } from "@mui/material/styles/createTypography";

interface IGraphTitleWrapper {
    title?: string;
    filter?: string;
    hasData?: boolean;
    variant?: Variant;
    showFilter?: boolean;
    children: React.ReactNode;
    handleFilterChange?: Function;
}

const GraphTitleWrapper = ({
    children,
    title = "",
    filter = "DAY",
    hasData = false,
    showFilter = true,
    variant = "subtitle1",
    handleFilterChange,
}: IGraphTitleWrapper) => {
    return (
        <Grid item container spacing={2} my={2} width="100%">
            {(title || showFilter) && (
                <Grid item container width="100%">
                    {title && (
                        <Grid item xs={6}>
                            <Typography variant={variant} fontWeight={500}>
                                {title}
                            </Typography>
                        </Grid>
                    )}
                    {showFilter && (
                        <Grid
                            item
                            xs={6}
                            display="flex"
                            justifyContent="flex-end"
                        >
                            <TimeFilter
                                filter={filter}
                                handleFilterSelect={(v: string) =>
                                    handleFilterChange && handleFilterChange(v)
                                }
                            />
                        </Grid>
                    )}
                </Grid>
            )}
            <Grid item width="100%">
                {hasData ? (
                    children
                ) : (
                    <EmptyView
                        size="large"
                        title="No activity yet!"
                        icon={BarChartIcon}
                    />
                )}
            </Grid>
        </Grid>
    );
};

export default GraphTitleWrapper;
