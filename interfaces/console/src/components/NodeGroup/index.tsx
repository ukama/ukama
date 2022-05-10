import { Grid, Link, Typography } from "@mui/material";
import { OrgNodeDto } from "../../generated";
import LoadingWrapper from "../LoadingWrapper";
interface INodeGroup {
    nodes: OrgNodeDto[];
    loading: boolean;
    handleNodeAction: Function;
}

const NodeGroup = ({ nodes, loading, handleNodeAction }: INodeGroup) => {
    return (
        <Grid container spacing={2} alignItems="center">
            <Grid item xs={5}>
                <Typography fontWeight={500} variant="body2">
                    Node Group
                </Typography>
            </Grid>
            <Grid item xs={7}>
                <LoadingWrapper isLoading={loading} height={24} radius="small">
                    {nodes.length > 0 ? (
                        nodes.map(item => (
                            <Link
                                variant="body2"
                                fontWeight={500}
                                key={item.nodeId}
                                underline="always"
                                sx={{ textTransform: "capitalize" }}
                                onClick={() => handleNodeAction(item.nodeId)}
                            >
                                {item.name}
                            </Link>
                        ))
                    ) : (
                        <Typography fontWeight={500} variant="body2">
                            N/A
                        </Typography>
                    )}
                </LoadingWrapper>
            </Grid>
        </Grid>
    );
};

export default NodeGroup;
