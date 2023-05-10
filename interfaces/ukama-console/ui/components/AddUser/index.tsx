import {
    Stack,
    Dialog,
    IconButton,
    DialogTitle,
    DialogContent,
} from "@mui/material";
import ESimQR from "./ESimQR";
import Userform from "./Userform";
import CloseIcon from "@mui/icons-material/Close";
import { useState, useEffect } from "react";
interface IAddUser {
    isOpen: boolean;
    isPsimAdded: boolean;
    handleClose: Function;
    loading: boolean;
    qrCodeId: any;
    addedUserName: any;
    iSeSimAdded: boolean;
    handleEsimInstallation: Function;
    handlePhysicalSimInstallationFlow1: Function;
    handlePhysicalSimInstallationFlow2: Function;
    step: number;
}

const getDescription = (id: number, addUserName?: any) => {
    switch (id) {
        case 0:
            return "What SIM do you want to assign to this user?";
        case 1:
            return "Start accessing high quality and fast data now. You’ll be able to add more users to the network later.";
        case 2:
            return `You have successfully added ${addUserName} as a user to your network, and an eSIM installation invitation has been sent out to them. If they would rather install now, have them scan the QR code below.`;
        case 3:
            return "Enter security code for Physical SIM lorem ipsum. Instructions for remembering to install SIM after?";
        case 4:
            return `You have successfully added ${addUserName} as a user to your network. Instructions for installing physical SIM (might need more thinking if this process is complex).`;
        default:
            return "";
    }
};

const AddUser = ({
    isOpen,
    qrCodeId,
    handleClose,
    iSeSimAdded,
    isPsimAdded,
    addedUserName,
    loading = false,
    handleEsimInstallation,
}: IAddUser) => {
    const [step, setstep] = useState(1);
    const getTitle = (esimSuccess: boolean) => {
        if (esimSuccess || isPsimAdded) {
            return "Add User Succesful";
        } else {
            return `Add User`;
        }
    };
    const getSimType = (sim: any) => {
        if (!sim) return;
        // setSelectedSimType(sim);
    };
    useEffect(() => {
        if (iSeSimAdded) {
            setstep(step + 1);
        }
    }, [iSeSimAdded]);
    const getAddUserComponent = (step: number) => {
        switch (step) {
            case 1:
                return (
                    <Userform
                        getSimType={getSimType}
                        loading={loading}
                        handleClose={handleClose}
                        description={getDescription(1)}
                        handleSimInstallation={handleEsimInstallation}
                    />
                );
            case 2:
                return (
                    <ESimQR
                        description={getDescription(2, addedUserName)}
                        qrCodeId={qrCodeId}
                        handleClose={handleClose}
                    />
                );

            default:
                return (
                    <Userform
                        getSimType={getSimType}
                        loading={loading}
                        handleClose={handleClose}
                        description={getDescription(1)}
                        handleSimInstallation={handleEsimInstallation}
                    />
                );
        }
    };
    return (
        <Dialog
            fullWidth
            open={isOpen}
            maxWidth="sm"
            onClose={() => handleClose()}
            onBackdropClick={() => handleClose()}
        >
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                <DialogTitle>{!loading && getTitle(iSeSimAdded)}</DialogTitle>
                <IconButton
                    onClick={() => handleClose()}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>
            <DialogContent sx={{ overflowX: "hidden" }}>
                {getAddUserComponent(step)}
            </DialogContent>
        </Dialog>
    );
};

export default AddUser;
