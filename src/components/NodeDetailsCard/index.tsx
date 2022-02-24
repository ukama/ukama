import { LoadingWrapper } from "..";
import colors from "../../theme/colors";
import DeviceModalView from "../DeviceModalView";
import { Chip, Paper, Stack, Typography } from "@mui/material";
import { HorizontalContainerJustify, LinkStyle } from "../../styles";

interface INodeDetailsCard {
    loading: boolean;
    nodeTitle: string;
    isUpdateAvailable: boolean;
    handleUpdateNode: Function;
}

const NodeDetailsCard = ({
    loading,
    nodeTitle,
    isUpdateAvailable,
}: INodeDetailsCard) => {
    return (
        <LoadingWrapper
            width="100%"
            height="100%"
            radius={"small"}
            isLoading={loading}
        >
            <Paper sx={{ p: 3, gap: 1 }}>
                <Stack spacing={3}>
                    <HorizontalContainerJustify>
                        <Stack direction={"row"} spacing={2}>
                            <Typography variant="h6">{nodeTitle}</Typography>
                            {isUpdateAvailable && (
                                <Chip
                                    variant="outlined"
                                    sx={{
                                        color: colors.primaryMain,
                                        border: `1px solid ${colors.primaryMain}`,
                                    }}
                                    label={
                                        <LinkStyle
                                            underline="hover"
                                            sx={{ fontSize: "14px" }}
                                        >
                                            Software update available -- view
                                            notes
                                        </LinkStyle>
                                    }
                                />
                            )}
                        </Stack>
                    </HorizontalContainerJustify>
                    {/* <DeviceModalView /> */}
                </Stack>
            </Paper>
        </LoadingWrapper>
    );
};

export default NodeDetailsCard;
