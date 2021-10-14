import PasswordRequirement from "../PasswordRequirement";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import colors from "../../theme/colors";
import { Grid } from "@mui/material";
import { makeStyles } from "@material-ui/core/styles";
const useStyles = makeStyles(() => ({
    progressIcon: {
        verticalAlign: "middle",
        display: "inline-flex",
    },
}));
type requirementProps = {
    long: boolean;
    containLetters: boolean;
    containSpecialCharacter: boolean;
};
const Requirements = ({
    long,
    containLetters,
    containSpecialCharacter,
}: requirementProps) => {
    const classes = useStyles();

    return (
        <section className="strength-meter">
            <Grid item container>
                <Grid xs={6} item>
                    <PasswordRequirement
                        isvalid={long}
                        invalidMessage={
                            <CheckCircleOutlineIcon
                                className={classes.progressIcon}
                            />
                        }
                        label="Be a minimum of 8 characters"
                        validMessage={
                            <CheckCircleOutlineIcon
                                className={classes.progressIcon}
                                style={{ color: colors.green }}
                            />
                        }
                    />
                </Grid>

                <Grid xs={6} item>
                    <PasswordRequirement
                        isvalid={containLetters}
                        invalidMessage={
                            <CheckCircleOutlineIcon
                                className={classes.progressIcon}
                            />
                        }
                        label="Upper & lowercase letters "
                        validMessage={
                            <CheckCircleOutlineIcon
                                className={classes.progressIcon}
                                style={{ color: colors.green }}
                            />
                        }
                    />
                </Grid>
                <Grid xs={6} item>
                    <PasswordRequirement
                        isvalid={containSpecialCharacter}
                        invalidMessage={
                            <CheckCircleOutlineIcon
                                className={classes.progressIcon}
                            />
                        }
                        label="At least one special character"
                        validMessage={
                            <CheckCircleOutlineIcon
                                className={classes.progressIcon}
                                style={{ color: colors.green }}
                            />
                        }
                    />
                </Grid>
            </Grid>
        </section>
    );
};

export default Requirements;
