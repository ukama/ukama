import { STATS_PERIOD } from "../../constants";
import NodePropertyStats from "../NodePropertyStats";
import { NodePlaceholder } from "../../assets/images";
import { Box, Button, styled, Typography, Paper } from "@mui/material";

const StyledBtn = styled(Button)({
    whiteSpace: "nowrap",
    minWidth: "max-content",
});

const Image = styled("img")({
    width: "100%",
});

interface NodeDetailsPanelProps {
    index: number;
    value: number;
    title: string;
    renderPropertyStats?: boolean;
    image?: { src: string; alt: string };
    button?: { label: string; onClick: Function };
    properties: Array<{ name: string; value: string | number }>;
}

const NodeDetailsPanel = ({
    image,
    title,
    value,
    index,
    button,
    properties,
    renderPropertyStats = true,
}: NodeDetailsPanelProps) => {
    return (
        <Box
            width="100%"
            overflow="hidden"
            borderRadius="5px"
            role="detailpanel"
            hidden={value !== index}
            id={`detailpanel-${index}`}
            aria-labelledby={`detailpanel-${index}`}
        >
            {value === index && (
                <Paper sx={{ p: "26px 28px" }}>
                    <Box
                        display="flex"
                        justifyContent="space-between"
                        mb="18px"
                    >
                        <Typography variant={"h6"}>{title}</Typography>
                        {button && (
                            <StyledBtn
                                variant="contained"
                                onClick={() => button.onClick()}
                            >
                                {button.label}
                            </StyledBtn>
                        )}
                    </Box>
                    {image && <Image alt={image.alt} src={NodePlaceholder} />}
                    {renderPropertyStats &&
                        properties.map((propery, keyIndex) => (
                            <NodePropertyStats
                                key={keyIndex}
                                propery={propery}
                                periodOptions={STATS_PERIOD}
                            />
                        ))}
                </Paper>
            )}
        </Box>
    );
};

export default NodeDetailsPanel;
