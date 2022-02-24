import { LoadingWrapper } from "..";
import { colors } from "../../theme";
import React, { useState } from "react";
import MenuOpenIcon from "@mui/icons-material/MenuOpen";
import { ContainerJustifySpaceBtw } from "../../styles";
import { IconButton, Paper, Typography } from "@mui/material";
interface INodeStatsContainer {
    index: number;
    title: string;
    loading: boolean;
    selected?: number;
    isAlert?: boolean; //Pass true to show red border
    isClickable?: boolean;
    isCollapsable?: boolean;
    handleAction?: Function;
    children: React.ReactNode;
}

const NodeStatsContainer = ({
    index,
    title,
    loading,
    children,
    handleAction,
    selected = -1,
    isAlert = false,
    isClickable = false,
    isCollapsable = false,
}: INodeStatsContainer) => {
    const [isCollapse, setIsCollapse] = useState(false);

    return (
        <LoadingWrapper
            width="100%"
            height="100px"
            radius="small"
            isLoading={loading}
        >
            <Paper
                sx={{
                    minWidth: isCollapse ? "fit-content" : 340,
                    padding: "24px 24px 24px 0px",
                    cursor:
                        isCollapsable || !isClickable ? "defautl" : "pointer",
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
                <ContainerJustifySpaceBtw sx={{ pb: 1 }}>
                    {!isCollapse && (
                        <Typography variant="h6">{title}</Typography>
                    )}
                    {isCollapsable && (
                        <IconButton
                            sx={{
                                p: 0,
                                transform: isCollapse
                                    ? "rotate(180deg)"
                                    : "none",
                            }}
                            onClick={() => setIsCollapse(!isCollapse)}
                        >
                            <MenuOpenIcon />
                        </IconButton>
                    )}
                </ContainerJustifySpaceBtw>
                {!isCollapse && children}
            </Paper>
        </LoadingWrapper>
    );
};

export default NodeStatsContainer;
