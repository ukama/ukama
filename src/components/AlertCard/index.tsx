import { colors } from "../../theme";
import { SkeletonRoundedCard } from "../../styles";
import { Typography, Card, Grid, List, ListItem, Stack } from "@mui/material";
import { AlertItemType } from "../../types";
import { CloudOffIcon } from "../../assets/svg";

import moment from "moment";
type AlertCardProps = {
    alertOptions: any;
    loading?: boolean;
    iSalertsData?: boolean;
};
const AlertCard = ({ alertOptions, loading }: AlertCardProps) => {
    return (
        <>
            {loading ? (
                <SkeletonRoundedCard variant="rectangular" height={64} />
            ) : (
                <List
                    sx={{
                        pr: "4px",
                        maxHeight: 305,
                        overflow: "auto",
                        position: "relative",
                    }}
                >
                    {alertOptions.map(
                        ({
                            id,
                            alertDate,
                            description,
                            title,
                        }: AlertItemType) => (
                            <ListItem
                                key={id}
                                style={{
                                    padding: 1,
                                    marginBottom: "4px",
                                }}
                            >
                                <Card
                                    sx={{
                                        width: "100%",
                                        marginBottom: "3px",
                                        padding: "0px 10px 10px 10px",
                                    }}
                                    elevation={1}
                                >
                                    <Grid container spacing={1}>
                                        <Grid item xs={8} container>
                                            <Stack direction="row">
                                                <CloudOffIcon />
                                                <Typography
                                                    variant="body1"
                                                    color="initial"
                                                >
                                                    {title}
                                                </Typography>
                                            </Stack>
                                        </Grid>

                                        <Grid
                                            item
                                            container
                                            justifyContent="flex-end"
                                            xs={4}
                                        >
                                            <Typography
                                                variant="caption"
                                                color={colors.empress}
                                            >
                                                {moment(alertDate).format(
                                                    "DD/MM/YY h A"
                                                )}
                                            </Typography>
                                        </Grid>
                                        <Grid item xs={12}>
                                            <Typography
                                                variant="body2"
                                                color={colors.empress}
                                                sx={{
                                                    position: "relative",
                                                    bottom: "5px",
                                                    left: "25px",
                                                }}
                                            >
                                                {description}
                                            </Typography>
                                        </Grid>
                                    </Grid>
                                    {/* <Grid
                                        spacing={2}
                                        container
                                        direction="row"
                                        justifyContent="center"
                                    >
                                        <Grid
                                            item
                                            display="flex"
                                            alignItems="center"
                                            sx={{
                                                position: "relative",
                                                left: "5px",
                                            }}
                                        >
                                           
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
                                                    <Typography
                                                        variant="body1"
                                                        color="initial"
                                                    >
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
                                                        {moment(
                                                            alertDate
                                                        ).format(
                                                            "DD/MM/YY h A"
                                                        )}
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
                                                            alignItems:
                                                                "center",
                                                        }}
                                                    >
                                                        {description}
                                                    </Typography>
                                                </Grid>
                                            </Grid>
                                        </Grid>
                                    </Grid> */}
                                </Card>
                            </ListItem>
                        )
                    )}
                </List>
            )}
        </>
    );
};
export default AlertCard;
