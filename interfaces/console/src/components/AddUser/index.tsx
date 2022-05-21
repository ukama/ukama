import {
    Stack,
    Dialog,
    IconButton,
    DialogTitle,
    CircularProgress,
    DialogContent,
} from "@mui/material";
import ESimQR from "./ESimQR";
import Success from "./Success";
import { useState } from "react";
import Userform from "./Userform";
import ChooseSim from "./ChooseSim";
import PhysicalSimform from "./PhysicalSimform";
import CloseIcon from "@mui/icons-material/Close";
import { CenterContainer } from "../../styles";

interface IAddUser {
    isOpen: boolean;
    handleClose: Function;
    loading?: boolean;
    qrCodeId: any;
    handlePhysicalSimInstallation: Function;
    addedUserName: any;
    iSeSimAdded: boolean;
    handleEsimInstallation: Function;
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

const getTitle = (iSeSimAdded: boolean, type: any) =>
    !iSeSimAdded ? `Add User${type && ` - ${type}`}` : "Add User Succesful";

const AddUser = ({
    isOpen,
    handleClose,
    loading = false,
    qrCodeId,
    handlePhysicalSimInstallation,
    addedUserName,
    handleEsimInstallation,
    iSeSimAdded,
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
            case 2:
                if (type !== "eSIM") {
                    setFlow(4);
                }
                break;
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
                    {!loading ? getTitle(iSeSimAdded, simType) : null}
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
                {simType == "eSIM" && (
                    <>
                        {loading ? (
                            <CenterContainer>
                                <CircularProgress />
                            </CenterContainer>
                        ) : (
                            <Userform
                                handleEsimInstallation={handleEsimInstallation}
                                description={getDescription(1)}
                            />
                        )}
                    </>
                )}
                {simType == "Physical SIM" && (
                    <>
                        {loading ? (
                            <CenterContainer>
                                <CircularProgress />
                            </CenterContainer>
                        ) : (
                            <PhysicalSimform
                                handlePhysicalSimInstallation={
                                    handlePhysicalSimInstallation
                                }
                                description={getDescription(3)}
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
                {flow === 4 && (
                    <Success description={getDescription(4, addedUserName)} />
                )}
            </DialogContent>
        </Dialog>
    );
};

export default AddUser;
