import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import EditableTextField from "../../EditableTextField";
import { useState, useEffect } from "react";
import { GetUserDto } from "../../../generated";
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
    getUser: any;
    isOpen: boolean;
    btnLabel: string;
    handleClose: any;
    isClosable?: boolean;
    closeBtnLabel?: string;
    saveBtnLabel?: string;
    handleSaveSimUser?: any;
    userDetailsTitle: string;
    simDetailsTitle: string;
};

const UserDetailsDialog = ({
    getUser,
    userDetailsTitle,
    isOpen,
    simDetailsTitle,
    closeBtnLabel,
    saveBtnLabel,
    handleClose,
    isClosable = true,
    handleSaveSimUser = () => {
        /* Default empty function */
    },
}: BasicDialogProps) => {
    const classes = useStyles();
    const [roaming, setRoaming] = useState(false);
    const [name, setName] = useState<any>();
    const [currentEmail, setCurrentEmail] = useState<any>();
    const [phone, setPhone] = useState<any>();
    const [isOff, setIsOff] = useState(true);
    const [formData, setFormData] = useState<any>();
    const handleRoaming = (e: any) => {
        setRoaming(e.target.checked);
    };
    const handleSave = () => {
        const datas = {
            name: name,
            email: currentEmail,
            phone: phone,
            service: isOff,
            roaming: roaming,
        };
        setFormData(datas);
    };
    useEffect(() => {
        handleSaveSimUser(formData);
    }, [formData]);

    return (
        <>
            {getUser.map(
                ({
                    id,
                    status,
                    name,
                    eSimNumber,
                    iccid,
                    email,
                    phone,
                    roaming,
                    dataPlan,
                    dataUsage,
                }: GetUserDto) => (
                    <Dialog
                        key={id}
                        open={isOpen}
                        onBackdropClick={() => isClosable && handleClose()}
                        hideBackdrop
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
                                    <Typography variant="h5">{name}</Typography>
                                    <Typography
                                        variant="body1"
                                        sx={{ color: colors.darkGrey }}
                                    >
                                        - {dataUsage} GB data used, {dataPlan}
                                        free GB left
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
                                        {userDetailsTitle} {id}
                                    </Typography>
                                    <Divider sx={{ width: "450px" }} />
                                </Grid>
                            </Grid>
                            <Grid item container>
                                <Grid item container xs={12} spacing={2}>
                                    <Grid item xs={9}>
                                        <EditableTextField
                                            value={name || ""}
                                            label={"NAME"}
                                            handleOnChange={(value: string) => {
                                                setName(value);
                                            }}
                                        />
                                    </Grid>
                                    <Grid item xs={9}>
                                        <EditableTextField
                                            value={email || ""}
                                            label={"EMAIL"}
                                            handleOnChange={(value: string) => {
                                                setCurrentEmail(value);
                                            }}
                                        />
                                    </Grid>
                                    <Grid item xs={9}>
                                        <EditableTextField
                                            value={phone || ""}
                                            label={"PHONE"}
                                            handleOnChange={(value: string) => {
                                                setPhone(value);
                                            }}
                                        />
                                    </Grid>
                                </Grid>
                            </Grid>
                            <Grid container sx={{ mt: 1 }}>
                                <Grid item xs={12}>
                                    <Typography variant="subtitle2">
                                        {simDetailsTitle}
                                    </Typography>
                                    <Divider sx={{ width: "450px" }} />
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
                                                sx={{
                                                    color: colors.lightGrey,
                                                }}
                                            >
                                                STATUS
                                            </Typography>
                                            <Typography variant="body2">
                                                {status}
                                            </Typography>
                                        </Stack>
                                    </Grid>
                                    <Grid item>
                                        <Button
                                            onClick={() => setIsOff(!isOff)}
                                            size="small"
                                            variant="outlined"
                                            sx={{
                                                color: "red",
                                                border: "1px solid red",
                                                position: "relative",
                                                top: "12px",
                                                mr: 2,
                                                "&:hover": {
                                                    borderColor: "red",
                                                },
                                            }}
                                        >
                                            {isOff
                                                ? "PAUSE SERVICE"
                                                : "RESUME SERVICE "}
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
                                                sx={{
                                                    color: colors.lightGrey,
                                                }}
                                            >
                                                {iccid}
                                            </Typography>
                                            <Typography variant="body2">
                                                {eSimNumber}
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
                                                sx={{
                                                    color: colors.lightGrey,
                                                }}
                                            >
                                                ICCID
                                            </Typography>
                                            <Typography variant="body2">
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
                                                sx={{
                                                    color: colors.lightGrey,
                                                }}
                                            >
                                                ROAMING
                                            </Typography>

                                            <Stack
                                                direction="row"
                                                alignItems="center"
                                                spacing={1}
                                            >
                                                <Typography variant="body2">
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
                                                mr: 1,
                                            }}
                                            size="small"
                                            checked={roaming}
                                            onClick={handleRoaming}
                                            value="active"
                                            inputProps={{
                                                "aria-label":
                                                    "secondary checkbox",
                                            }}
                                        />
                                    </Grid>
                                </Grid>
                            </Grid>

                            <DialogActions>
                                <Button
                                    onClick={handleClose}
                                    color="primary"
                                    sx={{ mr: 2, justifyItems: "center" }}
                                >
                                    {closeBtnLabel}
                                </Button>
                                <Button
                                    onClick={handleSave}
                                    color="primary"
                                    variant="contained"
                                >
                                    {saveBtnLabel}
                                </Button>
                            </DialogActions>
                        </Box>
                    </Dialog>
                )
            )}
        </>
    );
};

export default UserDetailsDialog;
