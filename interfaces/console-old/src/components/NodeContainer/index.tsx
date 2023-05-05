import { Box } from "@mui/material";
import { NodeDto } from "../../generated";
import { EmptyView, NodeSlider } from "..";
import RouterIcon from "@mui/icons-material/Router";
import React from "react";
type NodeContainerProps = {
    items: NodeDto[];
    handleItemAction: Function;
};

const NodeContainer = ({ items, handleItemAction }: NodeContainerProps) => {
    return (
        <Box
            component="div"
            sx={{
                display: "flex",
                minHeight: "246px",
                alignItems: "center",
            }}
        >
            {items.length > 0 ? (
                <NodeSlider items={items} handleItemAction={handleItemAction} />
            ) : (
                <EmptyView
                    size="large"
                    title="No nodes yet!"
                    icon={RouterIcon}
                />
            )}
        </Box>
    );
};

export default React.memo(NodeContainer);
