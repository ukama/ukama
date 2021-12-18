import {
    Button,
    Grid,
    LinearProgress,
    Stack,
    Paper,
    Typography,
} from "@mui/material";
import colors from "../../theme/colors";
import { UserDetailsDialog } from "../../components";
import { RoundedCard } from "../../styles";
import { useState } from "react";
type UserCardProps = {
    userDetails?: any;
    children?: any;
};
const UserCard = ({ userDetails, children }: UserCardProps) => {
    const [showSimDialog, setShowSimDialog] = useState(false);
    const showMore = () => {
        setShowSimDialog(true);
    };
    const handleSimDialog = () => {
        setShowSimDialog(false);
    };

    return (
        <>
            <RoundedCard sx={{ height: "100%" }}>
                {children}
                <Grid
                    container
                    spacing={3}
                    justifyContent="center"
                    alignItems="center"
                >
                    {userDetails.map(
                        ({ id, name, imsi, dataPack, remaingData }: any) => (
                            <Grid item xs={12} md={6} lg={3} key={id}>
                                <Paper
                                    elevation={2}
                                    sx={{ boxShadow: 2, p: 2 }}
                                >
                                    <Grid
                                        container
                                        item
                                        spacing={1}
                                        direction="column"
                                    >
                                        <Grid item>
                                            <Typography variant="body2">
                                                {imsi}
                                            </Typography>
                                        </Grid>
                                        <Grid item>
                                            <Typography variant="h5">
                                                {name}
                                            </Typography>
                                        </Grid>
                                    </Grid>
                                    <Grid
                                        item
                                        container
                                        spacing={2}
                                        sx={{ mt: 1 }}
                                    >
                                        <Grid
                                            item
                                            xs={12}
                                            md={6}
                                            lg={3}
                                            container
                                            justifyContent="flex-start"
                                        >
                                            <Stack
                                                direction="row"
                                                spacing={1 / 2}
                                            >
                                                <Typography
                                                    variant="h5"
                                                    sx={{
                                                        position: "relative",
                                                        bottom: "9px",
                                                    }}
                                                >
                                                    {dataPack}
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
                                                {remaingData} GB free data left
                                            </Typography>
                                        </Grid>
                                    </Grid>
                                    <Grid
                                        container
                                        spacing={2}
                                        direction="column"
                                    >
                                        <Grid item>
                                            <LinearProgress
                                                variant="determinate"
                                                value={dataPack - remaingData}
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
                                                onClick={showMore}
                                            >
                                                VIEW MORE
                                            </Button>
                                        </Grid>
                                        <UserDetailsDialog
                                            userName="John Doe"
                                            data="- 1.5 GB data used, 0.5 free GB left"
                                            isOpen={showSimDialog}
                                            userDetailsTitle="User Details"
                                            btnLabel="Submit"
                                            handleClose={handleSimDialog}
                                            simDetailsTitle="SIM Details"
                                            saveBtnLabel="save"
                                            closeBtnLabel="close"
                                        />
                                    </Grid>
                                </Paper>
                            </Grid>
                        )
                    )}
                </Grid>
            </RoundedCard>
        </>
    );
};
export default UserCard;
