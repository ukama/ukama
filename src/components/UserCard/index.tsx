import { colors } from "../../theme";
import { GetUserDto } from "../../generated";
import { Grid, Button, Typography, LinearProgress, Stack } from "@mui/material";

type UserCardProps = {
    user: GetUserDto;
    // eslint-disable-next-line no-unused-vars
    handleMoreUserdetails: (user: GetUserDto) => void;
};

const UserCard = ({ user, handleMoreUserdetails }: UserCardProps) => {
    return (
        <Grid container spacing={2}>
            <Grid item xs={12}>
                <Typography variant="body2" color="textSecondary">
                    {user.eSimNumber}
                </Typography>
            </Grid>
            <Grid item xs={12}>
                <Typography variant="h5">{user.name}</Typography>
            </Grid>
            <Grid item xs={6}>
                <Stack direction="row" spacing={"4px"} alignItems="baseline">
                    <Typography variant="h5">{user.dataPlan}</Typography>
                    <Typography variant="body2" textAlign={"end"}>
                        GB
                    </Typography>
                </Stack>
            </Grid>
            <Grid item xs={6} alignSelf="end" mb={"2px"}>
                <Typography variant="body2" textAlign={"end"}>
                    {`${user.dataUsage} GB free data left`}
                </Typography>
            </Grid>
            <Grid item xs={12}>
                <LinearProgress
                    variant="determinate"
                    value={user.dataPlan - user.dataUsage}
                    sx={{
                        height: "8px",
                        borderRadius: "2px",
                        backgroundColor: colors.silver,
                    }}
                />
            </Grid>
            <Grid item xs={12}>
                <Button
                    variant="text"
                    onClick={() => handleMoreUserdetails(user)}
                >
                    VIEW MORE
                </Button>
            </Grid>
        </Grid>
    );
};
export default UserCard;
