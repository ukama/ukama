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
    Alert,
} from "@mui/material";
import React, { useState } from "react";
import { colors } from "../../../theme";
import { makeStyles } from "@mui/styles";
import { IMaskInput } from "react-imask";
import { Node_Type } from "../../../generated";
import { globalUseStyles } from "../../../styles";
import ErrorIcon from "@mui/icons-material/Error";
import CloseIcon from "@mui/icons-material/Close";
import { MASK_BY_TYPE, MASK_PLACEHOLDERS } from "../../../constants";
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

const TextMaskCustom = React.forwardRef<HTMLInputElement, CustomProps>(
    // eslint-disable-next-line no-unused-vars
    function TextMaskCustom(props, ref) {
        const { onChange, ...other } = props;
        return (
            <IMaskInput
                {...other}
                overwrite
                unmask={false}
                mask={MASK_BY_TYPE[props.name]}
                placeholder={MASK_PLACEHOLDERS[props.name]}
                definitions={{
                    "#": /[a-zA-Z0-9]/,
                }}
                onAccept={(value: any) =>
                    onChange({ target: { name: props.name, value } })
                }
            />
        );
    }
);

type NodeDialogProps = {
    nodeData: any;
    isOpen: boolean;
    action?: string;
    subTitle: string;
    handleClose: any;
    subTitle2?: string;
    dialogTitle: string;
    handleNodeSubmitAction: Function;
};

const NodeDialog = ({
    isOpen,
    subTitle,
    nodeData,
    dialogTitle,
    action = "",
    handleClose,
    handleNodeSubmitAction,
}: NodeDialogProps) => {
    const classes = useStyles();
    const gclasses = globalUseStyles();
    const [formData, setFormData] = useState({
        type: nodeData.type,
        name: nodeData.name,
        nodeId: nodeData.nodeId,
        orgId: nodeData.orgId,
    });
    const [error, setError] = useState("");

    const handleRegisterNode = () => {
        if (!formData.name || !formData.nodeId) {
            setError("Please fill all require vields");
            return;
        }

        handleNodeSubmitAction(formData);
    };

    const handleNodeTypeChange = (e: SelectChangeEvent) =>
        setFormData({ ...formData, nodeId: "", type: e.target.value });

    return (
        <Dialog open={isOpen} onClose={handleClose} maxWidth="sm" fullWidth>
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
            {error && (
                <Alert
                    sx={{
                        mx: 3,
                        mb: 1,
                        color: theme => theme.palette.text.primary,
                    }}
                    severity={"error"}
                    icon={<ErrorIcon sx={{ color: colors.red }} />}
                >
                    {error}
                </Alert>
            )}
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
                <Grid container rowSpacing={2.75} columnSpacing={2.75} mt={2}>
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
                                value={formData.type}
                                variant="outlined"
                                onChange={handleNodeTypeChange}
                                disabled={action == "editNode"}
                                input={
                                    <OutlinedInput
                                        notched
                                        label="NODE TYPE"
                                        name="node_type"
                                        id="outlined-age-always-notched"
                                    />
                                }
                                MenuProps={{
                                    disablePortal: false,
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
                            required
                            fullWidth
                            value={formData.name}
                            label={"NODE NAME"}
                            InputLabelProps={{ shrink: true }}
                            InputProps={{
                                classes: {
                                    input: gclasses.inputFieldStyle,
                                },
                            }}
                            onChange={(e: any) =>
                                setFormData({
                                    ...formData,
                                    name: e.target.value,
                                })
                            }
                        />
                    </Grid>
                    <Grid item xs={12}>
                        <TextField
                            fullWidth
                            required
                            value={formData.nodeId}
                            label={"NODE NUMBER"}
                            onChange={(e: any) =>
                                setFormData({
                                    ...formData,
                                    nodeId: e.target.value.replace(/ /g, ""),
                                })
                            }
                            disabled={action == "editNode"}
                            InputLabelProps={{ shrink: true }}
                            name={formData.type}
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
                        <Grid item xs={12}>
                            <TextField
                                fullWidth
                                value={formData.orgId}
                                disabled={true}
                                label={"ORGANIZATION ID"}
                                InputLabelProps={{ shrink: true }}
                                InputProps={{
                                    classes: {
                                        input: gclasses.inputFieldStyle,
                                    },
                                }}
                            />
                        </Grid>
                    )}
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button sx={{ mr: 2 }} onClick={handleClose}>
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

export default NodeDialog;
