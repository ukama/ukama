import {
    Box,
    Stack,
    Theme,
    Select,
    Button,
    Divider,
    MenuItem,
    SelectChangeEvent,
    Typography,
} from "@mui/material";
import { LoadingWrapper } from "..";
import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";
import { getStatusByType, hexToRGB } from "../../utils";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import { NodeDto, Org_Node_State } from "../../generated";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";

const useStyles = makeStyles<Theme>(() => ({
    selectStyle: () => ({
        width: "fit-content",
    }),
}));

const getStatusIcon = (status: Org_Node_State) => {
    switch (status) {
        case Org_Node_State.Onboarded:
            return <CheckCircleIcon htmlColor={colors.green} />;
        case Org_Node_State.Pending:
            return <InfoIcon htmlColor={colors.yellow} />;
        default:
            return <InfoIcon htmlColor={colors.error} />;
    }
};

interface INodeDropDown {
    loading: boolean;
    onAddNode: Function;
    nodes: NodeDto[] | [];
    onNodeSelected: Function;
    selectedNode: NodeDto | undefined;
}

const NodeDropDown = ({
    nodes = [],
    onAddNode,
    selectedNode,
    loading = true,
    onNodeSelected,
}: INodeDropDown) => {
    const classes = useStyles();
    const handleChange = (e: SelectChangeEvent<string>) => {
        const { target } = e;
        target.value &&
            onNodeSelected(
                nodes.find((item: NodeDto) => item.title === target.value)
            );
    };
    return (
        <LoadingWrapper isLoading={loading} height={40} width={"30%"}>
            <Stack direction={"row"} spacing={1}>
                {selectedNode && (
                    <Box component="div" display={"flex"} alignItems={"center"}>
                        {getStatusIcon(selectedNode?.status)}
                    </Box>
                )}
                <Select
                    disableUnderline
                    variant="standard"
                    onChange={handleChange}
                    value={selectedNode?.title}
                    SelectDisplayProps={{
                        style: {
                            fontWeight: 600,
                            display: "flex",
                            fontSize: "20px",
                            marginRight: "4px",
                            alignItems: "center",
                            minWidth: "fit-content",
                        },
                    }}
                    MenuProps={{
                        disablePortal: true,
                        anchorOrigin: {
                            vertical: "bottom",
                            horizontal: "left",
                        },
                        transformOrigin: {
                            vertical: "top",
                            horizontal: "left",
                        },
                        PaperProps: {
                            sx: {
                                boxShadow:
                                    "0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)",
                                borderRadius: "4px",
                            },
                        },
                    }}
                    className={classes.selectStyle}
                    renderValue={selected => selected}
                >
                    {nodes.map(({ id, title }) => (
                        <MenuItem
                            key={id}
                            value={title}
                            sx={{
                                m: 0,
                                p: "6px 16px",

                                backgroundColor: `${
                                    id === selectedNode?.id
                                        ? hexToRGB(colors.secondaryLight, 0.25)
                                        : "inherit"
                                } !important`,
                                ":hover": {
                                    backgroundColor: `${hexToRGB(
                                        colors.secondaryLight,
                                        0.25
                                    )} !important`,
                                },
                            }}
                        >
                            <Typography variant="body1">{title}</Typography>
                        </MenuItem>
                    ))}
                    <Divider />
                    <MenuItem
                        onClick={e => {
                            onAddNode();
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
                            Add node
                        </Button>
                    </MenuItem>
                </Select>

                <Box
                    component="div"
                    display="flex"
                    flexDirection="row"
                    alignItems="center"
                >
                    <Typography variant={"h6"}>
                        {getStatusByType(selectedNode?.status as string)}
                    </Typography>

                    {selectedNode?.status !== "UNDEFINED" && (
                        <Typography ml="8px" variant={"h6"} color="secondary">
                            21 days 5 hours 1 minute
                        </Typography>
                    )}
                </Box>
            </Stack>
        </LoadingWrapper>
    );
};

export default NodeDropDown;
