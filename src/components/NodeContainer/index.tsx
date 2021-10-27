import {
    Typography,
    Grid,
    Button,
    ButtonGroup,
    Menu,
    MenuItem,
    Paper,
} from "@mui/material";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import { ActivityIcon } from "../../assets/svg";
import { globalUseStyles } from "../../styles";

const NodeContainer = () => {
    const classes = globalUseStyles();
    return (
        <>
            <Grid item xs={8}>
                <Paper className={classes.GridContainer}>
                    <Grid container spacing={1}>
                        <Grid
                            item
                            xs={12}
                            sm={6}
                            container
                            justifyContent="flex-start"
                        >
                            <Typography variant="h6" color="initial">
                                My Node
                            </Typography>
                        </Grid>
                        <Grid
                            item
                            xs={12}
                            sm={6}
                            container
                            justifyContent="flex-end"
                        >
                            <Button size="small" variant="contained">
                                ADD NODE
                            </Button>
                        </Grid>
                    </Grid>
                </Paper>
            </Grid>
        </>
    );
};

export default NodeContainer;
