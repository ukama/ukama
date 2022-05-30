import {
    Stack,
    Dialog,
    IconButton,
    DialogTitle,
    CircularProgress,
    DialogContent,
} from "@mui/material";
import ESimQR from "./ESimQR";
import { useState } from "react";
import Userform from "./Userform";
import ChooseSim from "./ChooseSim";
import CloseIcon from "@mui/icons-material/Close";
import { CenterContainer } from "../../styles";
import PhysicalSimFlow from "./PhysicalSimFlow";
interface IAddUser {
    isOpen: boolean;
    isPsimAdded: boolean;
    handleClose: Function;
    loading?: boolean;
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
            return "Add user xyz. They will be emailed the SIM installation link/QR code shortly after.";
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
    step,
    handlePhysicalSimInstallationFlow2,
    handlePhysicalSimInstallationFlow1,
}: IAddUser) => {
    const [flow, setFlow] = useState(0);
    const [simType, setSimType] = useState("");

    const handleAction = ({ type = simType }: { type?: string }) => {
        switch (flow) {
            case 0:
                setSimType(type);
                setFlow(flow + 1);
                break;
            case 1:
                setFlow(2);
                break;
        }
    };
    const getTitle = (iSeSimAdded: boolean, type: any) => {
        if (iSeSimAdded || isPsimAdded) {
            return "Add User Succesful";
        } else {
            return `Add User${type && ` - ${type}`}`;
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
                <DialogTitle>
                    {!loading && getTitle(iSeSimAdded, simType)}
                </DialogTitle>
                <IconButton
                    onClick={() => handleClose()}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>
            <DialogContent sx={{ overflowX: "hidden" }}>
                {flow === 0 && (
                    <ChooseSim
                        description={getDescription(0)}
                        handleSimType={handleAction}
                    />
                )}
                {simType == "eSIM" && !iSeSimAdded && (
                    <>
                        {loading ? (
                            <CenterContainer>
                                <CircularProgress />
                            </CenterContainer>
                        ) : (
                            <Userform
                                handleClose={handleClose}
                                description={getDescription(1)}
                                handleEsimInstallation={handleEsimInstallation}
                            />
                        )}
                    </>
                )}
                {simType == "Physical SIM" && !iSeSimAdded && (
                    <>
                        {loading ? (
                            <CenterContainer>
                                <CircularProgress />
                            </CenterContainer>
                        ) : (
                            <PhysicalSimFlow
                                handleClose={handleClose}
                                step={step}
                                handlePhysicalSimInstallationFlow1={
                                    handlePhysicalSimInstallationFlow1
                                }
                                handlePhysicalSimInstallationFlow2={
                                    handlePhysicalSimInstallationFlow2
                                }
                            />
                        )}
                    </>
                )}

                {iSeSimAdded && (
                    <ESimQR
                        description={getDescription(2, addedUserName)}
                        qrCodeId={qrCodeId}
                    />
                )}
            </DialogContent>
        </Dialog>
    );
};

export default AddUser;
