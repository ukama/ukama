import {
    Stack,
    Button,
    Dialog,
    IconButton,
    DialogTitle,
    DialogActions,
    DialogContent,
} from "@mui/material";
import ESimQR from "./ESimQR";
import Success from "./Success";
import { useState } from "react";
import Userform from "./Userform";
import ChooseSim from "./ChooseSim";
import PhysicalSimform from "./PhysicalSimform";
import CloseIcon from "@mui/icons-material/Close";
import { isEmailValid } from "../../utils";

interface IAddUser {
    isOpen: boolean;
    handleClose: Function;
    handleSubmitAction: Function;
}

const getDescription = (id: number) => {
    switch (id) {
        case 0:
            return "What SIM do you want to assign to this user?";
        case 1:
            return "Add user xyz. They will be emailed the SIM installation link/QR code shortly after.";
        case 2:
            return "You have successfully added [Name] as a user to your network, and an eSIM installation invitation has been sent out to them. If they would rather install now, have them scan the QR code below.";
        case 3:
            return "Enter security code for Physical SIM lorem ipsum. Instructions for remembering to install SIM after?";
        case 4:
            return "You have successfully added [Name] as a user to your network. Instructions for installing physical SIM (might need more thinking if this process is complex).";
        default:
            return "";
    }
};

const getTitle = (id: number, type: string) =>
    id === 0 || id === 1 || id === 2
        ? `Add User${type ? ` - ${type}` : ""}`
        : "Add User Succesful";

const AddUser = ({ isOpen, handleClose, handleSubmitAction }: IAddUser) => {
    const [flow, setFlow] = useState(0);
    const [formError, setError] = useState("");
    const [simType, setSimType] = useState("");
    const [form, setForm] = useState({
        name: "",
        code: "",
        email: "",
        iccid: "",
        roaming: false,
    });

    const handleAction = ({ type = simType }: { type?: string }) => {
        setError("");
        switch (flow) {
            case 0:
                setSimType(type);
                setFlow(flow + 1);
                break;
            case 1:
                if (!form.email || !form.name) {
                    setError("Please file require fileds.");
                    return;
                }

                if (!isEmailValid(form.email)) {
                    setError("Please enter valid email!");
                    return;
                }

                if (type === "eSIM") {
                    setFlow(3);
                    handleSubmitAction(form);
                } else setFlow(flow + 1);
                break;
            case 2:
                if (!form.code || !form.iccid) {
                    setError("Please file require fileds.");
                    return;
                }

                if (type !== "eSIM") {
                    setFlow(4);
                    handleSubmitAction(form);
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
                <DialogTitle>{getTitle(flow, simType)}</DialogTitle>
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
                {flow === 1 && (
                    <Userform
                        formData={form}
                        formError={formError}
                        setFormData={setForm}
                        description={getDescription(1)}
                    />
                )}
                {flow === 2 && (
                    <PhysicalSimform
                        formData={form}
                        formError={formError}
                        setFormData={setForm}
                        description={getDescription(3)}
                    />
                )}
                {flow === 3 && <ESimQR description={getDescription(2)} />}
                {flow === 4 && <Success description={getDescription(4)} />}
            </DialogContent>
            {(flow === 1 || flow === 2) && (
                <DialogActions>
                    <Button
                        onClick={() => handleClose()}
                        sx={{ mr: 2, justifyItems: "center" }}
                    >
                        Cancel
                    </Button>
                    <Button
                        variant="contained"
                        onClick={() => handleAction({})}
                    >
                        {simType !== "eSIM" && flow === 1
                            ? "Continue"
                            : "ADD USER"}
                    </Button>
                </DialogActions>
            )}
        </Dialog>
    );
};

export default AddUser;
