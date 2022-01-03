import {
    Grid,
    Switch,
    Divider,
    MenuItem,
    TextField,
    Typography,
    FormControlLabel,
} from "@mui/material";
import { useState } from "react";
import { RF_NODES } from "../../constants";

const LineDivider = () => (
    <Grid item xs={12}>
        <Divider sx={{ width: "100%" }} />
    </Grid>
);

const NodeSettings = () => {
    const [rf, setRf] = useState(false);
    const [node, setNodes] = useState("off");

    return (
        <Grid container spacing={2}>
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
                                    checked={rf}
                                    onChange={e => setRf(e.target.checked)}
                                />
                            }
                            label="Enable roaming for all"
                            sx={{ typography: "body1" }}
                        />
                    </Grid>
                    <Grid item xs={12} sm={10} md={8}>
                        <TextField
                            select
                            id="nodes"
                            InputProps={{
                                disabled: !rf,
                                disableUnderline: true,
                            }}
                            value={node}
                            sx={{ mt: "18px" }}
                            variant={"standard"}
                            onChange={e => setNodes(e.target.value)}
                        >
                            {RF_NODES.map(({ value, text }: any) => (
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

export default NodeSettings;
