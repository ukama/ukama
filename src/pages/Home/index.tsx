import { styled } from "@mui/material/styles";
import {
    Box,
    Typography,
    Grid,
    Button,
    ButtonGroup,
    Menu,
    MenuItem,
    Paper,
} from "@mui/material";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import { ActivityIcon, AlertIcon } from "../../assets/svg";
import React, { useState } from "react";
import { STATS_OPTIONS, STATS_PERIOD } from "../../constants";
import { globalUseStyles } from "../../styles";
const ITEM_HEIGHT = 48;
const Item = styled(Paper)(({ theme }) => ({
    ...theme.typography.body2,
    padding: theme.spacing(1),
    textAlign: "center",
    color: theme.palette.text.secondary,
}));

const Home = () => {
    const classes = globalUseStyles();
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    const [isActive, setIsActive] = useState(false);

    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };
    const handleClose = () => {
        setAnchorEl(null);
    };
    return (
        <>
            <Box sx={{ flexGrow: 1 }}>
                <Grid container spacing={2}>
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
                                    <div
                                        onClick={handleClick}
                                        style={{
                                            display: "flex",
                                            alignItems: "center",
                                            flexWrap: "wrap",
                                            cursor: "pointer",
                                        }}
                                    >
                                        <Typography
                                            variant="h6"
                                            color="initial"
                                        >
                                            Throughput
                                        </Typography>

                                        <ArrowDropDownIcon />
                                    </div>
                                    <Menu
                                        id="long-menu"
                                        MenuListProps={{
                                            "aria-labelledby": "long-button",
                                        }}
                                        anchorEl={anchorEl}
                                        open={open}
                                        onClose={handleClose}
                                        PaperProps={{
                                            style: {
                                                maxHeight: ITEM_HEIGHT * 4.5,
                                                width: "20ch",
                                            },
                                        }}
                                    >
                                        {STATS_OPTIONS.map(option => (
                                            <MenuItem
                                                key={option}
                                                selected={option === "Close"}
                                                onClick={handleClose}
                                            >
                                                {option}
                                            </MenuItem>
                                        ))}
                                    </Menu>
                                </Grid>
                                <Grid
                                    item
                                    xs={12}
                                    sm={6}
                                    container
                                    justifyContent="flex-end"
                                >
                                    <ButtonGroup
                                        size="small"
                                        variant="outlined"
                                        aria-label="outlined primary button group"
                                    >
                                        {STATS_PERIOD.map(period => (
                                            <Button key={period}>
                                                {period}
                                            </Button>
                                        ))}
                                    </ButtonGroup>
                                </Grid>

                                <Grid
                                    item
                                    xs={12}
                                    container
                                    justifyContent="center"
                                >
                                    <Typography variant="body2" color="initial">
                                        Graph goes here!
                                    </Typography>
                                </Grid>

                                <Grid
                                    item
                                    xs={12}
                                    container
                                    justifyContent="center"
                                >
                                    <ActivityIcon />
                                </Grid>
                            </Grid>
                        </Paper>
                    </Grid>

                    <Grid xs={4} item>
                        <Paper className={classes.GridContainer}>
                            <Grid
                                item
                                xs={12}
                                container
                                justifyContent="flex-start"
                            >
                                <Typography variant="h6">Alerts</Typography>{" "}
                            </Grid>
                            <Grid
                                item
                                xs={12}
                                container
                                justifyContent="center"
                            >
                                <AlertIcon />
                            </Grid>
                        </Paper>
                    </Grid>
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
                </Grid>
                <Grid item xs={4}>
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
                                    Residents
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
                                    ACTIVATE
                                </Button>
                            </Grid>
                        </Grid>
                    </Paper>
                </Grid>
            </Box>
        </>
    );
};

export default Home;
