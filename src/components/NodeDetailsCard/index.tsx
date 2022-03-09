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
    getNodeUpdateInfos: Function;
}

const NodeDetailsCard = ({
    loading,
    nodeTitle,
    isUpdateAvailable,
    getNodeUpdateInfos,
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
                                        <>
                                            <Stack direction="row">
                                                <Typography variant="body2">
                                                    Software update available —
                                                    view view
                                                </Typography>

                                                <LinkStyle
                                                    underline="hover"
                                                    onClick={() =>
                                                        getNodeUpdateInfos()
                                                    }
                                                    sx={{
                                                        fontSize: "14px",
                                                        paddingLeft: 1,
                                                        cursor: "pointer",
                                                        textDecoration:
                                                            "underline",
                                                        color: colors.primaryDark,
                                                    }}
                                                >
                                                    notes
                                                </LinkStyle>
                                            </Stack>
                                        </>
                                    }
                                />
                            )}
                        </Stack>
                    </HorizontalContainerJustify>
                    <DeviceModalView />
                </Stack>
            </Paper>
        </LoadingWrapper>
    );
};

export default NodeDetailsCard;
