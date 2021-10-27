import { styled } from "@mui/material/styles";
import { Grid, Paper, Box } from "@mui/material";
import { GraphContainer, AlertContainer } from "../../components";
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

                    <Grid item xs={4}>
                        <Item>xs=4</Item>
                    </Grid>
                    <Grid item xs={8}>
                        <Item>xs=8</Item>
                    </Grid>
                </Grid>
            </Box>
        </>
    );
};

export default Home;
