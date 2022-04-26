import {
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
import { hexToRGB } from "../../utils";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import CircleIcon from "@mui/icons-material/Circle";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { GetNodeStatusRes, NodeDto, Org_Node_State } from "../../generated";

const useStyles = makeStyles<Theme>(() => ({
    selectStyle: () => ({
        width: "fit-content",
    }),
}));

const getStatusIcon = (status: Org_Node_State) => {
    switch (status) {
        case Org_Node_State.Onboarded:
            return (
                <CheckCircleIcon htmlColor={colors.green} fontSize={"small"} />
            );
        case Org_Node_State.Pending:
            return <InfoIcon htmlColor={colors.yellow} fontSize={"small"} />;
        default:
            return <CircleIcon htmlColor={colors.black38} fontSize={"small"} />;
    }
};

interface INodeDropDown {
    loading: boolean;
    onAddNode: Function;
    nodes: NodeDto[] | [];
    onNodeSelected: Function;
    nodeStatusLoading: boolean;
    selectedNode: NodeDto | undefined;
    nodeStatus: GetNodeStatusRes | undefined;
}

const NodeDropDown = ({
    nodes = [],
    onAddNode,
    nodeStatus = {
        status: Org_Node_State.Undefined,
        uptime: new Date().getTime(),
    },
    selectedNode,
    loading = true,
    onNodeSelected,
    nodeStatusLoading,
}: INodeDropDown) => {
    const classes = useStyles();
    const handleChange = (e: SelectChangeEvent<string>) => {
        const { target } = e;
        target.value &&
            onNodeSelected(
                nodes.find((item: NodeDto) => item.name === target.value)
            );
    };
    return (
        <Stack direction={"row"} spacing={2} alignItems="center">
            {getStatusIcon(nodeStatus.status)}

            <LoadingWrapper
                height={38}
                isLoading={loading}
                width={loading ? "144px" : "fit-content"}
            >
                <Select
                    disableUnderline
                    variant="standard"
                    onChange={handleChange}
                    value={selectedNode?.name}
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
                    {nodes.map(({ id, name }) => (
                        <MenuItem
                            key={id}
                            value={name}
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
                            <Typography variant="body1">{name}</Typography>
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
            </LoadingWrapper>

            <LoadingWrapper
                height={38}
                isLoading={nodeStatusLoading}
                width={nodeStatusLoading ? "200px" : "fit-content"}
            >
                {nodeStatus.status !== Org_Node_State.Undefined && (
                    <Typography ml="8px" variant={"h6"} color="secondary">
                        {nodeStatus.uptime}
                    </Typography>
                )}
            </LoadingWrapper>
        </Stack>
    );
};

export default NodeDropDown;
