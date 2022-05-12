import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import { SIM_TYPES } from "../../constants";
import { HorizontalContainerJustify } from "../../styles";
import ArrowIcon from "@mui/icons-material/ArrowForwardIos";
import { Stack, Typography, Box, IconButton, Grid } from "@mui/material";

const useStyles = makeStyles(() => ({
    cardStyle: {
        cursor: "pointer",
        borderRadius: "4px",
        background: colors.white,
        border: "1px solid rgba(0, 0, 0, 0.23)",
        "& :hover": {
            borderRadius: "4px",
            background: colors.solitude,
        },
    },
}));

interface IChooseSim {
    description: string;
    handleSimType: Function;
}

const ChooseSim = ({ description, handleSimType }: IChooseSim) => {
    const classes = useStyles({});
    return (
        <Grid container spacing={3} mb={2}>
            <Grid item xs={12}>
                <Typography variant="body1">{description}</Typography>
            </Grid>
            <Grid item xs={12}>
                <Stack spacing={2}>
                    {SIM_TYPES.map(({ id, title, type }) => (
                        <Box
                            key={id}
                            component="div"
                            className={classes.cardStyle}
                            onClick={() => handleSimType({ type })}
                        >
                            <HorizontalContainerJustify sx={{ p: 2 }}>
                                <Typography variant="body1">{title}</Typography>
                                <IconButton sx={{ p: 0 }}>
                                    <ArrowIcon sx={{ height: 18 }} />
                                </IconButton>
                            </HorizontalContainerJustify>
                        </Box>
                    ))}
                </Stack>
            </Grid>
        </Grid>
    );
};

export default ChooseSim;
