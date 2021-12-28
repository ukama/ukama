import {
    Box,
    Menu,
    AppBar,
    Divider,
    Popover,
    Toolbar,
    MenuItem,
    IconButton,
    Typography,
} from "@mui/material";
import {
    useGetAlertsQuery,
    GetLatestAlertsDocument,
    GetLatestAlertsSubscription,
} from "../../generated";
import { colors } from "../../theme";
import { RoundedCard } from "../../styles";
import { useHistory } from "react-router-dom";
import { MoreVert } from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import { cloneDeep } from "@apollo/client/utilities";
import { Alerts, LoadingWrapper } from "../../components";
import React, { useEffect, useRef, useState } from "react";
import { AccountIcon, NotificationIcon, SettingsIcon } from "../../assets/svg";

type HeaderProps = {
    pageName: string;
    isLoading: boolean;
    handlePageChange: Function;
    handleDrawerToggle: Function;
};

const Header = ({
    pageName,
    handlePageChange,
    handleDrawerToggle,
    isLoading,
}: HeaderProps) => {
    const history = useHistory();
    const showDivider = pageName !== "Billing" ? true : false;
    const ref = useRef(null);
    const menuId = "account-popup-menu";
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const isMenuOpen = Boolean(anchorEl);
    const handleMenuClose = () => {
        setAnchorEl(null);
    };
    const handleMobileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const [notificationAnchorEl, setNotificationAnchorEl] =
        useState<HTMLButtonElement | null>(null);
    const handleNotificationClick = () => {
        setNotificationAnchorEl(ref.current);
    };
    const handleNotificationClose = () => {
        setNotificationAnchorEl(null);
    };
    const open = Boolean(notificationAnchorEl);
    const notificationAnchorElId = open ? "simple-popover" : undefined;

    const handleSettingsClick = () => {
        handleMenuClose();
        handlePageChange("Settings");
        history.push("/settings");
    };

    const { data: alertsInfoRes, subscribeToMore: subscribeToLatestAlerts } =
        useGetAlertsQuery({
            variables: {
                data: {
                    pageNo: 1,
                    pageSize: 50,
                },
            },
        });

    useEffect(() => {
        if (alertsInfoRes) {
            subscribeToLatestAlerts<GetLatestAlertsSubscription>({
                document: GetLatestAlertsDocument,
                updateQuery: (prev, { subscriptionData }) => {
                    let data = cloneDeep(prev);
                    const latestAlert = subscriptionData.data.getAlerts;
                    if (latestAlert.__typename === "AlertDto")
                        data.getAlerts.alerts = [
                            latestAlert,
                            ...data.getAlerts.alerts,
                        ];
                    return data;
                },
            });
        }
    }, [alertsInfoRes]);

    const renderMenu = (
        <Menu
            id={menuId}
            keepMounted
            anchorEl={anchorEl}
            anchorOrigin={{
                vertical: "top",
                horizontal: "right",
            }}
            transformOrigin={{
                vertical: "top",
                horizontal: "right",
            }}
            open={isMenuOpen}
            onClose={handleMenuClose}
        >
            <MenuItem onClick={handleSettingsClick}>Settings</MenuItem>
            <Divider />
            <MenuItem onClick={handleMenuClose}>Notifications</MenuItem>
        </Menu>
    );

    return (
        <Box>
            <Popover
                open={open}
                id={notificationAnchorElId}
                anchorEl={notificationAnchorEl}
                onClose={handleNotificationClose}
                anchorOrigin={{
                    vertical: "bottom",
                    horizontal: "right",
                }}
                transformOrigin={{
                    vertical: "top",
                    horizontal: "center",
                }}
            >
                <RoundedCard>
                    <Typography variant="h6" sx={{ mb: "14px" }}>
                        Alerts
                    </Typography>
                    <Alerts alertOptions={alertsInfoRes?.getAlerts?.alerts} />
                </RoundedCard>
            </Popover>
            <AppBar
                elevation={0}
                position="relative"
                sx={{
                    backgroundColor: colors.solitude,
                }}
            >
                <Toolbar sx={{ padding: "33px 0px 12px 0px !important" }}>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={() => handleDrawerToggle()}
                        sx={{ mr: 2, display: { sm: "none" } }}
                    >
                        <MenuIcon color={"primary"} />
                    </IconButton>

                    <LoadingWrapper
                        height={30}
                        width={82}
                        isLoading={isLoading}
                    >
                        <Typography
                            noWrap
                            variant="h5"
                            component="div"
                            color="black"
                        >
                            {pageName}
                        </Typography>
                    </LoadingWrapper>

                    <Box sx={{ flexGrow: 1 }} />

                    <LoadingWrapper
                        height={30}
                        width={120}
                        isLoading={isLoading}
                    >
                        <Box sx={{ display: { xs: "none", md: "flex" } }}>
                            <IconButton
                                size="medium"
                                color="inherit"
                                sx={{ padding: "0px 18px" }}
                                onClick={handleSettingsClick}
                            >
                                <SettingsIcon />
                            </IconButton>
                            <IconButton
                                size="medium"
                                color="inherit"
                                sx={{ padding: "0px 18px" }}
                                onClick={handleNotificationClick}
                            >
                                <NotificationIcon />
                            </IconButton>
                            <IconButton
                                size="medium"
                                color="inherit"
                                sx={{ padding: "0px 18px" }}
                            >
                                <AccountIcon />
                            </IconButton>
                        </Box>
                    </LoadingWrapper>

                    <Box sx={{ display: { xs: "flex", md: "none" } }}>
                        <IconButton
                            size="large"
                            color="inherit"
                            aria-controls={menuId}
                            onClick={handleMobileMenuOpen}
                        >
                            <MoreVert color={"primary"} />
                        </IconButton>
                    </Box>
                </Toolbar>
                {showDivider && <Divider ref={ref} sx={{ m: "0px" }} />}
            </AppBar>
            {renderMenu}
        </Box>
    );
};

export default Header;
