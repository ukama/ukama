import {
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
    DialogTitle,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import EditableTextField from "../../EditableTextField";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import { CenterContainer, ContainerJustifySpaceBtw } from "../../../styles";
import { GetUserDto } from "../../../generated";
import LoadingWrapper from "../../LoadingWrapper";
import { ReactEventHandler } from "react";

type BasicDialogProps = {
    type: string;
    user: GetUserDto;
    isOpen: boolean;
    setUserForm: any;
    loading: boolean;
    closeBtnLabel?: string;
    saveBtnLabel?: string;
    userDetailsTitle: string;
    simDetailsTitle: string;
    userStatusLoading: boolean;
    handleSubmitAction: Function;
    handleServiceAction: Function;
    handleClose: ReactEventHandler;
};

const UserDetailsDialog = ({
    type,
    user,
    isOpen,
    setUserForm,
    handleClose,
    saveBtnLabel,
    closeBtnLabel,
    loading = true,
    simDetailsTitle,
    userDetailsTitle,
    userStatusLoading,
    handleSubmitAction,
    handleServiceAction,
}: BasicDialogProps) => {
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

    const statusText = status ? "ACTIVE" : "INACTIVE";
    const title = type === "add" ? "Add User" : "Edit User";
    const statusAction = status ? "PAUSE SERVICE" : "RESUME SERVICE";

    return (
        <Dialog
            key={id}
            open={isOpen}
            onBackdropClick={handleClose}
            maxWidth="sm"
            fullWidth
        >
            {loading ? (
                <CenterContainer>
                    <CircularProgress />
                </CenterContainer>
            ) : (
                <>
                    <Stack
                        direction="row"
                        alignItems="center"
                        justifyContent="space-between"
                    >
                        <DialogTitle>{title}</DialogTitle>
                        <IconButton
                            onClick={handleClose}
                            sx={{ position: "relative", right: 8 }}
                        >
                            <CloseIcon />
                        </IconButton>
                    </Stack>
                    <DialogContent sx={{ overflowX: "hidden" }}>
                        <Grid container spacing={1.5}>
                            <Grid item xs={12}>
                                <Typography variant="subtitle2">
                                    {userDetailsTitle}
                                </Typography>
                                <Divider />
                            </Grid>
                            <Grid item container spacing={1.5}>
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
                                <Grid item xs={12}>
                                    <EditableTextField
                                        value={name}
                                        label={"NAME"}
                                        handleOnChange={(value: string) =>
                                            setUserForm({
                                                ...user,
                                                name: value,
                                            })
                                        }
                                    />
                                </Grid>
                                <Grid item xs={12}>
                                    <EditableTextField
                                        value={email}
                                        label={"EMAIL"}
                                        handleOnChange={(value: string) =>
                                            setUserForm({
                                                ...user,
                                                email: value?.toLowerCase(),
                                            })
                                        }
                                    />
                                </Grid>
                                <Grid item xs={12}>
                                    <EditableTextField
                                        value={phone}
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
                            <Grid item container spacing={1.5}>
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
                                                {statusText}
                                            </Typography>
                                            <LoadingWrapper
                                                height={34}
                                                width={148}
                                                isLoading={userStatusLoading}
                                            >
                                                <Button
                                                    size="small"
                                                    color="error"
                                                    variant="outlined"
                                                    onClick={() => {
                                                        if (id && iccid)
                                                            handleServiceAction(
                                                                id,
                                                                iccid,
                                                                !status
                                                            );
                                                    }}
                                                >
                                                    {statusAction}
                                                </Button>
                                            </LoadingWrapper>
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
                    <DialogActions>
                        <Button
                            onClick={handleClose}
                            sx={{ mr: 2, justifyItems: "center" }}
                        >
                            {closeBtnLabel}
                        </Button>
                        <Button
                            onClick={() => handleSubmitAction()}
                            variant="contained"
                        >
                            {saveBtnLabel}
                        </Button>
                    </DialogActions>
                </>
            )}
        </Dialog>
    );
};

export default UserDetailsDialog;
