import colors from "../../theme/colors";
import { makeStyles } from "@mui/styles";
import {
    Grid,
    IconButton,
    InputAdornment,
    TextField,
    Typography,
} from "@mui/material";
import { PasswordRules } from "../../constants";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import { useState } from "react";
import { Visibility, VisibilityOff } from "@mui/icons-material";
const useStyles = makeStyles(() => ({
    progressIcon: {
        verticalAlign: "middle",
        display: "inline-flex",
        marginRight: "5px",
    },
}));
type PasswordFieldWithIndicatorProps = {
    value: any;
    errors: any;
    touched: any;
    label: string;
    fieldStyle: any;
    handleChange: any;
    withIndicator: boolean;
};

type PasswordRulesProps = {
    id: number;
    label: string;
    validator: Function;
};

const PasswordFieldWithIndicator = ({
    value,
    errors,
    label,
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
                onFocus={() => setFocused(true)}
                onChange={event => {
                    handleChange(event);
                }}
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
                                    <Visibility />
                                ) : (
                                    <VisibilityOff />
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
                        marginTop: "8px",
                    }}
                >
                    {PasswordRules.map((rules: PasswordRulesProps) => {
                        return (
                            <Grid
                                xs={6}
                                item
                                key={rules.id}
                                id="passwordIndicator"
                            >
                                <Typography variant="body2">
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
