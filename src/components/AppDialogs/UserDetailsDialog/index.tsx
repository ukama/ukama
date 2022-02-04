import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import EditableTextField from "../../EditableTextField";
import {
    Tooltip,
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
    DialogContent,
} from "@mui/material";
import colors from "../../../theme/colors";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import { ContainerJustifySpaceBtw } from "../../../styles";
import { GetUserDto, Get_User_Status_Type } from "../../../generated";

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
    user: GetUserDto;
    isOpen: boolean;
    setUserForm: any;
    isClosable?: boolean;
    handleClose: Function;
    closeBtnLabel?: string;
    saveBtnLabel?: string;
    handleSaveSimUser?: any;
    userDetailsTitle: string;
    simDetailsTitle: string;
};

const UserDetailsDialog = ({
    user,
    isOpen,
    setUserForm,
    handleClose,
    saveBtnLabel,
    closeBtnLabel,
    simDetailsTitle,
    userDetailsTitle,
    isClosable = true,
    handleSaveSimUser = () => {
        /* Default empty function */
    },
}: BasicDialogProps) => {
    const classes = useStyles();
    const {
        id,
        name,
        email,
        phone,
        iccid,
        status,
        roaming,
        dataPlan,
        dataUsage,
        eSimNumber,
    } = user;

    return (
        <Dialog
            key={id}
            open={isOpen}
            hideBackdrop
            onBackdropClick={() => isClosable && handleClose()}
        >
            <Box
                sx={{
                    width: { xs: "100%", md: "500px" },
                    padding: "16px 24px",
                }}
            >
                <Box className={classes.basicDialogHeaderStyle}>
                    <Stack
                        direction="row"
                        sx={{ alignItems: "center" }}
                        spacing={1}
                    >
                        <Typography variant="h5">{name}</Typography>
                    </Stack>
                    {isClosable && (
                        <IconButton
                            onClick={() => handleClose()}
                            sx={{ ml: "24px", p: "8px" }}
                        >
                            <CloseIcon />
                        </IconButton>
                    )}
                </Box>
                <DialogContent sx={{ padding: 0, mb: 4 }}>
                    <Grid>
                        <Grid container>
                            <Grid item xs={12}>
                                <Typography variant="subtitle2">
                                    {userDetailsTitle}
                                </Typography>
                                <Divider />
                            </Grid>
                        </Grid>
                        <Grid item container spacing={1}>
                            <Grid item xs={12}>
                                <Stack direction="row" spacing={1}>
                                    <Typography
                                        variant="body1"
                                        sx={{ color: colors.vulcan }}
                                    >
                                        {dataUsage} GB data used,
                                    </Typography>
                                    <Typography
                                        variant="body1"
                                        sx={{ color: colors.vulcan }}
                                    >
                                        {dataPlan}
                                    </Typography>
                                    <Typography
                                        variant="body1"
                                        sx={{ color: colors.vulcan }}
                                    >
                                        free GB left
                                    </Typography>
                                </Stack>
                            </Grid>
                            <Grid item xs={9}>
                                <EditableTextField
                                    value={name || ""}
                                    label={"NAME"}
                                    handleOnChange={(value: string) =>
                                        setUserForm({ ...user, name: value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={9}>
                                <EditableTextField
                                    value={email || ""}
                                    label={"EMAIL"}
                                    handleOnChange={(value: string) =>
                                        setUserForm({ ...user, email: value })
                                    }
                                />
                            </Grid>
                            <Grid item xs={9}>
                                <EditableTextField
                                    value={phone || ""}
                                    label={"PHONE"}
                                    handleOnChange={(value: string) =>
                                        setUserForm({ ...user, email: value })
                                    }
                                />
                            </Grid>
                        </Grid>
                        <Grid container sx={{ mt: 1 }} spacing={1}>
                            <Grid item xs={12}>
                                <Typography variant="subtitle2">
                                    {simDetailsTitle}
                                </Typography>
                                <Divider />
                            </Grid>
                            <Grid item container>
                                <Grid item xs={12}>
                                    <Typography
                                        variant="caption"
                                        sx={{
                                            color: colors.lightGrey,
                                        }}
                                    >
                                        STATUS
                                    </Typography>
                                </Grid>
                                <Grid item xs={12}>
                                    <ContainerJustifySpaceBtw>
                                        <Typography variant="body2">
                                            {status}
                                        </Typography>
                                        <Button
                                            size="small"
                                            color="error"
                                            variant="outlined"
                                            onClick={() =>
                                                setUserForm({
                                                    ...user,
                                                    status:
                                                        status ===
                                                        Get_User_Status_Type.Active
                                                            ? Get_User_Status_Type.Inactive
                                                            : Get_User_Status_Type.Active,
                                                })
                                            }
                                        >
                                            {status ===
                                            Get_User_Status_Type.Active
                                                ? "PAUSE SERVICE"
                                                : "RESUME SERVICE "}
                                        </Button>
                                    </ContainerJustifySpaceBtw>
                                </Grid>
                            </Grid>
                            <Grid item container xs={12}>
                                <Grid item xs={12}>
                                    <Typography
                                        variant="caption"
                                        sx={{
                                            color: colors.lightGrey,
                                        }}
                                    >
                                        IMEI NUMBER
                                    </Typography>
                                </Grid>
                                <Grid item xs={12}>
                                    <Typography variant="body2">
                                        {eSimNumber}
                                    </Typography>
                                </Grid>
                            </Grid>
                            <Grid item container xs={12}>
                                <Grid item xs={12}>
                                    <Typography
                                        variant="caption"
                                        sx={{
                                            color: colors.lightGrey,
                                        }}
                                    >
                                        ICCID
                                    </Typography>
                                </Grid>
                                <Grid item xs={12}>
                                    <Typography variant="body2">
                                        {iccid}
                                    </Typography>
                                </Grid>
                            </Grid>
                            <Grid item container xs={12}>
                                <ContainerJustifySpaceBtw>
                                    <Typography
                                        variant="caption"
                                        sx={{
                                            color: colors.lightGrey,
                                            alignSelf: "end",
                                        }}
                                    >
                                        ROAMING
                                        <Tooltip
                                            title="Explain roaming policy for CS folks."
                                            placement="right"
                                            arrow
                                        >
                                            <IconButton>
                                                <InfoIcon />
                                            </IconButton>
                                        </Tooltip>
                                    </Typography>
                                    <Switch
                                        size="small"
                                        value="active"
                                        checked={roaming}
                                        onClick={(e: any) =>
                                            setUserForm({
                                                ...user,
                                                roaming: e.target.checked,
                                            })
                                        }
                                    />
                                </ContainerJustifySpaceBtw>
                            </Grid>
                        </Grid>
                    </Grid>
                </DialogContent>
                <DialogActions sx={{ padding: 0 }}>
                    <Button
                        onClick={() => handleClose()}
                        sx={{ mr: 2, justifyItems: "center" }}
                    >
                        {closeBtnLabel}
                    </Button>
                    <Button onClick={handleSaveSimUser} variant="contained">
                        {saveBtnLabel}
                    </Button>
                </DialogActions>
            </Box>
        </Dialog>
    );
};

export default UserDetailsDialog;
