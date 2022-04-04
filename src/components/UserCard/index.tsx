import { colors } from "../../theme";
import { GetUsersDto } from "../../generated";
import { Grid, Button, Typography, LinearProgress, Stack } from "@mui/material";

type UserCardProps = {
    user: GetUsersDto;
    // eslint-disable-next-line no-unused-vars
    handleMoreUserdetails: (user: GetUsersDto) => void;
};

const UserCard = ({ user, handleMoreUserdetails }: UserCardProps) => {
    return (
        <Grid container spacing={2}>
            <Grid item xs={12}>
                <Typography variant="body2" color="textSecondary">
                    {user.id}
                </Typography>
            </Grid>
            <Grid item xs={12}>
                <Typography variant="h5">{user.name}</Typography>
            </Grid>
            <Grid item xs={4}>
                <Stack direction="row" spacing={"4px"} alignItems="baseline">
                    <Typography variant="h5">{user.dataPlan}</Typography>
                    <Typography variant="body2" textAlign={"end"}>
                        MB
                    </Typography>
                </Stack>
            </Grid>
            <Grid item xs={8} alignSelf="end" mb={"2px"}>
                <Typography variant="body2" textAlign={"end"}>
                    {`${user.dataUsage} MB free data left`}
                </Typography>
            </Grid>
            <Grid item xs={12}>
                <LinearProgress
                    variant="determinate"
                    value={(user.dataUsage * 100) / user.dataPlan}
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
