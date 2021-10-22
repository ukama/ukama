import colors from "../../theme/colors";
import { Grid, Typography } from "@mui/material";
import { passwordRules } from "../../constants";
import { makeStyles } from "@mui/styles";
import CheckCircleRoundedIcon from "@mui/icons-material/CheckCircleRounded";
import CheckCircleSharpIcon from "@mui/icons-material/CheckCircleSharp";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
const useStyles = makeStyles(() => ({
    progressIcon: {
        verticalAlign: "middle",
        display: "inline-flex",
    },
}));
type PasswordRequirementIndicatorProps = {
    password: string;
};
type PasswordRulesProps = {
    id: number;
    label: string;
    validator: Function;
};

const PasswordRequirementIndicator = ({
    password,
}: PasswordRequirementIndicatorProps) => {
    const classes = useStyles();

    return (
        <section className="strength-meter">
            <Grid
                item
                container
                sx={{ display: password.length < 1 ? "none" : "flex" }}
            >
                {passwordRules.map((rules: PasswordRulesProps) => {
                    return (
                        <Grid xs={6} item key={rules.id}>
                            <Typography variant="body2">
                                <CheckCircleIcon
                                    fontSize="small"
                                    className={classes.progressIcon}
                                    style={{
                                        color: rules.validator(password)
                                            ? colors.green
                                            : colors.grey,
                                    }}
                                />
                                {rules.label}
                            </Typography>
                        </Grid>
                    );
                })}
            </Grid>
        </section>
    );
};

export default PasswordRequirementIndicator;
