import { useState } from "react";
import { Box, Grid } from "@mui/material";
import { NodeDetails, NodeStatus, LoadingWrapper } from "../../components";
import { NodeDetailsStub } from "../../constants/stubData";
import { useGetNodesByOrgQuery } from "../../generated";
import { useRecoilValue } from "recoil";
import { isSkeltonLoading, organizationId } from "../../recoil";

const Nodes = () => {
    const [selectedNodeIndex, setSelectedNodeIndex] = useState(0);
    const orgId = useRecoilValue(organizationId);
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const { data: nodesRes, loading: nodesLoading } = useGetNodesByOrgQuery({
        variables: { orgId: orgId || "" },
    });

    return (
        <Box sx={{ height: "calc(100vh - 8vh)" }}>
            <Box sx={{ flexGrow: 1 }}>
                <Grid container spacing={3}>
                    <Grid xs={12} item>
                        <LoadingWrapper
                            width="100%"
                            height="30px"
                            isLoading={isSkeltonLoad || nodesLoading}
                        >
                            <NodeStatus
                                nodes={nodesRes?.getNodesByOrg?.nodes}
                                selectedNodeIndex={selectedNodeIndex}
                                setSelectedNodeIndex={setSelectedNodeIndex}
                            />
                        </LoadingWrapper>
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
