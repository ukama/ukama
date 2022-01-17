import { useState } from "react";
import { Box, Grid } from "@mui/material";
import { NodeDetails, NodeStatus } from "../../components";
import { NodeDetailsStub, NodesData } from "../../constants/stubData";

const Nodes = () => {
    const [selectedNodeIndex, setSelectedNodeIndex] = useState(0);

    return (
        <Box sx={{ height: "calc(100vh - 8vh)", p: "28px 0px" }}>
            <Box sx={{ flexGrow: 1, pb: "18px" }}>
                <Grid container spacing={3}>
                    <Grid xs={12} item>
                        <NodeStatus
                            nodes={NodesData}
                            selectedNodeIndex={selectedNodeIndex}
                            setSelectedNodeIndex={setSelectedNodeIndex}
                        />
                    </Grid>
                    <Grid xs={12} item container spacing={3}>
                        <NodeDetails detailsList={NodeDetailsStub} />
                    </Grid>
                </Grid>
            </Box>
        </Box>
    );
};

export default Nodes;
