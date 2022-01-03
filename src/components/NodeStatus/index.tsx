import { colors } from "../../theme";
import { styled } from "@mui/styles";
import TextSelect from "../TextSelect";
import InfoIcon from "@mui/icons-material/Info";
import CancelIcon from "@mui/icons-material/Cancel";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { Box, Typography, Grid, Button, Stack } from "@mui/material";

const StyledBtn = styled(Button)({
    whiteSpace: "nowrap",
    minWidth: "max-content",
});

const getStatusDetails = (status: string) => {
    switch (status) {
        case "BEING_CONFIGURED":
            return {
                showDuration: true,
                text: " is being configured.",
                icon: <InfoIcon htmlColor={colors.yellow} />,
            };

        case "ONLINE":
            return {
                showDuration: true,
                text: " is online and well for ",
                icon: <CheckCircleIcon htmlColor={colors.green} />,
            };

        default:
            return {
                showDuration: false,
                text: " has something wrong.",
                icon: <CancelIcon htmlColor={colors.red} />,
            };
    }
};
type Node = { name: string; duration: string; statusType: string };

type NodeStatusProps = {
    nodes: Node[];
    selectedNodeIndex: number;
    setSelectedNodeIndex: (_index: number) => void;
};

const NodeStatus = ({
    nodes,
    selectedNodeIndex,
    setSelectedNodeIndex,
}: NodeStatusProps) => {
    const { icon, text, showDuration } = getStatusDetails(
        nodes[selectedNodeIndex].statusType
    );

    const handleRestartNode = () => {
        return;
    };
    const handleTurnRFOff = () => {
        return;
    };
    const handleTurnNodeOff = () => {
        return;
    };

    return (
        <Grid width="100%" container pt="18px">
            <Grid item xs={12} lg={7}>
                <Box
                    height="100%"
                    display="flex"
                    flexWrap="wrap"
                    flexDirection="row"
                    alignItems="center"
                >
                    {icon}
                    <TextSelect
                        value={selectedNodeIndex}
                        setValue={setSelectedNodeIndex}
                        options={nodes.map(({ name }) => name)}
                    />
                    <Typography variant={"h6"}>{text}</Typography>
                    {showDuration && (
                        <Typography
                            ml="8px"
                            variant={"h6"}
                            color={colors.empress}
                        >
                            {nodes[selectedNodeIndex].duration}
                        </Typography>
                    )}
                </Box>
            </Grid>
            <Grid item xs={12} lg={5} display="flex" justifyContent="flex-end">
                <Stack
                    direction={{
                        xs: "column",
                        sm: "row",
                    }}
                    spacing={2}
                    width="100%"
                    justifyContent="flex-end"
                >
                    <StyledBtn
                        variant="contained"
                        onClick={() => handleRestartNode()}
                    >
                        restart
                    </StyledBtn>
                    <StyledBtn
                        variant="contained"
                        onClick={() => handleTurnRFOff()}
                    >
                        turn rf off
                    </StyledBtn>
                    <StyledBtn
                        variant="contained"
                        onClick={() => handleTurnNodeOff()}
                    >
                        turn node off
                    </StyledBtn>
                </Stack>
            </Grid>
        </Grid>
    );
};

export default NodeStatus;
