import React from "react";
import { Box } from "@mui/material";

interface TabPanelProps {
    id: string;
    index: number;
    value: number;
    children?: React.ReactNode;
}

const TabPanel = ({ id, index, value, children }: TabPanelProps) => {
    return (
        <div
            id={id}
            role="tabpanel"
            aria-labelledby={id}
            hidden={value !== index}
        >
            {value === index && <Box component="div">{children}</Box>}
        </div>
    );
};

export default TabPanel;
