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
import { LoadingWrapper } from "..";

const StyledBtn = styled(Button)({
    whiteSpace: "nowrap",
    minWidth: "max-content",
});

type StyleProps = { color: string };

const useStyles = makeStyles(() => ({
    selectStyle: ({ color = colors.black }: StyleProps) => ({
        width: "fit-content",
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
    loading: boolean;
    onNodeRFClick: Function;
    onNodeSelected: Function;
    onNodeSwitchClick: Function;
    onRestartNodeClick: Function;
    selectedNode: NodeDto | undefined;
    nodes?: GetNodesByOrgQuery["getNodesByOrg"]["nodes"] | any;
}

const NodeStatus = ({
    nodes,
    loading = false,
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
    const styleProps = { color: _isDarkMod ? colors._white : colors.black };
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
            <LoadingWrapper isLoading={loading} height={40} width={"30%"}>
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
                        placeholder="Select Node"
                        className={classes.selectStyle}
                    >
                        {nodes?.map(({ id, title }: NodeDto) => (
                            <MenuItem key={id} value={id}>
                                <Typography variant="body1">{title}</Typography>
                            </MenuItem>
                        ))}
                    </Select>
                </Stack>
            </LoadingWrapper>
            <Stack direction={"row"} spacing={2}>
                <LoadingWrapper isLoading={loading} height={40} width={100}>
                    <StyledBtn variant="contained" onClick={handleRestartNode}>
                        restart
                    </StyledBtn>
                </LoadingWrapper>
                <LoadingWrapper isLoading={loading} height={40} width={100}>
                    <StyledBtn variant="contained" onClick={handleRFClick}>
                        turn rf off
                    </StyledBtn>
                </LoadingWrapper>
                <LoadingWrapper isLoading={loading} height={40} width={100}>
                    <StyledBtn variant="contained" onClick={handleNodeSwitch}>
                        turn node off
                    </StyledBtn>
                </LoadingWrapper>
            </Stack>
        </HorizontalContainerJustify>
    );
};

export default NodeStatus;
