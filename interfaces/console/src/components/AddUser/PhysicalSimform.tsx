import { colors } from "../../theme";
import { globalUseStyles } from "../../styles";
import ErrorIcon from "@mui/icons-material/Error";
import { Typography, Grid, TextField, Alert } from "@mui/material";

interface IPhysicalSimform {
    formData?: any;
    formError: string;
    setFormData?: any;
    description: String;
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
            {formError && (
                <Grid item xs={12}>
                    <Alert
                        sx={{
                            mb: 1,
                            color: colors.black,
                        }}
                        severity={"error"}
                        icon={<ErrorIcon sx={{ color: colors.red }} />}
                    >
                        {formError}
                    </Alert>
                </Grid>
            )}
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
                    />
                </Grid>
            </Grid>
        </Grid>
    );
};

export default PhysicalSimform;
