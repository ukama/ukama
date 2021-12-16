import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import EditIcon from "@mui/icons-material/Edit";
import { useState } from "react";
import {
    Box,
    Switch,
    Button,
    Dialog,
    IconButton,
    Typography,
    Stack,
    DialogActions,
    Divider,
    Grid,
} from "@mui/material";
import { InfoIcon } from "../../../assets/svg";
import colors from "../../../theme/colors";

const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        padding: "0px 0px 18px 0px",
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
}));

type BasicDialogProps = {
    userName: string;
    data: string;
    isOpen: boolean;
    userDetailsTitle: string;
    simDetailsTitle: string;
    btnLabel: string;
    handleClose: any;
    isClosable?: boolean;
    closeBtnLabel?: string;
    saveBtnLabel?: string;
};

const UserDetailsDialog = ({
    userName,
    data,
    isOpen,
    simDetailsTitle,
    userDetailsTitle,
    closeBtnLabel,
    saveBtnLabel,
    handleClose,
    isClosable = true,
}: BasicDialogProps) => {
    const classes = useStyles();
    const [roaming, setRoaming] = useState(false);
    const handleRoaming = (e: any) => {
        setRoaming(e.target.checked);
    };
    return (
        <Dialog
            open={isOpen}
            onBackdropClick={() => isClosable && handleClose()}
        >
            <Box
                sx={{
                    width: { xs: "100%", md: "500px" },
                    padding: "16px 8px 8px 24px",
                }}
            >
                <Box className={classes.basicDialogHeaderStyle}>
                    <Stack
                        direction="row"
                        sx={{ alignItems: "center" }}
                        spacing={1}
                    >
                        <Typography variant="h5">{userName}</Typography>
                        <Typography
                            variant="body1"
                            sx={{ color: colors.darkGrey }}
                        >
                            {data}
                        </Typography>
                    </Stack>
                    {isClosable && (
                        <IconButton
                            onClick={handleClose}
                            sx={{ ml: "24px", p: "8px" }}
                        >
                            <CloseIcon />
                        </IconButton>
                    )}
                </Box>

                <Grid container>
                    <Grid item xs={12}>
                        <Typography variant="subtitle2">
                            {userDetailsTitle}
                        </Typography>
                        <Divider />
                    </Grid>
                    <Grid
                        item
                        container
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Grid item>
                            <Stack direction="column">
                                <Typography
                                    variant="caption"
                                    sx={{ color: colors.lightGrey }}
                                >
                                    NAME
                                </Typography>
                                <Typography variant="caption">
                                    John Doe
                                </Typography>
                            </Stack>
                        </Grid>
                        <Grid item>
                            <IconButton>
                                <EditIcon
                                    fontSize="small"
                                    sx={{
                                        color: colors.lightGrey,
                                        position: "relative",
                                        top: "10px",
                                    }}
                                />
                            </IconButton>
                        </Grid>
                    </Grid>
                    <Grid
                        item
                        container
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Grid item>
                            <Stack direction="column">
                                <Typography
                                    variant="caption"
                                    sx={{ color: colors.lightGrey }}
                                >
                                    EMAIL
                                </Typography>
                                <Typography variant="caption">
                                    John@gmail.com
                                </Typography>
                            </Stack>
                        </Grid>
                        <Grid item>
                            <IconButton>
                                <EditIcon
                                    fontSize="small"
                                    sx={{
                                        color: colors.lightGrey,
                                        position: "relative",
                                        top: "10px",
                                    }}
                                />
                            </IconButton>
                        </Grid>
                    </Grid>
                    <Grid
                        item
                        container
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Grid item>
                            <Stack direction="column">
                                <Typography
                                    variant="caption"
                                    sx={{ color: colors.lightGrey }}
                                >
                                    PHONE
                                </Typography>
                                <Typography variant="caption">
                                    (111) 111-1111
                                </Typography>
                            </Stack>
                        </Grid>
                        <Grid item>
                            <IconButton>
                                <EditIcon
                                    fontSize="small"
                                    sx={{
                                        color: colors.lightGrey,
                                        position: "relative",
                                        top: "10px",
                                    }}
                                />
                            </IconButton>
                        </Grid>
                    </Grid>
                </Grid>

                <Grid container sx={{ mt: 1 }}>
                    <Grid item xs={12}>
                        <Typography variant="subtitle2">
                            {simDetailsTitle}
                        </Typography>
                        <Divider />
                    </Grid>
                    <Grid
                        item
                        container
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Grid item>
                            <Stack direction="column">
                                <Typography
                                    variant="caption"
                                    sx={{ color: colors.lightGrey }}
                                >
                                    STATUS
                                </Typography>
                                <Typography variant="caption">
                                    Active
                                </Typography>
                            </Stack>
                        </Grid>
                        <Grid item>
                            <Button
                                size="small"
                                variant="outlined"
                                sx={{
                                    color: "red",
                                    border: "1px solid red",
                                    position: "relative",
                                    top: "12px",
                                }}
                            >
                                PAUSE SERVICE
                            </Button>
                        </Grid>
                    </Grid>
                    <Grid
                        item
                        container
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Grid item>
                            <Stack direction="column">
                                <Typography
                                    variant="caption"
                                    sx={{ color: colors.lightGrey }}
                                >
                                    IMEI NUMBER
                                </Typography>
                                <Typography variant="caption">
                                    1234812374109374139434173470
                                </Typography>
                            </Stack>
                        </Grid>
                    </Grid>
                    <Grid
                        item
                        container
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Grid item>
                            <Stack direction="column">
                                <Typography
                                    variant="caption"
                                    sx={{ color: colors.lightGrey }}
                                >
                                    ICCID
                                </Typography>
                                <Typography variant="caption">
                                    1234812374109374139434173470
                                </Typography>
                            </Stack>
                        </Grid>
                    </Grid>
                    <Grid
                        item
                        container
                        direction="row"
                        justifyContent="space-between"
                    >
                        <Grid item>
                            <Stack direction="column">
                                <Typography
                                    variant="caption"
                                    sx={{ color: colors.lightGrey }}
                                >
                                    ROAMING
                                </Typography>

                                <Stack
                                    direction="row"
                                    alignItems="center"
                                    spacing={1}
                                >
                                    <Typography variant="caption">
                                        Roaming enabled
                                    </Typography>
                                    <InfoIcon />
                                </Stack>
                            </Stack>
                        </Grid>
                        <Grid item>
                            <Switch
                                sx={{
                                    color: colors.lightGrey,
                                    position: "relative",
                                    top: "12px",
                                }}
                                size="small"
                                checked={roaming}
                                onClick={handleRoaming}
                                value="active"
                                inputProps={{
                                    "aria-label": "secondary checkbox",
                                }}
                            />
                        </Grid>
                    </Grid>
                </Grid>

                <DialogActions>
                    <Button
                        onClick={handleClose}
                        color="primary"
                        sx={{ mr: 2 }}
                    >
                        {closeBtnLabel}
                    </Button>
                    <Button
                        onClick={handleClose}
                        color="primary"
                        variant="contained"
                    >
                        {saveBtnLabel}
                    </Button>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default UserDetailsDialog;
