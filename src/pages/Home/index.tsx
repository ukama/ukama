import { styled } from "@mui/material/styles";
import { Grid, Paper, Box } from "@mui/material";
import {
    GraphContainer,
    AlertContainer,
    NodeContainer,
    ResidentContainer,
} from "../../components";
const Item = styled(Paper)(({ theme }) => ({
    ...theme.typography.body2,
    padding: theme.spacing(1),
    textAlign: "center",
    color: theme.palette.text.secondary,
}));

const Home = () => {
    return (
        <>
            <Box sx={{ flexGrow: 1 }}>
                <Grid container spacing={2}>
                    <GraphContainer />

                    <AlertContainer />
                    <NodeContainer />
                    <ResidentContainer />
                </Grid>
            </Box>
        </>
    );
};

export default Home;
