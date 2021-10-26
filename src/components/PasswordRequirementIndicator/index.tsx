import colors from "../../theme/colors";
import { Grid, Typography } from "@mui/material";
import { passwordRules } from "../../constants";
import { makeStyles } from "@mui/styles";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
const useStyles = makeStyles(() => ({
    progressIcon: {
        verticalAlign: "middle",
        display: "inline-flex",
        marginRight: "5px",
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
                                {rules.validator(password) ? (
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
        </section>
    );
};

export default PasswordRequirementIndicator;
