import {
    Box,
    Button,
    Grid,
    LinearProgress,
    Stack,
    Typography,
} from "@mui/material";
import { RoundedCard } from "../../styles";
import colors from "../../theme/colors";
const UserCard = () => {
    return (
        <>
            <Box>
                <Grid container spacing={2}>
                    <Grid item>
                        <RoundedCard sx={{ boxShadow: 2 }}>
                            <Stack spacing={1} direction="column">
                                <Typography variant="body2">
                                    19123192313128381239128381239
                                </Typography>
                                <Typography variant="h5">John Doe</Typography>
                            </Stack>
                            <Grid
                                container
                                spacing={3}
                                style={{ alignItems: "center" }}
                            >
                                <Grid item>
                                    <Stack direction="row" spacing={1}>
                                        <Typography variant="h5">
                                            1.5
                                        </Typography>
                                        <Typography
                                            variant="body2"
                                            sx={{
                                                position: "relative",
                                                top: "8px",
                                            }}
                                        >
                                            GB
                                        </Typography>
                                    </Stack>
                                </Grid>
                                <Grid item>
                                    <Typography
                                        variant="body2"
                                        sx={{
                                            position: "relative",
                                            top: "2px",
                                        }}
                                    >
                                        0.5 GB free data left
                                    </Typography>
                                </Grid>
                            </Grid>
                            <Grid container spacing={1}>
                                <Grid item xs={12}>
                                    <Box sx={{ width: "100%" }}>
                                        <LinearProgress
                                            variant="determinate"
                                            value={50}
                                            sx={{
                                                height: "8px",
                                                backgroundColor:
                                                    colors.darkGray,
                                            }}
                                        />
                                    </Box>
                                </Grid>
                                <Grid item xs={12}>
                                    <Button
                                        variant="text"
                                        sx={{ color: colors.darkGrey }}
                                    >
                                        VIEW MORE
                                    </Button>
                                </Grid>
                            </Grid>
                        </RoundedCard>
                    </Grid>
                </Grid>
            </Box>
        </>
    );
};
export default UserCard;
