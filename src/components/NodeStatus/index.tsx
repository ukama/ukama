import { colors } from "../../theme";
import {
    Box,
    Button,
    Stack,
    Select,
    MenuItem,
    Typography,
    SelectChangeEvent,
} from "@mui/material";
import { useRecoilValue } from "recoil";
import { isDarkmode } from "../../recoil";
import { makeStyles, styled } from "@mui/styles";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import { HorizontalContainerJustify } from "../../styles";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { GetNodesByOrgQuery, NodeDto, Org_Node_State } from "../../generated";

const StyledBtn = styled(Button)({
    whiteSpace: "nowrap",
    minWidth: "max-content",
});

type StyleProps = { color: string };

const useStyles = makeStyles(() => ({
    selectStyle: ({ color = colors.black }: StyleProps) => ({
        width: "144px",
        "& p": {
            fontSize: "20px",
            fontWeight: "500",
            lineHeight: "160%",
            fontFamily: "Rubik",
            letterSpacing: "0.15px",
            color: color,
        },
    }),
}));

const getStatusIcon = (status: Org_Node_State) => {
    switch (status) {
        case Org_Node_State.Onboarded:
            return <CheckCircleIcon htmlColor={colors.green} />;
        case Org_Node_State.Pending:
            return <InfoIcon htmlColor={colors.yellow} />;
        default:
            return <InfoIcon />;
    }
};

interface INodeStatus {
    nodes?: GetNodesByOrgQuery["getNodesByOrg"]["nodes"] | any;
    selectedNode: NodeDto | undefined;
    onNodeRFClick: Function;
    onNodeSelected: Function;
    onNodeSwitchClick: Function;
    onRestartNodeClick: Function;
}

const NodeStatus = ({
    nodes,
    selectedNode = {
        id: "1",
        totalUser: 4,
        title: "Node 1",
        status: Org_Node_State.Undefined,
        description: "Node 1 description",
    },
    onNodeSelected,
    onNodeRFClick,
    onNodeSwitchClick,
    onRestartNodeClick,
}: INodeStatus) => {
    const _isDarkMod = useRecoilValue(isDarkmode);
    const styleProps = { color: _isDarkMod ? colors.white : colors.black };
    const classes = useStyles(styleProps);

    const handleRestartNode = () =>
        onRestartNodeClick(
            nodes.find((item: NodeDto) => item.id === selectedNode.id)
        );

    const handleRFClick = () =>
        onNodeRFClick(
            nodes.find((item: NodeDto) => item.id === selectedNode.id)
        );

    const handleNodeSwitch = () =>
        onNodeSwitchClick(
            nodes.find((item: NodeDto) => item.id === selectedNode.id)
        );

    const handleChange = (e: SelectChangeEvent<string>) =>
        onNodeSelected(
            nodes.find((item: NodeDto) => item.id === e.target.value)
        );

    return (
        <HorizontalContainerJustify>
            <Stack direction={"row"} spacing={2}>
                <Box display={"flex"} alignItems={"center"}>
                    {getStatusIcon(selectedNode?.status)}
                </Box>
                <Select
                    sx={{
                        width: "auto",
                    }}
                    disableUnderline
                    variant="standard"
                    value={selectedNode.id}
                    onChange={handleChange}
                    className={classes.selectStyle}
                >
                    {nodes.map(({ id, title }: NodeDto) => (
                        <MenuItem key={id} value={id}>
                            <Typography variant="body1">{title}</Typography>
                        </MenuItem>
                    ))}
                </Select>
            </Stack>
            <Stack direction={"row"} spacing={2}>
                <StyledBtn variant="contained" onClick={handleRestartNode}>
                    restart
                </StyledBtn>
                <StyledBtn variant="contained" onClick={handleRFClick}>
                    turn rf off
                </StyledBtn>
                <StyledBtn variant="contained" onClick={handleNodeSwitch}>
                    turn node off
                </StyledBtn>
            </Stack>
        </HorizontalContainerJustify>
    );
};

export default NodeStatus;
