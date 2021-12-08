import { FormControlCheckboxes } from "..";
import { useCallback, useState } from "react";
import { Grid, Divider, Typography } from "@mui/material";

const LineDivider = () => (
    <Grid item xs={12}>
        <Divider sx={{ width: "90%" }} />
    </Grid>
);

const AlertSettings = () => {
    const [alertList, setAlertList] = useState<Object>({});

    const handleAlertChange = useCallback((key: string, value: boolean) => {
        setAlertList(prevState => ({
            ...prevState,
            [key]: value,
        }));
    }, []);

    return (
        <Grid container spacing={2}>
            <Grid item container>
                <Grid item xs={3}>
                    <Typography variant="h6">Common Events</Typography>
                </Grid>
                <Grid item container xs={9} spacing={1}>
                    {[1, 2].map(i => (
                        <Grid key={`${i}-`} item>
                            <FormControlCheckboxes
                                values={alertList}
                                handleChange={handleAlertChange}
                                checkboxList={[
                                    {
                                        id: 1,
                                        label: `Event Log ${i}`,
                                        value: `event${i}`,
                                    },
                                    {
                                        id: 2,
                                        label: `Alerts ${i}`,
                                        value: `alert${i}`,
                                    },
                                    {
                                        id: 3,
                                        label: `Email ${i}`,
                                        value: `email${i}`,
                                    },
                                ]}
                            />
                        </Grid>
                    ))}
                </Grid>
            </Grid>
            <LineDivider />
            <Grid item container>
                <Grid item xs={3}>
                    <Typography variant="h6">Cloud Events</Typography>
                </Grid>
                <Grid item container xs={9} spacing={1}>
                    {[3, 4].map(i => (
                        <Grid key={`${i}-`} item>
                            <FormControlCheckboxes
                                values={alertList}
                                handleChange={handleAlertChange}
                                checkboxList={[
                                    {
                                        id: 1,
                                        label: `Event Log ${i}`,
                                        value: `event${i}`,
                                    },
                                    {
                                        id: 2,
                                        label: `Alerts ${i}`,
                                        value: `alert${i}`,
                                    },
                                    {
                                        id: 3,
                                        label: `Email ${i}`,
                                        value: `email${i}`,
                                    },
                                ]}
                            />
                        </Grid>
                    ))}
                </Grid>
            </Grid>
            <LineDivider />
        </Grid>
    );
};

export default AlertSettings;
