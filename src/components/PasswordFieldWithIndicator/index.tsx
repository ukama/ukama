import {
    Grid,
    IconButton,
    InputAdornment,
    TextField,
    Typography,
} from "@mui/material";
import {
    Visibility,
    CheckCircle,
    VisibilityOff,
    CheckCircleOutline,
} from "@mui/icons-material";
import { useState } from "react";
import colors from "../../theme/colors";
import { makeStyles } from "@mui/styles";
import { PasswordRules } from "../../constants";
import { FormikErrors, FormikTouched } from "formik";

const useStyles = makeStyles(() => ({
    progressIcon: {
        verticalAlign: "middle",
        display: "inline-flex",
        marginRight: "5px",
    },
}));

type PasswordFieldWithIndicatorProps = {
    value: string;
    label: string;
    fieldStyle: any;
    onBlur: Function;
    handleChange: Function;
    withIndicator: boolean;
    errors: FormikErrors<any>;
    touched: FormikTouched<any>;
};

type PasswordRulesProps = {
    id: number;
    label: string;
    idLabel: string;
    validator: Function;
};

const PasswordFieldWithIndicator = ({
    value,
    errors,
    label,
    onBlur,
    touched,
    fieldStyle,
    handleChange,
    withIndicator,
}: PasswordFieldWithIndicatorProps) => {
    const classes = useStyles();
    const [focused, setFocused] = useState(false);
    const [togglePassword, setTogglePassword] = useState(false);
    const handleTogglePassword = () => {
        setTogglePassword(prev => !prev);
    };
    return (
        <section>
            <TextField
                fullWidth
                id="password"
                name="password"
                label={label}
                value={value}
                onBlur={e => onBlur(e)}
                onChange={e => handleChange(e)}
                onFocus={() => setFocused(true)}
                InputLabelProps={{ shrink: true }}
                type={togglePassword ? "text" : "password"}
                sx={{ mt: "12px", mb: withIndicator ? "12px" : "0px" }}
                error={touched.password && Boolean(errors.password)}
                helperText={touched.password && errors.password}
                InputProps={{
                    classes: { input: fieldStyle },
                    endAdornment: (
                        <InputAdornment position="end">
                            <IconButton
                                edge="end"
                                onClick={handleTogglePassword}
                            >
                                {togglePassword ? (
                                    <VisibilityOff />
                                ) : (
                                    <Visibility />
                                )}
                            </IconButton>
                        </InputAdornment>
                    ),
                }}
            />
            {withIndicator && (
                <Grid
                    item
                    container
                    sx={{
                        display: focused ? "flex" : "none",
                        mb: "12px",
                    }}
                >
                    {PasswordRules.map((rules: PasswordRulesProps) => {
                        return (
                            <Grid xs={12} sm={6} item key={rules.id}>
                                <Typography
                                    variant="caption"
                                    id={rules.idLabel}
                                    fontFamily="Work Sans"
                                    letterSpacing="0.4px"
                                >
                                    {rules.validator(value) ? (
                                        <CheckCircle
                                            fontSize="small"
                                            className={classes.progressIcon}
                                            style={{
                                                color: colors.green,
                                            }}
                                        />
                                    ) : (
                                        <CheckCircleOutline
                                            fontSize="small"
                                            className={classes.progressIcon}
                                            style={{
                                                color: colors.grey,
                                            }}
                                        />
                                    )}

                                    {rules.label}
                                </Typography>
                            </Grid>
                        );
                    })}
                </Grid>
            )}
        </section>
    );
};

export default PasswordFieldWithIndicator;
