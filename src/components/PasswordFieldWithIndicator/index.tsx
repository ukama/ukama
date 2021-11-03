import {
    Grid,
    IconButton,
    InputAdornment,
    TextField,
    Typography,
} from "@mui/material";
import { useState } from "react";
import colors from "../../theme/colors";
import { makeStyles } from "@mui/styles";
import { PasswordRules } from "../../constants";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { Visibility, VisibilityOff } from "@mui/icons-material";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import { FormikErrors, FormikTouched } from "formik";

const useStyles = makeStyles(() => ({
    progressIcon: {
        verticalAlign: "middle",
        display: "inline-flex",
        marginRight: "5px",
    },
}));

type PasswordFieldWithIndicatorProps = {
    onBlur: any;
    value: string;
    label: string;
    fieldStyle: any;
    handleChange: any;
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
        <section className="strength-meter">
            <TextField
                fullWidth
                id="password"
                name="password"
                label={label}
                value={value}
                onBlur={onBlur}
                onChange={handleChange}
                onFocus={() => setFocused(true)}
                InputLabelProps={{ shrink: true }}
                type={togglePassword ? "text" : "password"}
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
                    spacing={"6px"}
                    sx={{
                        display: focused ? "flex" : "none",
                        marginTop: "12px",
                    }}
                >
                    {PasswordRules.map((rules: PasswordRulesProps) => {
                        return (
                            <Grid xs={12} sm={6} item key={rules.id}>
                                <Typography
                                    variant="body2"
                                    id={rules.idLabel}
                                    fontFamily="Work Sans"
                                    letterSpacing="0.4px"
                                >
                                    {rules.validator(value) ? (
                                        <CheckCircleIcon
                                            fontSize="small"
                                            className={classes.progressIcon}
                                            style={{
                                                color: colors.green,
                                            }}
                                        />
                                    ) : (
                                        <CheckCircleOutlineIcon
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
