import { FieldArray, Form, Formik, getIn } from "formik";
import { makeStyles } from "@mui/styles";
import * as Yup from "yup";
import CloseIcon from "@mui/icons-material/Close";
import { colors } from "../../../theme";
import React, { useState } from "react";
import { MASK_BY_TYPE, MASK_PLACEHOLDERS, NODE_TYPE } from "../../../constants";
import { globalUseStyles } from "../../../styles";
import RemoveCircleOutlineIcon from "@mui/icons-material/RemoveCircleOutline";
import { IMaskInput } from "react-imask";
import AddIcon from "@mui/icons-material/Add";
import { Node_Type } from "../../../generated";
import { ExportOptionsType } from "../../../types";

import {
    Grid,
    Stack,
    Button,
    Dialog,
    Select,
    MenuItem,
    TextField,
    IconButton,
    Divider,
    InputLabel,
    Typography,
    DialogTitle,
    FormControl,
    DialogActions,
    DialogContent,
    OutlinedInput,
    DialogContentText,
} from "@mui/material";
const validationSchema = Yup.object().shape({
    nodes: Yup.array().of(
        Yup.object().shape({
            nodeName: Yup.string().required("Node name is required"),
            nodeNumber: Yup.string().required("Node number is required"),
            nodeType: Yup.string().required("Node type is required"),
        })
    ),
});
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
const debug = true;
const useStyles = makeStyles((theme: any) => ({
    basicDialogHeaderStyle: {
        padding: "0px 0px 18px 0px",
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
    field: {
        margin: theme.spacing(1),
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
const removeNodeType = function (arr: any, attr: any, value: any) {
    var i = arr.length;
    while (i--) {
        if (
            arr[i] &&
            // eslint-disable-next-line no-prototype-builtins
            arr[i].hasOwnProperty(attr) &&
            arguments.length > 2 &&
            arr[i][attr] === value
        ) {
            arr.splice(i, 1);
        }
    }
    return arr;
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
    const [formData, setFormData] = useState({
        type: nodeData.type,
        orgId: nodeData.orgId,
    });
    const [nodeType, setNodeType] = useState<any>();
    console.log("NODE", nodeType);
    const gclasses = globalUseStyles();
    const [isArrayfield, setIsArrayfield] = useState(0);
    console.log("NUMBER", isArrayfield);
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

            <Formik
                initialValues={{
                    nodes: [
                        {
                            id: Math.random(),
                            nodeName: "",
                            nodeNumber: "",
                            nodeType: "HOME",
                        },
                    ],
                }}
                validationSchema={validationSchema}
                onSubmit={values => {
                    // const data = {
                    //     ...values,
                    //     type: formData.type,
                    //     orgId: formData.orgId,
                    // };s
                    // console.log("onSubmit", JSON.stringify(data));
                    console.log("onSubmit", JSON.stringify(values, null, 2));
                }}
            >
                {({
                    values,
                    touched,
                    errors,
                    handleChange,
                    handleBlur,
                    isValid,
                }) => (
                    <Form noValidate autoComplete="off">
                        <FieldArray name="nodes">
                            {({ push, remove }) => (
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
                                    <div>
                                        {values.nodes.map((n, index) => {
                                            const nodeName = `nodes[${index}].nodeName`;
                                            const touchedNodeName = getIn(
                                                touched,
                                                nodeName
                                            );
                                            const errorNodeName = getIn(
                                                errors,
                                                nodeName
                                            );

                                            const nodeNumber = `nodes[${index}].nodeNumber`;
                                            const touchedNodeNumber = getIn(
                                                touched,
                                                nodeNumber
                                            );
                                            const errorNodeNumber = getIn(
                                                errors,
                                                nodeNumber
                                            );

                                            const nodeType = `nodes[${index}].nodeType`;
                                            const touchedNodeType = getIn(
                                                touched,
                                                nodeType
                                            );
                                            const errorNodeType = getIn(
                                                errors,
                                                nodeName
                                            );
                                            setNodeType(n.nodeType);

                                            return (
                                                <>
                                                    <Grid
                                                        container
                                                        rowSpacing={2.75}
                                                        columnSpacing={2.75}
                                                        mt={2}
                                                        key={n.id}
                                                    >
                                                        <Grid
                                                            item
                                                            xs={12}
                                                            md={6}
                                                        >
                                                            <FormControl
                                                                variant="outlined"
                                                                className={
                                                                    classes.formControl
                                                                }
                                                            >
                                                                <InputLabel
                                                                    shrink
                                                                    variant="outlined"
                                                                    htmlFor="outlined-age-always-notched"
                                                                >
                                                                    NODE TYPE
                                                                </InputLabel>
                                                                <Select
                                                                    value={
                                                                        n.nodeType
                                                                    }
                                                                    name={
                                                                        nodeType
                                                                    }
                                                                    // error={Boolean(
                                                                    //     touchedNodeType &&
                                                                    //         errorNodeType
                                                                    // )}
                                                                    // helperText={
                                                                    //     touchedNodeType &&
                                                                    //     errorNodeType
                                                                    //         ? errorNodeType
                                                                    //         : ""
                                                                    // }
                                                                    variant="outlined"
                                                                    disabled={
                                                                        action ==
                                                                        "editNode"
                                                                    }
                                                                    onChange={
                                                                        handleChange
                                                                    }
                                                                    input={
                                                                        <OutlinedInput
                                                                            notched
                                                                            label="NODE TYPE"
                                                                            name={
                                                                                nodeType
                                                                            }
                                                                            id="outlined-age-always-notched"
                                                                        />
                                                                    }
                                                                    MenuProps={{
                                                                        disablePortal:
                                                                            false,
                                                                        PaperProps:
                                                                            {
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
                                                                    {NODE_TYPE.map(
                                                                        ({
                                                                            id,
                                                                            value,
                                                                            label,
                                                                        }: ExportOptionsType) => (
                                                                            <MenuItem
                                                                                key={
                                                                                    id
                                                                                }
                                                                                value={
                                                                                    value
                                                                                }
                                                                                sx={{
                                                                                    m: 0,
                                                                                    p: "6px 16px",
                                                                                }}
                                                                            >
                                                                                <Typography variant="body1">
                                                                                    {
                                                                                        label
                                                                                    }
                                                                                </Typography>
                                                                            </MenuItem>
                                                                        )
                                                                    )}
                                                                </Select>
                                                            </FormControl>
                                                        </Grid>
                                                        <Grid
                                                            item
                                                            xs={12}
                                                            md={6}
                                                        >
                                                            <TextField
                                                                label="NODE NAME"
                                                                name={nodeName}
                                                                value={
                                                                    n.nodeName
                                                                }
                                                                helperText={
                                                                    touchedNodeName &&
                                                                    errorNodeName
                                                                        ? errorNodeName
                                                                        : ""
                                                                }
                                                                InputLabelProps={{
                                                                    shrink: true,
                                                                }}
                                                                fullWidth
                                                                error={Boolean(
                                                                    touchedNodeName &&
                                                                        errorNodeName
                                                                )}
                                                                onChange={
                                                                    handleChange
                                                                }
                                                                onBlur={
                                                                    handleBlur
                                                                }
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
                                                                required
                                                                label={
                                                                    "NODE NUMBER"
                                                                }
                                                                value={
                                                                    n.nodeNumber
                                                                }
                                                                disabled={
                                                                    action ==
                                                                    "editNode"
                                                                }
                                                                helperText={
                                                                    touchedNodeNumber &&
                                                                    errorNodeNumber
                                                                        ? errorNodeNumber
                                                                        : ""
                                                                }
                                                                error={Boolean(
                                                                    touchedNodeNumber &&
                                                                        errorNodeNumber
                                                                )}
                                                                onChange={
                                                                    handleChange
                                                                }
                                                                onBlur={
                                                                    handleBlur
                                                                }
                                                                name={
                                                                    nodeNumber
                                                                }
                                                                InputLabelProps={{
                                                                    shrink: true,
                                                                }}
                                                                // id="formatted-text-mask-input"
                                                                // spellCheck={false}
                                                                InputProps={{
                                                                    // inputComponent:
                                                                    //     TextMaskCustom as any,
                                                                    classes: {
                                                                        input: gclasses.inputFieldStyle,
                                                                    },
                                                                }}
                                                            />
                                                        </Grid>
                                                        {action ==
                                                            "editNode" && (
                                                            <Grid item xs={12}>
                                                                <TextField
                                                                    fullWidth
                                                                    value={
                                                                        formData.orgId
                                                                    }
                                                                    disabled={
                                                                        true
                                                                    }
                                                                    label={
                                                                        "ORGANIZATION ID"
                                                                    }
                                                                    InputLabelProps={{
                                                                        shrink: true,
                                                                    }}
                                                                    InputProps={{
                                                                        classes:
                                                                            {
                                                                                input: gclasses.inputFieldStyle,
                                                                            },
                                                                    }}
                                                                />
                                                            </Grid>
                                                        )}
                                                        {isArrayfield >= 1 && (
                                                            <>
                                                                <Grid
                                                                    item
                                                                    xs={12}
                                                                >
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
                                                                                fontStyle:
                                                                                    "Bold",
                                                                            }}
                                                                        >
                                                                            AMPLIFIER
                                                                            NODE
                                                                        </Typography>
                                                                        <IconButton
                                                                            color="primary"
                                                                            aria-label="remove-node"
                                                                            component="span"
                                                                            onClick={() =>
                                                                                remove(
                                                                                    index
                                                                                )
                                                                            }
                                                                        >
                                                                            <RemoveCircleOutlineIcon />
                                                                        </IconButton>
                                                                    </Stack>
                                                                </Grid>
                                                            </>
                                                        )}
                                                    </Grid>
                                                </>
                                            );
                                        })}

                                        {(nodeType == "TOWER" ||
                                            nodeType == "") && (
                                            <>
                                                <Grid
                                                    item
                                                    xs={12}
                                                    sx={{ pt: 2 }}
                                                >
                                                    <Button
                                                        variant="text"
                                                        startIcon={<AddIcon />}
                                                        sx={{
                                                            color: colors.primaryMain,
                                                            pointer: "cursor",
                                                        }}
                                                        onClick={() => {
                                                            push({
                                                                id: Math.random(),
                                                                nodeName: "",
                                                                nodeNumber: "",
                                                                nodeType: "",
                                                            });
                                                            setIsArrayfield(
                                                                isArrayfield + 1
                                                            );
                                                        }}
                                                    >
                                                        add amplifier node
                                                    </Button>
                                                </Grid>
                                            </>
                                        )}
                                    </div>
                                </DialogContent>
                            )}
                        </FieldArray>

                        <DialogActions>
                            <Button sx={{ mr: 2 }} onClick={handleClose}>
                                Cancel
                            </Button>
                            <Button
                                variant="contained"
                                type="submit"
                                className={classes.stepButtonStyle}
                                // disabled={!isValid || values.nodes.length === 0}
                            >
                                {action == "editNode"
                                    ? "UPDATE NODE"
                                    : "REGISTER NODE"}
                            </Button>
                            {/* {JSON.stringify(values.nodes[0], null, 2)} */}
                        </DialogActions>
                    </Form>
                )}
            </Formik>
        </Dialog>
    );
};

export default NodeDialog;
