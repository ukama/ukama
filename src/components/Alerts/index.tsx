import { format } from "date-fns";
import { AlertDto } from "../../generated";
import { getColorByType } from "../../utils";
import CloudOffIcon from "@mui/icons-material/CloudOff";
import { Typography, Grid, List, ListItem } from "@mui/material";

type AlertsProps = {
    alertOptions: AlertDto[] | undefined;
};

const PROPS = {
    display: "flex",
    alignItems: "center",
};

const Alerts = ({ alertOptions = [] }: AlertsProps) => {
    return (
        <List
            sx={{
                width: 372,
                maxHeight: 305,
                overflowY: "auto",
                overflowX: "hidden",
                position: "relative",
            }}
        >
            {alertOptions.map(
                ({ id, alertDate, description, title, type }: AlertDto) => (
                    <ListItem key={id} sx={{ p: 0, mb: 2 }}>
                        <Grid container>
                            <Grid item container spacing={4}>
                                <Grid item xs={1} {...PROPS}>
                                    <CloudOffIcon
                                        fontSize="small"
                                        color={getColorByType(type)}
                                    />
                                </Grid>
                                <Grid item xs={7} {...PROPS}>
                                    <Typography
                                        variant="body1"
                                        sx={{ fontWeight: 500 }}
                                    >
                                        {title}
                                    </Typography>
                                </Grid>
                                <Grid
                                    item
                                    xs={4}
                                    {...PROPS}
                                    justifyContent={"flex-end"}
                                >
                                    <Typography
                                        variant="caption"
                                        color={"textSecondary"}
                                    >
                                        {format(alertDate, "MMM dd ha")}
                                    </Typography>
                                </Grid>
                            </Grid>
                            <Grid item container spacing={3}>
                                <Grid item xs={1} />
                                <Grid item xs={11}>
                                    <Typography
                                        variant="body2"
                                        color={"textSecondary"}
                                        {...PROPS}
                                    >
                                        {description}
                                    </Typography>
                                </Grid>
                            </Grid>
                        </Grid>
                    </ListItem>
                )
            )}
        </List>
    );
};
export default Alerts;
