import {
    Box,
    Grid,
    Stack,
    Button,
    Dialog,
    Switch,
    Divider,
    Tooltip,
    IconButton,
    Typography,
    DialogActions,
    DialogContent,
    CircularProgress,
} from "@mui/material";
import { makeStyles } from "@mui/styles";
import CloseIcon from "@mui/icons-material/Close";
import EditableTextField from "../../EditableTextField";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import { CenterContainer, ContainerJustifySpaceBtw } from "../../../styles";
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
    loading: boolean;
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
    loading = true,
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
            onBackdropClick={() => isClosable && handleClose()}
        >
            {loading ? (
                <CenterContainer>
                    <CircularProgress />
                </CenterContainer>
            ) : (
                <Box
                    component="div"
                    sx={{
                        width: { xs: "100%", md: "500px" },
                        padding: "16px 24px",
                    }}
                >
                    <Box
                        component="div"
                        className={classes.basicDialogHeaderStyle}
                    >
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
                                        <Typography variant="body1">
                                            {dataUsage} MB data used,
                                        </Typography>
                                        <Typography variant="body1">
                                            {dataPlan - dataUsage}
                                        </Typography>
                                        <Typography variant="body1">
                                            free MB left
                                        </Typography>
                                    </Stack>
                                </Grid>
                                <Grid item xs={9}>
                                    <EditableTextField
                                        value={name || "-"}
                                        label={"NAME"}
                                        handleOnChange={(value: string) =>
                                            setUserForm({
                                                ...user,
                                                name: value,
                                            })
                                        }
                                    />
                                </Grid>
                                <Grid item xs={9}>
                                    <EditableTextField
                                        value={email || "-"}
                                        label={"EMAIL"}
                                        handleOnChange={(value: string) =>
                                            setUserForm({
                                                ...user,
                                                email: value,
                                            })
                                        }
                                    />
                                </Grid>
                                <Grid item xs={9}>
                                    <EditableTextField
                                        value={phone || "-"}
                                        label={"PHONE"}
                                        handleOnChange={(value: string) =>
                                            setUserForm({
                                                ...user,
                                                phone: value,
                                            })
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
                                            color="textSecondary"
                                        >
                                            STATUS
                                        </Typography>
                                    </Grid>
                                    <Grid item xs={12}>
                                        <ContainerJustifySpaceBtw
                                            alignItems={"center"}
                                            sx={{ pb: "0px !important" }}
                                        >
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
                                            color="textSecondary"
                                        >
                                            IMEI NUMBER
                                        </Typography>
                                    </Grid>
                                    <Grid item xs={12} mt={1}>
                                        <Typography variant="body2">
                                            {eSimNumber}
                                        </Typography>
                                    </Grid>
                                </Grid>
                                <Grid item container xs={12}>
                                    <Grid item xs={12}>
                                        <Typography
                                            variant="caption"
                                            color="textSecondary"
                                        >
                                            ICCID
                                        </Typography>
                                    </Grid>
                                    <Grid item xs={12} mt={1}>
                                        <Typography variant="body2">
                                            {iccid}
                                        </Typography>
                                    </Grid>
                                </Grid>
                                <Grid item container xs={12}>
                                    <ContainerJustifySpaceBtw
                                        sx={{ alignItems: "center" }}
                                    >
                                        <Typography
                                            variant="caption"
                                            color="textSecondary"
                                            alignSelf={"end"}
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
            )}
        </Dialog>
    );
};

export default UserDetailsDialog;
