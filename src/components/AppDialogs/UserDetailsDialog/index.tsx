import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import EditableTextField from "../../EditableTextField";
import { useState, useEffect } from "react";
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
    handleSaveSimUser?: any;
    singleUserdetails?: any;
};

const UserDetailsDialog = ({
    userName,
    isOpen,
    simDetailsTitle,
    userDetailsTitle,
    closeBtnLabel,
    saveBtnLabel,
    handleClose,
    singleUserdetails,
    isClosable = true,
    handleSaveSimUser = () => {
        /* Default empty function */
    },
}: BasicDialogProps) => {
    const classes = useStyles();
    const [roaming, setRoaming] = useState(false);
    const [name, setName] = useState<string>("JohnDoe");
    const [email, setEmail] = useState<string>("john@ukama.com");
    const [phone, setPhone] = useState<string>("(111) 111-1111");
    const [isOff, setIsOff] = useState(true);
    const [formData, setFormData] = useState<any>();
    const handleRoaming = (e: any) => {
        setRoaming(e.target.checked);
    };
    const handleSave = () => {
        const data = {
            name: name,
            email: email,
            phone: phone,
            service: isOff,
            roaming: roaming,
        };
        setFormData(data);
    };
    useEffect(() => {
        handleSaveSimUser(formData);
    }, [formData]);

    return (
        <>
            {singleUserdetails &&
                singleUserdetails.map(({ data }: any) => (
                    <Dialog
                        key={data.id}
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
                                    <Typography variant="h5">
                                        {userName}
                                    </Typography>
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
                                    <Divider sx={{ width: "450px" }} />
                                </Grid>
                            </Grid>
                            <Grid item container>
                                <Grid item container xs={12} spacing={2}>
                                    <Grid item xs={9}>
                                        <EditableTextField
                                            value={name}
                                            label={"NAME"}
                                            handleOnChange={(value: string) => {
                                                setName(value);
                                            }}
                                        />
                                    </Grid>
                                    <Grid item xs={9}>
                                        <EditableTextField
                                            value={email}
                                            label={"EMAIL"}
                                            handleOnChange={(value: string) => {
                                                setEmail(value);
                                            }}
                                        />
                                    </Grid>
                                    <Grid item xs={9}>
                                        <EditableTextField
                                            value={phone}
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
                                                Active
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
                                                IMEI NUMBER
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
                ))}
        </>
    );
};

export default UserDetailsDialog;
