import { colors } from "../../theme";
import { ContainerJustifySpaceBtw, globalUseStyles } from "../../styles";
import { Stack, Typography, Grid, TextField, Switch } from "@mui/material";

interface IUserform {
    formData: any;
    setFormData: any;
    description: String;
}

const Userform = ({ description, formData, setFormData }: IUserform) => {
    const gclasses = globalUseStyles();
    return (
        <Grid container spacing={2}>
            <Grid item xs={12} mb={1}>
                <Typography variant="body1">{description}</Typography>
            </Grid>
            <Grid item xs={12}>
                <TextField
                    fullWidth
                    label={"NAME"}
                    value={formData.name}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{
                        classes: {
                            input: gclasses.inputFieldStyle,
                        },
                    }}
                    onChange={e =>
                        setFormData({ ...formData, name: e.target.value })
                    }
                />
            </Grid>
            <Grid item xs={12}>
                <TextField
                    fullWidth
                    label={"EMAIL"}
                    value={formData.email}
                    InputLabelProps={{ shrink: true }}
                    InputProps={{
                        classes: {
                            input: gclasses.inputFieldStyle,
                        },
                    }}
                    onChange={e =>
                        setFormData({ ...formData, email: e.target.value })
                    }
                />
            </Grid>
            <Grid item xs={12}>
                <ContainerJustifySpaceBtw sx={{ alignItems: "end" }}>
                    <Stack display="flex" alignItems="flex-start">
                        <Typography variant="caption" color={colors.black54}>
                            ROAMING
                        </Typography>
                        <Typography variant="body1">
                            Roaming allows user to do xyz. Insert billing
                            information.
                        </Typography>
                    </Stack>
                    <Switch
                        size="small"
                        value="active"
                        checked={formData.roaming}
                        onChange={e =>
                            setFormData({
                                ...formData,
                                roaming: e.target.checked,
                            })
                        }
                    />
                </ContainerJustifySpaceBtw>
            </Grid>
        </Grid>
    );
};

export default Userform;
