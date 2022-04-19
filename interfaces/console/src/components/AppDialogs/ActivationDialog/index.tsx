import {
    Grid,
    Stack,
    Button,
    Dialog,
    Select,
    MenuItem,
    TextField,
    IconButton,
    InputLabel,
    Typography,
    DialogTitle,
    FormControl,
    DialogActions,
    DialogContent,
    OutlinedInput,
    DialogContentText,
} from "@mui/material";
import { colors } from "../../../theme";
import { makeStyles } from "@mui/styles";
import { IMaskInput } from "react-imask";
import { Node_Type } from "../../../generated";
import { globalUseStyles } from "../../../styles";
import { MASK_BY_TYPE } from "../../../constants";
import CloseIcon from "@mui/icons-material/Close";
import React, { useState, useEffect } from "react";
import { SelectChangeEvent } from "@mui/material/Select/SelectInput";

const useStyles = makeStyles(() => ({
    basicDialogHeaderStyle: {
        padding: "0px 0px 18px 0px",
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
    actionContainer: {
        padding: "0px",
        marginTop: "16px",
        justifyContent: "space-between",
    },
    "&.MuiFormHelperText-root.Mui-error": {
        color: "red",
    },
    stepButtonStyle: {
        "&:disabled": {
            color: colors.white,
            backgroundColor: colors.nightGrey,
        },
    },
    selectStyle: () => ({
        width: "100%",
        height: "48px",
    }),
    formControl: {
        width: "100%",
        height: "48px",
    },
}));

interface CustomProps {
    // eslint-disable-next-line no-unused-vars
    onChange: (event: { target: { name: string; value: string } }) => void;
    name: Node_Type;
}

const TextMaskCustom = React.forwardRef<HTMLElement, CustomProps>(
    function TextMaskCustom(props) {
        const { onChange, ...other } = props;
        return (
            <IMaskInput
                {...other}
                mask={MASK_BY_TYPE[props.name]}
                unmask={false}
                lazy={false}
                overwrite
                onAccept={(value: any) =>
                    onChange({ target: { name: props.name, value } })
                }
            />
        );
    }
);

type ActivationDialogProps = {
    isOpen: boolean;
    subTitle: string;
    handleClose: any;
    subTitle2?: string;
    dialogTitle: string;
    nodeData?: any;
    handleActivationSubmit: Function;
    action?: string;
};

const ActivationDialog = ({
    isOpen,
    subTitle,
    nodeData,
    dialogTitle,
    action = "",
    handleClose,
    handleActivationSubmit,
}: ActivationDialogProps) => {
    const classes = useStyles();
    const gclasses = globalUseStyles();
    const [nodeType, setNodeType] = useState("HOME");
    const [nodeName, setNodeName] = useState("");
    const [nodeSerial, setNodeSerial] = useState("");
    const [nodeNameError, setNodeNameError] = useState("");
    const [nodeSerialError, setNodeSerialError] = useState("");
    const [orgIdError, setOrgIdError] = useState("");
    const [orgId, setOrgId] = useState("");
    useEffect(() => {
        if (action == "editNode" && nodeData) {
            setNodeName(nodeData.name);
            setNodeSerial(nodeData.id);
            setOrgId(nodeData.orgId);
        }
    }, [nodeData]);
    const handleRegisterNode = () => {
        if (action == "editNode" && nodeName && nodeSerial && orgId) {
            handleActivationSubmit({
                name: nodeName,
                nodeId: nodeSerial,
                orgId: orgId,
            });
        } else {
            handleActivationSubmit({
                name: nodeName,
                nodeId: nodeSerial,
            });
        }

        if (!nodeName) {
            setNodeNameError("Node Name is required!");
        }
        if (!nodeSerial) {
            setNodeSerialError("Node number is required!");
        }
        if (!orgId) {
            setOrgIdError("Organiation Id is required!");
        }
    };

    useEffect(() => {
        if (nodeName.length > 0) {
            setNodeNameError("");
        }
    }, [nodeName]);

    useEffect(() => {
        if (nodeSerial.length > 0) {
            setNodeSerialError("");
        }
    }, [nodeSerial]);

    useEffect(() => {
        if (orgId.length > 0) {
            setOrgIdError("");
        }
    }, [orgId]);

    const handleNodeTypeChange = (e: SelectChangeEvent) => {
        setNodeSerial("");
        setNodeType(e.target.value);
    };

    return (
        <Dialog open={isOpen} onClose={handleClose}>
            <Stack
                direction="row"
                alignItems="center"
                justifyContent="space-between"
            >
                <DialogTitle>{dialogTitle}</DialogTitle>
                <IconButton
                    onClick={handleClose}
                    sx={{ position: "relative", right: 8 }}
                >
                    <CloseIcon />
                </IconButton>
            </Stack>

            <DialogContent>
                <DialogContentText>
                    <Typography
                        component={"span"}
                        variant="body1"
                        color={"textPrimary"}
                    >
                        {subTitle}
                    </Typography>
                </DialogContentText>
                <Grid container spacing={2} mt={2}>
                    <Grid item xs={12} md={6}>
                        <FormControl
                            variant="outlined"
                            className={classes.formControl}
                        >
                            <InputLabel
                                shrink
                                variant="outlined"
                                htmlFor="outlined-age-always-notched"
                            >
                                NODE TYPE
                            </InputLabel>
                            <Select
                                value={nodeType}
                                variant="outlined"
                                onChange={handleNodeTypeChange}
                                input={
                                    <OutlinedInput
                                        notched
                                        label="NODE TYPE"
                                        name="node_type"
                                        id="outlined-age-always-notched"
                                    />
                                }
                                MenuProps={{
                                    disablePortal: true,
                                    PaperProps: {
                                        sx: {
                                            boxShadow:
                                                "0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)",
                                            borderRadius: "4px",
                                        },
                                    },
                                }}
                                className={classes.selectStyle}
                            >
                                {[
                                    { id: 1, label: "Home", value: "HOME" },
                                    {
                                        id: 2,
                                        label: "Amplifier",
                                        value: "AMPLIFIER",
                                    },
                                    { id: 3, label: "Tower", value: "TOWER" },
                                ].map(({ id, label, value }) => (
                                    <MenuItem
                                        key={id}
                                        value={value}
                                        sx={{
                                            m: 0,
                                            p: "6px 16px",
                                        }}
                                    >
                                        <Typography variant="body1">
                                            {label}
                                        </Typography>
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <TextField
                            error={nodeNameError ? true : false}
                            fullWidth
                            value={nodeName}
                            label={"NODE NAME"}
                            InputLabelProps={{ shrink: true }}
                            helperText={nodeNameError}
                            InputProps={{
                                classes: {
                                    input: gclasses.inputFieldStyle,
                                },
                            }}
                            onChange={(e: any) => setNodeName(e.target.value)}
                        />
                    </Grid>
                    <Grid item xs={12}>
                        <TextField
                            fullWidth
                            value={nodeSerial}
                            error={nodeSerialError ? true : false}
                            label={"NODE NUMBER"}
                            helperText={nodeSerialError}
                            onChange={(e: any) => {
                                setNodeSerial(e.target.value.replace(/ /g, ""));
                            }}
                            InputLabelProps={{ shrink: true }}
                            name={nodeType}
                            id="formatted-text-mask-input"
                            spellCheck={false}
                            InputProps={{
                                inputComponent: TextMaskCustom as any,
                                classes: {
                                    input: gclasses.inputFieldStyle,
                                },
                            }}
                        />
                    </Grid>
                    {action == "editNode" && (
                        <Grid item xs={12} md={12}>
                            <TextField
                                error={orgIdError ? true : false}
                                fullWidth
                                value={orgId}
                                label={"ORGANIZATION ID"}
                                InputLabelProps={{ shrink: true }}
                                helperText={orgIdError}
                                InputProps={{
                                    classes: {
                                        input: gclasses.inputFieldStyle,
                                    },
                                }}
                                onChange={(e: any) => setOrgId(e.target.value)}
                            />
                        </Grid>
                    )}
                </Grid>
            </DialogContent>
            <DialogActions sx={{ mr: 2, paddingBottom: 3 }}>
                <Button
                    sx={{ color: colors.primaryMain, mr: 2 }}
                    onClick={handleClose}
                >
                    Cancel
                </Button>
                <Button
                    variant="contained"
                    type="submit"
                    onClick={handleRegisterNode}
                    className={classes.stepButtonStyle}
                >
                    {action == "editNode" ? "UPDATE NODE" : "REGISTER NODE"}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default ActivationDialog;
