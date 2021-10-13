import { Box, Grid, Typography } from "@mui/material";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";
import colors from "../../theme/colors";
type PasswordRequirementProps = {
    passwordLength: boolean;
    containLetters: boolean;
    containSpecialCharacter: boolean;
};
const PasswordRequirement = ({
    passwordLength,
    containLetters,
    containSpecialCharacter,
}: PasswordRequirementProps) => {
    return (
        <Box width="100%">
            <Grid container spacing={1}>
                <Grid xs={6} item>
                    <Typography variant="caption">
                        <CheckCircleIcon
                            style={
                                passwordLength
                                    ? { color: colors.green }
                                    : { color: colors.lightGrey }
                            }
                        />
                        {`Be a minimum of 8 characters`}
                    </Typography>
                </Grid>
                <Grid xs={6} item>
                    <Typography variant="caption">
                        <CheckCircleIcon
                            style={
                                containLetters
                                    ? { color: "green" }
                                    : { color: "#70757e" }
                            }
                        />
                        {`Upper & lowercase letters`}
                    </Typography>
                </Grid>
                <Grid xs={6} item>
                    <Typography variant="caption">
                        <CheckCircleIcon
                            style={
                                containSpecialCharacter
                                    ? { color: "green" }
                                    : { color: "#70757e" }
                            }
                        />
                        {`At least one special character`}
                    </Typography>
                </Grid>
            </Grid>
        </Box>
    );
};

export default PasswordRequirement;
