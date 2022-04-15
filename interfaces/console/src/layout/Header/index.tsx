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
    Badge,
    Stack,
    Button,
} from "@mui/material";
import {
    GetLatestAlertsDocument,
    GetLatestAlertsSubscription,
    useGetAlertsQuery,
} from "../../generated";
import {
    MoreVert,
    Settings,
    Notifications,
    AccountCircle,
} from "@mui/icons-material";
import { colors } from "../../theme";
import { RoundedCard } from "../../styles";
import { routes } from "../../router/config";
import { useHistory } from "react-router-dom";
import MenuIcon from "@mui/icons-material/Menu";
import { cloneDeep } from "@apollo/client/utilities";
import React, { useEffect, useRef, useState } from "react";
import ExitToAppOutlined from "@mui/icons-material/ExitToAppOutlined";
import { Alerts, DarkModToggle, LoadingWrapper } from "../../components";
import { useRecoilValue, useResetRecoilState, useSetRecoilState } from "recoil";
import { isSkeltonLoading, user, pageName } from "../../recoil";
import { useCookies } from "react-cookie";

type HeaderProps = {
    pageName: string;
    isLoading: boolean;
    handlePageChange: Function;
    handleDrawerToggle: Function;
};

const Header = ({
    pageName: _pageName,
    handlePageChange,
    handleDrawerToggle,
    isLoading,
}: HeaderProps) => {
    const history = useHistory();
    const showDivider = _pageName !== "Billing" ? true : false;
    const ref = useRef(null);
    const menuId = "account-popup-menu";
    const _user = useRecoilValue(user);
    const resetPageName = useResetRecoilState(pageName);
    const resetData = useResetRecoilState(user);
    const [, , removeCookie] = useCookies(["orgId"]);
    const setSkeltonLoading = useSetRecoilState(isSkeltonLoading);
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
    const [userAnchorEl, setUserAnchorEl] = useState<HTMLButtonElement | null>(
        null
    );
    const handleNotificationClick = () => {
        setNotificationAnchorEl(ref.current);
    };
    const handleUserClick = () => {
        setUserAnchorEl(ref.current);
    };
    const handleNotificationClose = () => {
        setNotificationAnchorEl(null);
    };
    const handleUserClose = () => {
        setUserAnchorEl(null);
    };
    const open = Boolean(notificationAnchorEl);
    const openUserPopover = Boolean(userAnchorEl);
    const userAnchorElId = openUserPopover ? "user-popover" : undefined;
    const notificationAnchorElId = open ? "simple-popover" : undefined;

    const handleSettingsClick = () => {
        handleMenuClose();
        handlePageChange("Settings");
        history.push(routes.Settings.path);
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

    const alertSubscription = () =>
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

    useEffect(() => {
        let unsub = alertSubscription();
        return () => {
            unsub && unsub();
        };
    }, [alertsInfoRes]);

    const handleLogout = () => {
        removeCookie("orgId");
        handleUserClose();
        resetData();
        resetPageName();
        setSkeltonLoading(true);
        window.location.replace(`${process.env.REACT_APP_AUTH_URL}/logout`);
    };

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
        <Box component="div">
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
                PaperProps={{
                    style: {
                        background: "transparent",
                    },
                }}
            >
                <RoundedCard>
                    <Typography variant="h6" sx={{ mb: "14px" }}>
                        Alerts
                    </Typography>
                    <Alerts alertOptions={alertsInfoRes?.getAlerts?.alerts} />
                </RoundedCard>
            </Popover>
            <Popover
                open={openUserPopover}
                id={userAnchorElId}
                anchorEl={userAnchorEl}
                onClose={handleUserClose}
                anchorOrigin={{
                    vertical: "bottom",
                    horizontal: "right",
                }}
                transformOrigin={{
                    vertical: "top",
                    horizontal: "center",
                }}
                PaperProps={{
                    style: {
                        background: "transparent",
                    },
                }}
            >
                <RoundedCard sx={{ minWidth: "200px", p: 0 }}>
                    <Stack m={"12px 16px"}>
                        <Typography variant="body1">{_user.name}</Typography>
                        <Typography variant="caption" color={"textSecondary"}>
                            {_user.email}
                        </Typography>
                    </Stack>
                    <Divider sx={{ mt: "6px" }} />
                    <Button
                        onClick={handleLogout}
                        startIcon={<ExitToAppOutlined />}
                        sx={{
                            mb: "12px",
                            mx: "16px",
                            typography: "body1",
                            textTransform: "capitalize",
                            ":hover": {
                                backgroundColor: "white !important",
                                svg: {
                                    fill: colors.primaryDark,
                                },
                            },
                        }}
                    >
                        Sign out
                    </Button>
                </RoundedCard>
            </Popover>
            <AppBar
                elevation={0}
                position="relative"
                color="transparent"
                sx={{ boxShadow: "none !important" }}
            >
                <Toolbar sx={{ padding: "33px 0px 12px 0px !important" }}>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={() => handleDrawerToggle()}
                        sx={{ mr: 2, display: { sm: "none" } }}
                    >
                        <MenuIcon />
                    </IconButton>

                    <LoadingWrapper
                        height={30}
                        width={82}
                        isLoading={isLoading}
                    >
                        <Typography variant="h5">{_pageName}</Typography>
                    </LoadingWrapper>

                    <Box component="div" sx={{ flexGrow: 1 }} />

                    <LoadingWrapper
                        height={30}
                        width={120}
                        isLoading={isLoading}
                    >
                        <Stack
                            spacing={3}
                            direction="row"
                            sx={{
                                display: { xs: "none", md: "flex" },
                                justifyContent: "flex-end",
                            }}
                        >
                            <DarkModToggle />
                            <IconButton
                                size="small"
                                color="inherit"
                                aria-label="setting-btn"
                                onClick={handleSettingsClick}
                            >
                                <Settings />
                            </IconButton>
                            <IconButton
                                size="small"
                                color="inherit"
                                aria-label="notification-btn"
                                onClick={handleNotificationClick}
                                // aria-label={notificationsLabel(
                                //     alertsInfoRes?.getAlerts?.alerts.length
                                // )}
                            >
                                <Badge
                                    badgeContent={
                                        alertsInfoRes?.getAlerts?.alerts.length
                                    }
                                    sx={{
                                        "& .MuiBadge-badge": {
                                            color: "inherit",
                                            paddingLeft: "3px",
                                            paddingRight: "3px",
                                            backgroundColor:
                                                colors.secondaryMain,
                                        },
                                    }}
                                >
                                    <Notifications
                                        color={
                                            notificationAnchorEl
                                                ? "primary"
                                                : "inherit"
                                        }
                                    />
                                </Badge>
                            </IconButton>
                            <IconButton
                                size="small"
                                color="inherit"
                                aria-label="account-btn"
                                onClick={handleUserClick}
                            >
                                <AccountCircle />
                            </IconButton>
                        </Stack>
                    </LoadingWrapper>

                    <Box
                        component="div"
                        sx={{ display: { xs: "flex", md: "none" } }}
                    >
                        <IconButton
                            size="large"
                            color="inherit"
                            aria-controls={menuId}
                            onClick={handleMobileMenuOpen}
                        >
                            <MoreVert />
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
