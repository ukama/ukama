import React from "react";
import { LoadingWrapper } from "..";
import { colors } from "../../theme";
import MenuOpenIcon from "@mui/icons-material/MenuOpen";
import { IconButton, Paper, Stack, Typography } from "@mui/material";
interface INodeStatsContainer {
    index: number;
    title: string;
    loading: boolean;
    selected?: number;
    isAlert?: boolean; //Pass true to show red border
    isCollapse?: boolean;
    isClickable?: boolean;
    onCollapse?: Function;
    isCollapsable?: boolean;
    handleAction?: Function;
    children: React.ReactNode;
}

const NodeStatsContainer = ({
    index,
    title,
    loading,
    children,
    onCollapse,
    handleAction,
    selected = -1,
    isAlert = false,
    isCollapse = false,
    isClickable = false,
    isCollapsable = false,
}: INodeStatsContainer) => {
    return (
        <LoadingWrapper
            width="100%"
            height="100px"
            radius="small"
            isLoading={loading}
        >
            <Paper
                sx={{
                    padding: "24px 24px 24px 0px",
                    cursor:
                        isCollapsable || !isClickable ? "default" : "pointer",
                    paddingLeft:
                        isAlert && selected !== index ? "16px" : "24px",
                    borderLeft: {
                        md:
                            selected === index
                                ? `8px solid ${colors.secondaryMain}`
                                : isAlert
                                ? `1px solid ${colors.error}`
                                : `8px solid ${colors.silver}`,
                    },
                    border: isAlert ? `0.5px solid ${colors.error}` : "none",
                }}
                onClick={() =>
                    isClickable && handleAction && handleAction(index)
                }
            >
                <Stack
                    direction="row"
                    justifyContent="space-between"
                    alignItems="center"
                    spacing={2}
                    sx={{ mb: 1 }}
                >
                    {!isCollapse && (
                        <Typography variant="h6">{title}</Typography>
                    )}
                    {isCollapsable && (
                        <IconButton
                            sx={{
                                position: isCollapse ? "relative" : null,
                                right: isCollapse ? 10 : null,
                                transform: isCollapse
                                    ? "rotate(180deg)"
                                    : "none",
                            }}
                            onClick={() => onCollapse && onCollapse()}
                        >
                            <MenuOpenIcon fontSize="medium" />
                        </IconButton>
                    )}
                </Stack>
                {!isCollapse && children}
            </Paper>
        </LoadingWrapper>
    );
};

export default NodeStatsContainer;
