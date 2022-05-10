import { colors } from "../../theme";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import { Stack, Typography, Grid, TextField, Switch } from "@mui/material";

interface IPhysicalSimform {
    formData?: any;
    setFormData?: any;
    description: String;
}

const PhysicalSimform = ({
    description,
    formData,
    setFormData,
}: IPhysicalSimform) => {
    const gclasses = globalUseStyles();
    return (
        <Grid container spacing={3}>
            <Grid item xs={12}>
                <Typography variant="body1">{description}</Typography>
            </Grid>
            <Grid container item xs={12} spacing={1}>
                <Grid item xs={6}>
                    <TextField
                        fullWidth
                        label={"ICCID"}
                        value={formData.iccid}
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            classes: {
                                input: gclasses.inputFieldStyle,
                            },
                        }}
                        onChange={e =>
                            setFormData({ ...formData, iccid: e.target.value })
                        }
                    />
                </Grid>
                <Grid item xs={6}>
                    <TextField
                        fullWidth
                        label={"SECURITY CODE"}
                        value={formData.code}
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            classes: {
                                input: gclasses.inputFieldStyle,
                            },
                        }}
                        onChange={e =>
                            setFormData({ ...formData, code: e.target.value })
                        }
                    />
                </Grid>
            </Grid>
        </Grid>
    );
};

export default PhysicalSimform;
