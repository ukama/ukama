import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import colors from "../../theme/colors";
import { Grid, Typography } from "@mui/material";
import { passwordRules } from "../../constants";
import { makeStyles } from "@mui/styles";
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
            <Grid item container>
                {passwordRules.map((rules: PasswordRulesProps) => {
                    return (
                        <Grid xs={6} item key={rules.id}>
                            <Typography variant="body2">
                                <CheckCircleOutlineIcon
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
