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
} from "@mui/material";
import React, { Fragment, useState } from "react";
import { colors } from "../../../theme";
import { IMaskInput } from "react-imask";
import { makeStyles } from "@mui/styles";
import AddIcon from "@mui/icons-material/Add";
import { Node_Type } from "../../../generated";
import { globalUseStyles } from "../../../styles";
import CloseIcon from "@mui/icons-material/Close";
import { MASK_BY_TYPE, MASK_PLACEHOLDERS, NODE_TYPE } from "../../../constants";
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
    name: Node_Type;
}

const TextMaskCustom = React.forwardRef<HTMLInputElement, CustomProps>(
    function TextMaskCustom(props, _ref) {
        const { ...other } = props;
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
    towerNodesArrayList?: any;
    handleNodeSubmitAction: Function;
};

const NodeDialog = ({
    isOpen,
    subTitle,
    nodeData,
    dialogTitle,
    action = "",
    handleClose,
    towerNodesArrayList,

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
        isAssiociatedTowerNode: nodeData.isAssiociatedTowerNode,
    });
    const [attachedAmplierNode, setAttachedAmplierNode] = useState<any>([
        {
            nodeId: "",
            name: "",
        },
    ]);
    const handleInputChange = (e: any, index: number) => {
        const { id, value } = e.target;
        const list: any = [...attachedAmplierNode];
        list[index][id] = id == "nodeId" ? value.replace(/ /g, "") : value;
        setAttachedAmplierNode(list);
        setFormData({
            ...formData,
            associatedTowerNode: attachedAmplierNode,
        });
    };
    const [isSubmitted, setIsSubmitted] = useState<boolean>(false);
    const handleRemoveClick = (index: number) => {
        const list = [...attachedAmplierNode];
        list.splice(index, 1);
        setAttachedAmplierNode(list);
    };
    const onAddTowerNode = () => {
        setIsAssociatedTowerNode(true);
    };

    const handleAddClick = () => {
        setAttachedAmplierNode([
            ...attachedAmplierNode,
            { nodeId: "", name: "" },
        ]);
    };

    const handleNodeTypeChange = (e: SelectChangeEvent) => {
        setFormData({ ...formData, nodeId: "", type: e.target.value });
    };
    const [selectedToweNode, setSelectedToweNode] = useState("");
    const handleAssociatedTowerNode = (e: SelectChangeEvent) => {
        setSelectedToweNode(e.target.value);
        const getSelectedNodeInfo = towerNodesArrayList.filter(
            (item: { name: string }) => item.name === e.target.value
        );
        const result = getSelectedNodeInfo.map(({ name, id }: any) => ({
            name,
            id,
        }))[0];

        setFormData({
            ...formData,
            isAssiociatedTowerNode: true,
            associatedTowerNode: result,
        });
    };
    const [nType, setNtype] = useState<any>("AMPLIFIER");
    const [isAssociatedTowerNode, setIsAssociatedTowerNode] =
        useState<boolean>(false);
    const handleOptionalNodeType = (e: SelectChangeEvent) => {
        setNtype(e.target.value);
    };
    const removeNodeTypefromArray = (from: number, to: number) => {
        return NODE_TYPE.filter(function (value, index) {
            return [from, to].indexOf(index) == -1;
        });
    };
    const handleRegisterNode = () => {
        setIsSubmitted(true);
        if (!formData.name || !formData.nodeId) {
            return;
        }
        if (isAssociatedTowerNode) {
            if (!formData.associatedTowerNode) {
                return;
            }
        }
        handleNodeSubmitAction(formData);
    };

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
                                {NODE_TYPE.map(({ id, label, value }) => (
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
                            error={isSubmitted && formData.nodeId === ""}
                            helperText={
                                isSubmitted && formData.nodeId === ""
                                    ? "Node number is required !"
                                    : " "
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
                            error={isSubmitted && formData.name === ""}
                            helperText={
                                isSubmitted && formData.name === ""
                                    ? "Node name is required !"
                                    : " "
                            }
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
                                    value={selectedToweNode}
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
                                    {towerNodesArrayList.map(
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

                    {attachedAmplierNode.map((x: any, i: number) => {
                        return (
                            (formData.type == "TOWER" ||
                                isAssociatedTowerNode == true) && (
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
                                                {formData.type == "TOWER" &&
                                                    "AMPLIFIER NODE"}

                                                {formData.type == "AMPLIFIER" &&
                                                    "NEW TOWER NODE"}
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
                                                    variant="outlined"
                                                    onChange={
                                                        handleOptionalNodeType
                                                    }
                                                    value={nType}
                                                    disabled={
                                                        action == "editNode"
                                                    }
                                                    input={
                                                        <OutlinedInput
                                                            notched
                                                            label="NODE TYPE"
                                                            name={x.type}
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
                                                    {removeNodeTypefromArray(
                                                        0,
                                                        isAssociatedTowerNode
                                                            ? 1
                                                            : 2
                                                    ).map(
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
                                                value={x.nodeId}
                                                onChange={(e: any) => {
                                                    handleInputChange(e, i);
                                                }}
                                                InputLabelProps={{
                                                    shrink: true,
                                                }}
                                                error={
                                                    isSubmitted &&
                                                    x.nodeId === ""
                                                }
                                                helperText={
                                                    isSubmitted &&
                                                    x.nodeId === ""
                                                        ? "Node number is required !"
                                                        : " "
                                                }
                                                id={"nodeId"}
                                                name={"AMPLIFIER"}
                                                spellCheck={false}
                                                InputProps={{
                                                    inputComponent:
                                                        TextMaskCustom as any,
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
                                                name={"name"}
                                                value={x.name}
                                                id={"name"}
                                                onChange={e => {
                                                    handleInputChange(e, i);
                                                }}
                                                error={
                                                    isSubmitted && x.name === ""
                                                }
                                                helperText={
                                                    isSubmitted && x.name === ""
                                                        ? "Node name is required !"
                                                        : " "
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
                                            attachedAmplierNode.length < 2 &&
                                            !isAssociatedTowerNode && (
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
