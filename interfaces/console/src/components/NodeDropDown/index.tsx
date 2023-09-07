import {
    Stack,
    Select,
    Button,
    Divider,
    MenuItem,
    Typography,
    SelectChangeEvent,
} from "@mui/material";
import { LoadingWrapper } from "..";
import { colors } from "../../theme";
import CircleIcon from "@mui/icons-material/Circle";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import { hexToRGB, secToHoursNMints } from "../../utils";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { PaperProps, SelectDisplayProps, useStyles } from "./styles";
import { GetNodeStatusRes, NodeDto, Org_Node_State } from "../../generated";

const getStatus = (status: Org_Node_State, time: number) => {
    switch (status) {
        case Org_Node_State.Onboarded:
            return (
                <Stack display="flex" flexDirection="row" alignItems={"center"}>
                    <Typography variant={"h6"} mr={"6px"}>
                        is online and well for
                    </Typography>
                    <Typography variant={"h6"} color="primary">
                        {secToHoursNMints(time, " hours and ")}
                    </Typography>
                </Stack>
            );

        case Org_Node_State.Pending:
            return <Typography variant={"h6"}>is configuring.</Typography>;

        default:
            return "";
    }
};

const getStatusIcon = (status: Org_Node_State) => {
    switch (status) {
        case Org_Node_State.Onboarded:
            return (
                <CheckCircleIcon htmlColor={colors.green} fontSize={"small"} />
            );
        case Org_Node_State.Pending:
            return <InfoIcon htmlColor={colors.yellow} fontSize={"small"} />;
        case Org_Node_State.Error:
            return <InfoIcon htmlColor={colors.red} fontSize={"small"} />;
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
        <Stack direction={"row"} spacing={1} alignItems="center">
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
                    SelectDisplayProps={SelectDisplayProps}
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
                                ...PaperProps,
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
                {getStatus(nodeStatus.status, nodeStatus.uptime)}
            </LoadingWrapper>
        </Stack>
    );
};

export default NodeDropDown;
