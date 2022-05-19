import { globalUseStyles } from "../../styles";
import { Typography, Grid, TextField } from "@mui/material";

interface IPhysicalSimform {
    formData?: any;
    formError: string;
    setFormData?: any;
    description: string;
}

const PhysicalSimform = ({
    formError,
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
                        required
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
                        error={Boolean(formError)}
                    />
                </Grid>
                <Grid item xs={6}>
                    <TextField
                        required
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
                        helperText={formError && formError}
                        error={Boolean(formError)}
                    />
                </Grid>
            </Grid>
        </Grid>
    );
};

export default PhysicalSimform;
