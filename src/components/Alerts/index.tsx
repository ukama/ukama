import { format } from "date-fns";
import { AlertDto } from "../../generated";
import { getColorByType } from "../../utils";
import { CloudOffIcon } from "../../assets/svg";
import { Typography, Grid, List, ListItem, Stack } from "@mui/material";

type AlertsProps = {
    alertOptions: AlertDto[] | undefined;
};

const Alerts = ({ alertOptions = [] }: AlertsProps) => {
    return (
        <List
            sx={{
                width: 364,
                maxHeight: 305,
                overflowY: "auto",
                overflowX: "hidden",
                position: "relative",
            }}
        >
            {alertOptions.map(
                ({ id, alertDate, description, title, type }: AlertDto) => (
                    <ListItem key={id} sx={{ p: "0px 12px 8px 0px" }}>
                        <Grid container spacing={1}>
                            <Grid item xs={8} container>
                                <Stack direction="row" spacing={1}>
                                    <CloudOffIcon
                                        color={getColorByType(type)}
                                    />
                                    <Typography variant="body1">
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
                                    color={"textSecondary"}
                                >
                                    {format(alertDate, "M/d/YY h a")}
                                </Typography>
                            </Grid>
                            <Grid item xs={12}>
                                <Typography
                                    variant="body2"
                                    color={"textSecondary"}
                                    sx={{
                                        left: "36px",
                                        bottom: "5px",
                                        overflow: "hidden",
                                        position: "relative",
                                    }}
                                >
                                    {description}
                                </Typography>
                            </Grid>
                        </Grid>
                    </ListItem>
                )
            )}
        </List>
    );
};
export default Alerts;
