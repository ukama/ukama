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
    Divider,
    DialogContentText,
    Alert,
} from "@mui/material";
import React, { Fragment, useState } from "react";
import { colors } from "../../../theme";
import { makeStyles } from "@mui/styles";
import { IMaskInput } from "react-imask";
import AddIcon from "@mui/icons-material/Add";
import { Node_Type } from "../../../generated";
import { globalUseStyles } from "../../../styles";
import ErrorIcon from "@mui/icons-material/Error";
import CloseIcon from "@mui/icons-material/Close";
import { MASK_BY_TYPE, MASK_PLACEHOLDERS } from "../../../constants";
import RemoveCircleOutlineIcon from "@mui/icons-material/RemoveCircleOutline";
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
    function TextMaskCustom(props, _ref) {
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
    associatedTowerNodes?: any;
    handleNodeSubmitAction: Function;
};

const NodeDialog = ({
    isOpen,
    subTitle,
    nodeData,
    dialogTitle,
    action = "",
    handleClose,
    associatedTowerNodes,
    handleNodeSubmitAction,
}: NodeDialogProps) => {
    const classes = useStyles();
    const gclasses = globalUseStyles();
    const [formData, setFormData] = useState({
        type: nodeData.type,
        name: nodeData.name,
        nodeId: nodeData.nodeId,
        orgId: nodeData.orgId,
        associatedTowerNode: nodeData.associatedTowerNode,
    });

    const [error, setError] = useState("");
    const [attachedAmplierNode, setAttachedAmplierNode] = useState([
        {
            nodeId: "",
            nodeName: "",
        },
    ]);
    const handleInputChange = (e: any, index: number) => {
        const { name, value } = e.target;
        const list: any = [...attachedAmplierNode];
        list[index][name] = value;
        setAttachedAmplierNode(list);
    };

    const handleRemoveClick = (index: number) => {
        const list = [...attachedAmplierNode];
        list.splice(index, 1);
        setAttachedAmplierNode(list);
    };
    const onAddTowerNode = () => {
        //handle addTowerNode
    };

    const handleAddClick = () => {
        setAttachedAmplierNode([
            ...attachedAmplierNode,
            { nodeId: "", nodeName: "" },
        ]);
    };

    const handleRegisterNode = () => {
        if (!formData.name || !formData.nodeId) {
            setError("Please fill all require vields");
            return;
        }

        handleNodeSubmitAction(formData);
    };

    const handleNodeTypeChange = (e: SelectChangeEvent) => {
        setFormData({ ...formData, nodeId: "", type: e.target.value });
    };

    const handleAssociatedTowerNode = (e: SelectChangeEvent) => {
        setFormData({
            ...formData,
            nodeId: "",
            associatedTowerNode: e.target.value as string,
        });
    };

    console.log(attachedAmplierNode.length);
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
                            fullWidth
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
                    <Grid item xs={12}>
                        <TextField
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
                    {formData.type == "AMPLIFIER" && (
                        <Grid item xs={12}>
                            <FormControl
                                variant="outlined"
                                className={classes.formControl}
                            >
                                <InputLabel
                                    shrink
                                    variant="outlined"
                                    htmlFor="associatedTowerNode"
                                >
                                    ASSOCIATED TOWER NODE
                                </InputLabel>
                                <Select
                                    labelId="associatedTowerNodel"
                                    id="associatedTowerNode"
                                    sx={{
                                        "& legend": { width: "190px" },
                                    }}
                                    onChange={handleAssociatedTowerNode}
                                    value={formData.associatedTowerNode}
                                    variant="outlined"
                                    disabled={action == "editNode"}
                                    input={
                                        <OutlinedInput
                                            notched
                                            label="ASSOCIATED TOWER NODE"
                                            name="associatedTowerNode"
                                            id="associatedTowerNode"
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
                                    {associatedTowerNodes.map(
                                        ({ id, name }: any) => (
                                            <MenuItem value={name} key={id}>
                                                <Typography variant="body1">
                                                    {name}
                                                </Typography>
                                            </MenuItem>
                                        )
                                    )}
                                    <Divider />
                                    <MenuItem
                                        onClick={e => {
                                            onAddTowerNode();
                                            e.stopPropagation();
                                        }}
                                    >
                                        <Button
                                            variant="text"
                                            sx={{
                                                typography: "body1",
                                                textTransform: "none",
                                            }}
                                        >
                                            Add Tower Node
                                        </Button>
                                    </MenuItem>
                                </Select>
                            </FormControl>
                        </Grid>
                    )}
                    {attachedAmplierNode.map((x: any, i: any) => {
                        return (
                            formData.type == "TOWER" && (
                                <Fragment key={i}>
                                    <Grid item xs={12}>
                                        <Divider />
                                    </Grid>
                                    <Grid
                                        item
                                        xs={12}
                                        sx={{
                                            pt: 4,
                                        }}
                                    >
                                        <Stack
                                            direction="row"
                                            justifyContent="space-between"
                                        >
                                            <Typography
                                                variant="body2"
                                                sx={{
                                                    fontStyle: "Bold",
                                                }}
                                            >
                                                AMPLIFIER NODE
                                            </Typography>
                                            <IconButton
                                                color="primary"
                                                aria-label="remove-node"
                                                component="span"
                                                onClick={() =>
                                                    handleRemoveClick(i)
                                                }
                                            >
                                                <RemoveCircleOutlineIcon />
                                            </IconButton>
                                        </Stack>
                                    </Grid>

                                    <Fragment key={i}>
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
                                                    value={"AMPLIFIER"}
                                                    variant="outlined"
                                                    onChange={
                                                        handleNodeTypeChange
                                                    }
                                                    disabled={
                                                        action == "editNode"
                                                    }
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
                                                                borderRadius:
                                                                    "4px",
                                                            },
                                                        },
                                                    }}
                                                    className={
                                                        classes.selectStyle
                                                    }
                                                >
                                                    {[
                                                        {
                                                            id: 2,
                                                            label: "Amplifier",
                                                            value: "AMPLIFIER",
                                                        },
                                                    ].map(
                                                        ({
                                                            id,
                                                            label,
                                                            value,
                                                        }) => (
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
                                                        )
                                                    )}
                                                </Select>
                                            </FormControl>
                                        </Grid>
                                        <Grid item xs={12} md={6}>
                                            <TextField
                                                fullWidth
                                                label={"NODE NUMBER"}
                                                name="nodeId"
                                                value={x.nodeId}
                                                onChange={e =>
                                                    handleInputChange(e, i)
                                                }
                                                disabled={action == "editNode"}
                                                InputLabelProps={{
                                                    shrink: true,
                                                }}
                                                InputProps={{
                                                    classes: {
                                                        input: gclasses.inputFieldStyle,
                                                    },
                                                }}
                                            />
                                        </Grid>
                                        <Grid item xs={12}>
                                            <TextField
                                                fullWidth
                                                label={"NODE NAME"}
                                                name={"nodeName"}
                                                value={x.nodeName}
                                                onChange={e =>
                                                    handleInputChange(e, i)
                                                }
                                                disabled={action == "editNode"}
                                                InputLabelProps={{
                                                    shrink: true,
                                                }}
                                                InputProps={{
                                                    classes: {
                                                        input: gclasses.inputFieldStyle,
                                                    },
                                                }}
                                            />
                                        </Grid>
                                    </Fragment>

                                    <Grid item xs={12} sx={{ pt: 2 }}>
                                        {attachedAmplierNode.length - 1 === i &&
                                            attachedAmplierNode.length < 2 && (
                                                <Button
                                                    variant="text"
                                                    startIcon={<AddIcon />}
                                                    onClick={handleAddClick}
                                                    sx={{
                                                        color: colors.primaryMain,
                                                        pointer: "cursor",
                                                    }}
                                                >
                                                    add amplifier node
                                                </Button>
                                            )}
                                    </Grid>
                                </Fragment>
                            )
                        );
                    })}

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
