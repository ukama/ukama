import { RoundedCard } from "../../styles";
import { AmountBalanceImg } from "../../assets/svg";
import { Typography, Grid, Button, Chip } from "@mui/material";

type CurrentBillProps = {
    title: string;
    amount: string;
    dueDate?: string;
    periodOf: string;
    handleMakePayment: Function;
};

const CurrentBill = ({
    title,
    amount,
    dueDate,
    periodOf,
    handleMakePayment,
}: CurrentBillProps) => {
    return (
        <RoundedCard>
            <Grid container>
                <Grid item container>
                    <Grid xs={12} sm={7} item>
                        <Typography variant="h6">{title}</Typography>
                        <Typography variant="body2">{periodOf}</Typography>
                        <Typography variant="h3" sx={{ m: "18px 0px" }}>
                            {amount}
                        </Typography>
                        {dueDate && (
                            <Chip
                                label={dueDate}
                                sx={{
                                    backgroundColor: "rgba(227, 0, 0, 0.2)",
                                    typography: "caption",
                                }}
                            />
                        )}
                    </Grid>
                    <Grid
                        xs={12}
                        sm={5}
                        item
                        container
                        display="flex"
                        spacing={"8px"}
                        justifyContent="flex-end"
                    >
                        {dueDate && (
                            <Grid item>
                                <Button
                                    variant="contained"
                                    sx={{ width: "191px" }}
                                    onClick={() => handleMakePayment()}
                                >
                                    MAKE PAYMENT
                                </Button>
                            </Grid>
                        )}
                        <Grid item>
                            <AmountBalanceImg />
                        </Grid>
                    </Grid>
                </Grid>
                {!dueDate && (
                    <Grid xs={12} item>
                        <Typography variant="caption">
                            *Automatically charged to card on November 10, 2021
                        </Typography>
                    </Grid>
                )}
            </Grid>
        </RoundedCard>
    );
};
export default CurrentBill;
