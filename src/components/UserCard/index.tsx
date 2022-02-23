import { colors } from "../../theme";
import { Grid, Button, Typography, LinearProgress, Stack } from "@mui/material";

type UserCardProps = {
    id: string;
    name: string;
    dataPlan: number;
    dataUsage: number;
    eSimNumber: string;
    handleMoreUserdetails?: any;
};

const UserCard = ({
    id,
    name,
    dataPlan,
    dataUsage,
    eSimNumber,
    handleMoreUserdetails,
}: UserCardProps) => {
    return (
        <Grid container spacing={2}>
            <Grid item xs={12}>
                <Typography variant="body2" color="textSecondary">
                    {eSimNumber}
                </Typography>
            </Grid>
            <Grid item xs={12}>
                <Typography variant="h5">{name}</Typography>
            </Grid>
            <Grid item xs={6}>
                <Stack direction="row" spacing={"4px"} alignItems="baseline">
                    <Typography variant="h5">{`${dataPlan}`}</Typography>
                    <Typography variant="body2" textAlign={"end"}>
                        GB
                    </Typography>
                </Stack>
            </Grid>
            <Grid item xs={6} alignSelf="end" mb={"2px"}>
                <Typography variant="body2" textAlign={"end"}>
                    {`${dataUsage} GB free data left`}
                </Typography>
            </Grid>
            <Grid item xs={12}>
                <LinearProgress
                    variant="determinate"
                    value={dataPlan - dataUsage}
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
                    onClick={() => handleMoreUserdetails(id)}
                >
                    VIEW MORE
                </Button>
            </Grid>
        </Grid>
    );
};
export default UserCard;
