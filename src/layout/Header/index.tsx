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
import { RoundedCard } from "../../styles";
import { routes } from "../../router/config";
import { useHistory } from "react-router-dom";
import MenuIcon from "@mui/icons-material/Menu";
import React, { useEffect, useRef, useState } from "react";
import { Alerts, DarkModToggle, LoadingWrapper } from "../../components";
import { colors } from "../../theme";
import { cloneDeep } from "@apollo/client/utilities";

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
                        <Typography variant="h5">{pageName}</Typography>
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
