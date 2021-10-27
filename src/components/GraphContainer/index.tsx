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
import React, { useState } from "react";
import { STATS_OPTIONS, STATS_PERIOD } from "../../constants";
import { globalUseStyles } from "../../styles";
const ITEM_HEIGHT = 48;
const GraphContainer = () => {
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
                                <Typography variant="h6" color="initial">
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
                            {isActive ? (
                                <ButtonGroup
                                    size="small"
                                    variant="outlined"
                                    aria-label="outlined primary button group"
                                >
                                    {STATS_PERIOD.map(period => (
                                        <Button key={period}>{period}</Button>
                                    ))}
                                </ButtonGroup>
                            ) : null}
                        </Grid>
                        {isActive ? (
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
                        ) : (
                            <Grid
                                item
                                xs={12}
                                container
                                justifyContent="center"
                            >
                                <ActivityIcon />
                            </Grid>
                        )}
                    </Grid>
                </Paper>
            </Grid>
        </>
    );
};

export default GraphContainer;
