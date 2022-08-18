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
import React, { useState } from "react";
import { colors } from "../../../theme";
import { IMaskInput } from "react-imask";
import { makeStyles } from "@mui/styles";
import { Node_Type } from "../../../generated";
import { globalUseStyles } from "../../../styles";
import CloseIcon from "@mui/icons-material/Close";
import { MASK_BY_TYPE, MASK_PLACEHOLDERS, NODE_TYPE } from "../../../constants";
import { SelectChangeEvent } from "@mui/material/Select/SelectInput";
import AddNodeForm from "./addNodeForm";
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

    const [isSubmitted, setIsSubmitted] = useState<boolean>(false);

    const getNodeArray = (data: any) => {
        setFormData({
            ...formData,
            associatedTowerNode: data,
        });
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

    const handleRegisterNode = () => {
        setIsSubmitted(true);
        if (!formData.name || !formData.nodeId) {
            return;
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
                <Grid container columnSpacing={2} mt={2}>
                    <Grid item xs={12} md={4}>
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
                    <Grid item xs={12} md={8}>
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

                    {(formData.type == "TOWER" ||
                        formData.type == "AMPLIFIER") && (
                        <AddNodeForm
                            nodeType={formData.type}
                            nodeArray={getNodeArray}
                            towerNodesArrayList={towerNodesArrayList}
                            handleAssociatedTowerNode={
                                handleAssociatedTowerNode
                            }
                            isSubmitted={isSubmitted}
                            selectedToweNode={selectedToweNode}
                        />
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
