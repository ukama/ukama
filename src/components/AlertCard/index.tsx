import { format } from "date-fns";
import { colors } from "../../theme";
import { CloudOffIcon } from "../../assets/svg";
import { AlertDto, Alert_Type } from "../../generated";
import { Typography, Card, Grid, List, ListItem, Stack } from "@mui/material";

type AlertCardProps = {
    alertOptions: AlertDto[] | undefined;
};

const getColorByType = (type: Alert_Type) =>
    type === Alert_Type.Error
        ? colors.red
        : type === Alert_Type.Warning
        ? colors.yellow
        : colors.green;

const AlertCard = ({ alertOptions = [] }: AlertCardProps) => {
    return (
        <>
            <List
                sx={{
                    pr: "4px",
                    maxHeight: 305,
                    overflow: "auto",
                    position: "relative",
                }}
            >
                {alertOptions.map(
                    ({ id, alertDate, description, title, type }: AlertDto) => (
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
                                            <CloudOffIcon
                                                color={getColorByType(type)}
                                            />
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
                                            {format(
                                                new Date(alertDate),
                                                "yyyy-MM-dd"
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
                            </Card>
                        </ListItem>
                    )
                )}
            </List>
        </>
    );
};
export default AlertCard;
