import {
    Grid,
    Radio,
    Switch,
    Divider,
    MenuItem,
    TextField,
    Typography,
    RadioGroup,
    FormControlLabel,
} from "@mui/material";
import { useState } from "react";
import { ROAMING_SELECT } from "../../constants";

const LineDivider = () => (
    <Grid item xs={12}>
        <Divider sx={{ width: "90%" }} />
    </Grid>
);

const NetworkSettings = () => {
    const [roaming, setRoaming] = useState(false);
    const [esim, setEsim] = useState("all");

    const handleTimezoneChange = (event: any) => {
        setEsim(event.target.value);
    };
    return (
        <Grid container spacing={2}>
            <Grid item container spacing={2}>
                <Grid item xs={12} md={4}>
                    <Typography variant="h6">Network Name</Typography>
                </Grid>
                <Grid item xs={12} md={8}>
                    <Grid item xs={12} sm={10} md={8}>
                        <Typography
                            variant="body1"
                            sx={{
                                mb: "18px",
                                lineHeight: "19px",
                            }}
                        >
                            This is the name that shows up on xyz. You can edit
                            this again at any point.
                        </Typography>
                    </Grid>
                    <Grid item xs={12} sm={10} md={8}>
                        <TextField
                            fullWidth
                            id="name"
                            name="name"
                            disabled={true}
                            variant="standard"
                            value={"Network X"}
                            label={"NETWORK NAME"}
                            InputLabelProps={{ shrink: true }}
                            InputProps={{
                                disableUnderline: true,
                            }}
                        />
                    </Grid>
                </Grid>
            </Grid>
            <LineDivider />
            <Grid item container spacing={2}>
                <Grid item xs={12} md={4}>
                    <Typography variant="h6">Network Visibility(?)</Typography>
                </Grid>
                <Grid item container xs={12} md={8}>
                    <Grid item xs={12} sm={10} md={8}>
                        <Typography
                            variant="body1"
                            sx={{
                                mb: "18px",
                                lineHeight: "19px",
                            }}
                        >
                            Policy regarding network switching & explain how
                            it’ll change after hardware is actually shipped.
                        </Typography>
                    </Grid>
                    <Grid item xs={12} sm={10} md={8}>
                        <RadioGroup
                            aria-label="networkType"
                            defaultValue="private"
                            name="radio-buttons-group"
                        >
                            <FormControlLabel
                                value="public"
                                control={<Radio />}
                                label="Public Network"
                            />
                            <FormControlLabel
                                value="private"
                                control={<Radio />}
                                label="Private Network"
                            />
                        </RadioGroup>
                    </Grid>
                </Grid>
            </Grid>
            <LineDivider />
            <Grid item container spacing={2}>
                <Grid item xs={12} md={4}>
                    <Typography variant="h6">Roaming Options</Typography>
                </Grid>
                <Grid item container xs={12} md={8}>
                    <Grid item xs={12} sm={10} md={8}>
                        <Typography
                            variant="body1"
                            sx={{
                                mb: "18px",
                                lineHeight: "19px",
                            }}
                        >
                            Explanation of roaming & its rates. Your temporary
                            eSIM has roaming enabled by default, and cannot be
                            disabled.
                        </Typography>
                    </Grid>
                    <Grid item xs={12} sm={10} md={8}>
                        <FormControlLabel
                            control={
                                <Switch
                                    checked={roaming}
                                    onChange={e => setRoaming(e.target.checked)}
                                />
                            }
                            label="Enable roaming for all"
                            sx={{ typography: "body1" }}
                        />
                    </Grid>
                    <Grid item xs={12} sm={10} md={8}>
                        <TextField
                            select
                            id="eSims"
                            InputProps={{
                                disabled: !roaming,
                                disableUnderline: true,
                            }}
                            value={esim}
                            variant={"standard"}
                            sx={{ mt: "18px" }}
                            onChange={handleTimezoneChange}
                        >
                            {ROAMING_SELECT.map(({ value, text }: any) => (
                                <MenuItem key={value} value={value}>
                                    <Typography
                                        variant="body2"
                                        sx={{ fontWeight: 500 }}
                                    >
                                        {text}
                                    </Typography>
                                </MenuItem>
                            ))}
                        </TextField>
                    </Grid>
                </Grid>
            </Grid>
            <LineDivider />
        </Grid>
    );
};

export default NetworkSettings;
