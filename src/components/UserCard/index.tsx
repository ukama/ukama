import {
    Button,
    Grid,
    LinearProgress,
    Stack,
    Paper,
    Typography,
} from "@mui/material";
import colors from "../../theme/colors";
type UserCardProps = {
    userDetails?: any;
    children?: any;
    handleMoreUserdetails?: any;
};
const UserCard = ({ userDetails, handleMoreUserdetails }: UserCardProps) => {
    return (
        <>
            <Grid
                container
                spacing={3}
                justifyContent="center"
                alignItems="center"
            >
                {userDetails.map(
                    ({ id, name, eSimNumber, dataPlan, dataUsage }: any) => (
                        <Grid item xs={12} md={6} lg={3} key={id}>
                            <Paper elevation={2} sx={{ p: 2 }}>
                                <Grid
                                    container
                                    item
                                    spacing={1}
                                    direction="column"
                                >
                                    <Grid item>
                                        <Typography
                                            variant="body2"
                                            sx={{ color: colors.empress }}
                                        >
                                            {eSimNumber}
                                        </Typography>
                                    </Grid>
                                    <Grid item>
                                        <Typography variant="h5">
                                            {name}
                                        </Typography>
                                    </Grid>
                                </Grid>
                                <Grid item container spacing={2} sx={{ mt: 1 }}>
                                    <Grid
                                        item
                                        xs={12}
                                        md={6}
                                        lg={3}
                                        container
                                        justifyContent="flex-start"
                                    >
                                        <Stack direction="row" spacing={1 / 2}>
                                            <Typography
                                                variant="h5"
                                                sx={{
                                                    position: "relative",
                                                    bottom: "9px",
                                                }}
                                            >
                                                {dataPlan}
                                            </Typography>
                                            <Typography variant="body2">
                                                GB
                                            </Typography>
                                        </Stack>
                                    </Grid>
                                    <Grid
                                        item
                                        container
                                        justifyContent="flex-end"
                                        xs={12}
                                        md={6}
                                        lg={9}
                                    >
                                        <Typography variant="body2">
                                            {dataUsage} GB free data left
                                        </Typography>
                                    </Grid>
                                </Grid>
                                <Grid container spacing={2} direction="column">
                                    <Grid item>
                                        <LinearProgress
                                            variant="determinate"
                                            value={dataPlan - dataUsage}
                                            sx={{
                                                height: "8px",
                                                backgroundColor:
                                                    colors.darkGray,
                                            }}
                                        />
                                    </Grid>
                                    <Grid item>
                                        <Button
                                            variant="text"
                                            sx={{ color: colors.darkGrey }}
                                            onClick={() =>
                                                handleMoreUserdetails(id)
                                            }
                                        >
                                            VIEW MORE
                                        </Button>
                                    </Grid>
                                </Grid>
                            </Paper>
                        </Grid>
                    )
                )}
            </Grid>
        </>
    );
};
export default UserCard;
