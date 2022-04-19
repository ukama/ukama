import React from "react";
import { LoadingWrapper } from "..";
import OptionsPopover from "../OptionsPopover";
import { BASIC_MENU_ACTIONS } from "../../constants";
import UsersIcon from "@mui/icons-material/PeopleAlt";
import {
    Typography,
    Grid,
    Divider,
    styled,
    Card,
    IconButton,
    Tooltip,
} from "@mui/material";
import UpdateIcon from "@mui/icons-material/SystemUpdateAltRounded";
import { Node_Type } from "../../generated";

const CpuIcon = React.lazy(() =>
    import("../../assets/svg").then(module => ({
        default: module.CpuIcon,
    }))
);
const BatteryIcon = React.lazy(() =>
    import("../../assets/svg").then(module => ({
        default: module.BatteryIcon,
    }))
);
const ThermometerIcon = React.lazy(() =>
    import("../../assets/svg").then(module => ({
        default: module.ThermometerIcon,
    }))
);

const Line = styled(Divider)(() => ({
    margin: "18px -18px 4px -18px",
    background: "rgba(255, 255, 255, 0.12)",
}));

const NODE_IMAGES = {
    TOWER: "https://i.ibb.co/qRRFz9j/tower-node.png",
    AMPLIFIER: "https://i.ibb.co/NKtrHzG/amplifier-node.png",
    HOME: "https://i.ibb.co/G0gqqb5/home-node.png",
};

const IconStyle = {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
};

type NodeCardProps = {
    id: string;
    type: Node_Type;
    title: string;
    users?: number;
    loading?: boolean;
    subTitle: string;
    isConfigure?: boolean;
    updateShortNote: string;
    isUpdateAvailable: boolean;
    handleNodeUpdate: Function;
    handleOptionItemClick?: Function;
};

const NodeCard = ({
    id,
    type,
    title,
    users,
    subTitle,
    loading,
    updateShortNote = "",
    isUpdateAvailable = false,
    handleOptionItemClick = () => {
        /* Default empty function */
    },
    handleNodeUpdate = () => {
        /* Default empty function */
    },
}: NodeCardProps) => {
    return (
        <LoadingWrapper
            isLoading={loading}
            variant="rectangular"
            width={218}
            height={216}
        >
            <Card
                variant="outlined"
                sx={{
                    width: "218px",
                    height: "216px",
                    borderRadius: "10px",
                    justifyContent: "center",
                    padding: "15px 18px 8px 18px",
                }}
            >
                <Grid container>
                    <Grid item xs={10}>
                        <Typography
                            variant="subtitle1"
                            sx={{
                                height: "20px",
                                fontWeight: 500,
                                overflow: "hidden",
                                lineHeight: "19px",
                                display: "-webkit-box",
                                letterSpacing: "-0.02em",
                                textOverflow: "ellipsis",
                            }}
                        >
                            {title}
                        </Typography>
                    </Grid>
                    <Grid
                        item
                        xs={2}
                        display={"flex"}
                        alignItems="center"
                        justifyContent={"flex-end"}
                    >
                        {isUpdateAvailable ? (
                            <Tooltip title={updateShortNote} arrow>
                                <IconButton
                                    sx={{ p: 0, cursor: "pointer" }}
                                    onClick={() => handleNodeUpdate(id)}
                                >
                                    <UpdateIcon
                                        fontSize="small"
                                        color="primary"
                                    />
                                </IconButton>
                            </Tooltip>
                        ) : (
                            <OptionsPopover
                                style={{ cursor: "pointer" }}
                                cid={"node-card-options"}
                                menuOptions={BASIC_MENU_ACTIONS}
                                handleItemClick={(type: string) =>
                                    handleOptionItemClick(type)
                                }
                            />
                        )}
                    </Grid>
                    <Grid item xs={12} textAlign="initial">
                        <Typography variant="caption">{subTitle}</Typography>
                    </Grid>

                    <Grid
                        item
                        xs={12}
                        minHeight={"92px"}
                        sx={{ ...IconStyle, py: 1 }}
                    >
                        <img
                            src={NODE_IMAGES[type]}
                            alt="node-img"
                            style={{ maxWidth: "180px", maxHeight: "78px" }}
                        />
                    </Grid>
                    <Grid item xs={12} mb={0.8}>
                        <Line />
                    </Grid>
                    <Grid
                        item
                        xs={12}
                        container
                        spacing={1}
                        mb={"2px"}
                        sx={{ alignItems: "center" }}
                    >
                        <Grid
                            item
                            xs={6}
                            container
                            display="flex"
                            alignSelf="end"
                            pt="0px !important"
                            alignItems={"flex-end"}
                        >
                            <UsersIcon sx={{ width: "16px", height: "16px" }} />
                            <Typography
                                pl="8px"
                                variant="caption"
                                lineHeight={"14px"}
                            >
                                {users}
                            </Typography>
                        </Grid>
                        <Grid xs={2} item sx={{ ...IconStyle }}>
                            <ThermometerIcon />
                        </Grid>
                        <Grid xs={2} item sx={{ ...IconStyle }}>
                            <BatteryIcon />
                        </Grid>
                        <Grid xs={2} item sx={{ ...IconStyle }}>
                            <CpuIcon />
                        </Grid>
                    </Grid>
                </Grid>
            </Card>
        </LoadingWrapper>
    );
};

export default NodeCard;
