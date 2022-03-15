import React from "react";
import { LoadingWrapper } from "..";
import { hexToRGB } from "../../utils";
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

const Line = styled(Divider)(props => ({
    margin: "18px -18px 4px -18px",
    background: hexToRGB(props.theme.palette.text.primary, 0.3),
}));

const IconStyle = {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
};

type NodeCardProps = {
    id: string;
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
            height={206}
        >
            <Card
                variant="outlined"
                sx={{
                    width: "218px",
                    height: "206px",
                    borderRadius: "10px",
                    justifyContent: "center",
                    padding: "15px 18px 8px 18px",
                }}
            >
                <Grid container spacing={0.8}>
                    <Grid item xs={10}>
                        <Grid textAlign="initial">
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
                        <Grid textAlign="initial">
                            <Typography variant="caption">
                                {subTitle}
                            </Typography>
                        </Grid>
                    </Grid>
                    <Grid
                        item
                        xs={2}
                        m="4px 0px"
                        display={"flex"}
                        alignItems="flex-start"
                        justifyContent={"flex-end"}
                    >
                        {isUpdateAvailable ? (
                            <Tooltip title={updateShortNote}>
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
                                style={{
                                    cursor: "pointer",
                                    position: "relative",
                                    bottom: "13px",
                                }}
                                cid={"node-card-options"}
                                menuOptions={BASIC_MENU_ACTIONS}
                                handleItemClick={(type: string) =>
                                    handleOptionItemClick(type)
                                }
                            />
                        )}
                    </Grid>
                    <Grid item xs={12} sx={{ ...IconStyle }}>
                        <img
                            src="https://ovalloqu.sirv.com/Images/node.png"
                            width="86"
                            height="76"
                            alt="node-img"
                        />
                    </Grid>
                    <Grid item xs={12}>
                        <Line />
                    </Grid>
                    <Grid
                        item
                        xs={12}
                        container
                        spacing={1}
                        sx={{ alignItems: "center" }}
                    >
                        <Grid
                            item
                            xs={6}
                            container
                            display="flex"
                            alignSelf="end"
                            pt="0px !important"
                        >
                            <UsersIcon sx={{ width: "16px", height: "16px" }} />
                            <Typography variant="caption" pl="8px">
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
