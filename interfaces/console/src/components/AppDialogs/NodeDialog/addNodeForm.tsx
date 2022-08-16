import {
    Grid,
    Stack,
    Button,
    Select,
    MenuItem,
    TextField,
    IconButton,
    InputLabel,
    Typography,
    FormControl,
    OutlinedInput,
    Divider,
} from "@mui/material";
import { colors } from "../../../theme";
import { IMaskInput } from "react-imask";
import { makeStyles } from "@mui/styles";
import AddIcon from "@mui/icons-material/Add";
import { Node_Type } from "../../../generated";
import { globalUseStyles } from "../../../styles";
import React, { Fragment, useEffect, useState } from "react";
import { SelectChangeEvent } from "@mui/material/Select/SelectInput";
import RemoveCircleOutlineIcon from "@mui/icons-material/RemoveCircleOutline";
import { MASK_BY_TYPE, MASK_PLACEHOLDERS, NODE_TYPE } from "../../../constants";

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
    nodeArray: any;
    towerNodesArrayList?: any;
    selectedToweNode?: any;
    isSubmitted: boolean;
    nodeType: string;
    handleAssociatedTowerNode: Function;
};

const AddNodeForm = ({
    nodeArray,
    isSubmitted,
    nodeType,
    selectedToweNode,
    towerNodesArrayList,
    handleAssociatedTowerNode,
}: NodeDialogProps) => {
    const classes = useStyles();
    const gclasses = globalUseStyles();
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
        nodeArray(list);
    };
    const handleRemoveClick = (index: number) => {
        const list = [...attachedAmplierNode];
        list.splice(index, 1);
        setAttachedAmplierNode(list);
    };

    const handleAddClick = () => {
        setAttachedAmplierNode([
            ...attachedAmplierNode,
            { nodeId: "", name: "" },
        ]);
    };
    const [showInitialAmplifierNode, setShowInitialAmplifierNode] =
        useState<boolean>(false);
    const [isAssociatedTowerNode, setIsAssociatedTowerNode] =
        useState<boolean>(false);
    const handleOptionalNodeType = (e: SelectChangeEvent) => {
        setNtype(e.target.value);
    };

    const onAddTowerNode = () => {
        setIsAssociatedTowerNode(true);
    };

    const removeNodeTypefromArray = (from: number, to: number) => {
        return NODE_TYPE.filter(function (value, index) {
            return [from, to].indexOf(index) == -1;
        });
    };

    const [nType, setNtype] = useState<any>(
        isAssociatedTowerNode ? "TOWER" : "AMPLIFIER"
    );
    const showAddButton = () => {
        setShowInitialAmplifierNode(true);
    };
    useEffect(() => {
        if (attachedAmplierNode.length == 0) {
            setShowInitialAmplifierNode(false);
            setAttachedAmplierNode([
                ...attachedAmplierNode,
                { nodeId: "", name: "" },
            ]);
        }
    }, [attachedAmplierNode]);
    return (
        <>
            {nodeType == "AMPLIFIER" && (
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
                            onChange={(e: any) => {
                                handleAssociatedTowerNode(e);
                            }}
                            value={selectedToweNode}
                            variant="outlined"
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
                            {towerNodesArrayList &&
                                towerNodesArrayList.map(({ id, name }: any) => (
                                    <MenuItem value={name} key={id}>
                                        <Typography variant="body1">
                                            {towerNodesArrayList.length
                                                ? name
                                                : "hello"}
                                        </Typography>
                                    </MenuItem>
                                ))}
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
            {nodeType == "TOWER" && showInitialAmplifierNode == false && (
                <Grid item xs={12}>
                    <Button
                        variant="text"
                        startIcon={<AddIcon />}
                        onClick={showAddButton}
                        sx={{
                            color: colors.primaryMain,
                            pointer: "cursor",
                        }}
                    >
                        add amplifier node
                    </Button>
                </Grid>
            )}
            {attachedAmplierNode.map((x: any, i: number) => {
                return (
                    (showInitialAmplifierNode == true ||
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
                                            fontWeight: "600",
                                        }}
                                    >
                                        {nodeType == "TOWER" &&
                                            "AMPLIFIER NODE"}

                                        {nodeType == "AMPLIFIER" &&
                                            "New Tower Node"}
                                    </Typography>
                                    <IconButton
                                        color="primary"
                                        aria-label="remove-node"
                                        component="span"
                                        onClick={() => handleRemoveClick(i)}
                                        sx={{
                                            position: "relative",
                                            bottom: 10,
                                            left: 10,
                                        }}
                                    >
                                        <RemoveCircleOutlineIcon />
                                    </IconButton>
                                </Stack>
                            </Grid>

                            <Fragment key={i}>
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
                                            variant="outlined"
                                            onChange={handleOptionalNodeType}
                                            value={nType}
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
                                                        borderRadius: "4px",
                                                    },
                                                },
                                            }}
                                            className={classes.selectStyle}
                                        >
                                            {removeNodeTypefromArray(
                                                0,
                                                nodeType == "TOWER" ? 1 : 2
                                            ).map(({ id, label, value }) => (
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
                                        label={"NODE NUMBER"}
                                        value={x.nodeId}
                                        onChange={(e: any) => {
                                            handleInputChange(e, i);
                                        }}
                                        InputLabelProps={{
                                            shrink: true,
                                        }}
                                        error={isSubmitted && x.nodeId === ""}
                                        helperText={
                                            isSubmitted && x.nodeId === ""
                                                ? "Node number is required !"
                                                : " "
                                        }
                                        id={"nodeId"}
                                        name={
                                            isAssociatedTowerNode
                                                ? "TOWER"
                                                : "AMPLIFIER"
                                        }
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
                                        error={isSubmitted && x.name === ""}
                                        helperText={
                                            isSubmitted && x.name === ""
                                                ? "Node name is required !"
                                                : " "
                                        }
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
        </>
    );
};

export default AddNodeForm;
