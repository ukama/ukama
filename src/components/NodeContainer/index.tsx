import { Box } from "@mui/material";
import { NodeDto } from "../../generated";
import { EmptyView, NodeSlider } from "..";
import RouterIcon from "@mui/icons-material/Router";
type NodeContainerProps = {
    items: NodeDto[];
    handleItemAction: Function;
    handleNodeUpdate: Function;
};

const NodeContainer = ({
    items,
    handleItemAction,
    handleNodeUpdate,
}: NodeContainerProps) => {
    return (
        <Box
            component="div"
            sx={{ display: "flex", minHeight: "246px", alignItems: "center" }}
        >
            {items.length > 0 ? (
                <NodeSlider
                    items={items}
                    handleItemAction={handleItemAction}
                    handleNodeUpdate={handleNodeUpdate}
                />
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

export default NodeContainer;
