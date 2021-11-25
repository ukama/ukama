import { colors } from "../../theme";
import { Typography, Card, Grid } from "@mui/material";
type AlertCardProps = {
    Icon: any;
    id: number;
    date: string;
    title: string;
    description: string;
};

const AlertCard = ({ date, description, title, Icon }: AlertCardProps) => {
    return (
        <Card
            sx={{
                width: "100%",
                marginBottom: "3px",
                padding: "0px 10px 10px 10px",
            }}
            elevation={1}
        >
            <Grid spacing={2} container direction="row" justifyContent="center">
                <Grid
                    item
                    display="flex"
                    alignItems="center"
                    sx={{
                        position: "relative",

                        left: "5px",
                    }}
                >
                    <Icon />
                </Grid>
                <Grid
                    xs={12}
                    item
                    sm
                    container
                    direction="column"
                    sx={{
                        position: "relative",
                        top: "8px",
                    }}
                >
                    <Grid
                        sm
                        item
                        container
                        spacing={2}
                        display="flex"
                        direction="row"
                        alignItems="center"
                    >
                        <Grid item xs={8} md sm lg>
                            <Typography variant="body1" color="initial">
                                {title}
                            </Typography>
                        </Grid>
                        <Grid
                            item
                            xs={4}
                            display="flex"
                            justifyContent="flex-end"
                        >
                            <Typography
                                variant="caption"
                                color={colors.empress}
                            >
                                {date}
                            </Typography>
                        </Grid>
                    </Grid>
                    <Grid item sm container>
                        <Grid item xs={12} container>
                            <Typography
                                variant="body2"
                                color={colors.empress}
                                sx={{
                                    display: "flex",
                                    alignItems: "center",
                                }}
                            >
                                {description}
                            </Typography>
                        </Grid>
                    </Grid>
                </Grid>
            </Grid>
        </Card>
    );
};
export default AlertCard;
