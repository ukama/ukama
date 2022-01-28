import {
    CpuIcon,
    NodeImg,
    BatteryIcon,
    ThermometerIcon,
} from "../../assets/svg";
import OptionsPopover from "../OptionsPopover";
import { SkeletonRoundedCard } from "../../styles";
import { BASIC_MENU_ACTIONS } from "../../constants";
import UsersIcon from "@mui/icons-material/PeopleAlt";
import { Typography, Grid, Divider, Box, styled } from "@mui/material";

const Container = styled(Box)(() => ({
    width: "214px",
    height: "206px",
    display: "flex",
    alignItems: "center",
    padding: "15px 18px 8px 18px",
    borderRadius: "10px",
    justifyContent: "center",
    // HERE background:
    //     props.theme.palette.mode === "dark" ? colors.nightGrey12 : colors.white,
}));

const IconStyle = {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
};

type NodeCardProps = {
    title?: string;
    users?: number;
    loading?: boolean;
    subTitle?: string;
    isConfigure?: boolean;
    handleOptionItemClick?: Function;
};

const NodeCard = ({
    title,
    users,
    subTitle,
    loading,
    handleOptionItemClick = () => {
        /* Default empty function */
    },
}: NodeCardProps) => {
    return (
        <>
            {loading ? (
                <SkeletonRoundedCard
                    variant="rectangular"
                    width={204}
                    height={206}
                />
            ) : (
                <Container>
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
                        <Grid item xs={2} m="4px 0px">
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
                        </Grid>
                        <Grid item xs={12} sx={{ ...IconStyle }}>
                            <NodeImg />
                        </Grid>
                        <Grid item xs={12}>
                            <Divider sx={{ m: "18px -18px 4px -18px" }} />
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
                                <UsersIcon
                                    fontSize="small"
                                    width="16px"
                                    height="16px"
                                />
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
                </Container>
            )}
        </>
    );
};

export default NodeCard;
