import React from "react";
import { EmptyView, GraphLoading, TimeFilter } from "..";
import { Grid, Typography } from "@mui/material";
import BarChartIcon from "@mui/icons-material/BarChart";
import { Variant } from "@mui/material/styles/createTypography";

interface IGraphTitleWrapper {
    title?: string;
    filter?: string;
    hasData?: boolean;
    loading?: boolean;
    variant?: Variant;
    showFilter?: boolean;
    children: React.ReactNode;
    handleFilterChange?: Function;
}

const GraphTitleWrapper = ({
    children,
    title = "",
    filter = "LIVE",
    hasData = false,
    loading = true,
    showFilter = true,
    variant = "subtitle1",
    handleFilterChange,
}: IGraphTitleWrapper) => {
    return (
        <Grid item container width="100%">
            {(title || showFilter) && (
                <Grid item container width="100%" mb={2}>
                    {title && (
                        <Grid item xs={6}>
                            <Typography variant={variant} fontWeight={500}>
                                {title}
                            </Typography>
                        </Grid>
                    )}
                    {hasData && showFilter && (
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
            <Grid
                item
                container
                height={"300px"}
                alignItems={"center"}
                justifyContent="center"
            >
                {loading ? (
                    <GraphLoading />
                ) : hasData ? (
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
