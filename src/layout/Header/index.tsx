import {
    Divider,
    Menu,
    AppBar,
    MenuItem,
    Typography,
    IconButton,
    Toolbar,
    Box,
} from "@mui/material";
import React from "react";
import config from "../../config";
import { colors } from "../../theme";
import { HEADER_MENU } from "../../constants";
import { MoreVert } from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import { HeaderMenuItemType } from "../../types";
import { SkeletonRoundedCard } from "../../styles";

type HeaderProps = {
    pageName: string;
    isLoading: boolean;
    handleDrawerToggle: Function;
};

const Header = ({ pageName, handleDrawerToggle, isLoading }: HeaderProps) => {
    const showDivider =
        pageName !== "Billing" && pageName !== "User" ? true : false;
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

    const isMenuOpen = Boolean(anchorEl);

    const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
    };

    const handleMobileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleLogout = () => {
        handleMenuClose();
        window.close();
        window.location.replace(`${config.REACT_APP_AUTH_URL}logout`);
    };

    const handleHeaderMenu = (
        e: React.MouseEvent<HTMLElement>,
        name: string
    ) => {
        switch (name) {
            case "Setting":
                //GOTO Settings page
                handleMenuClose();
                break;
            case "Notification":
                //GOTO Notification page
                handleMenuClose();
                break;
            case "Account":
                handleProfileMenuOpen(e);
                break;
        }
    };

    const menuId = "account-popup-menu";
    const renderMenu = (
        <Menu
            anchorEl={anchorEl}
            anchorOrigin={{
                vertical: "top",
                horizontal: "right",
            }}
            id={menuId}
            keepMounted
            transformOrigin={{
                vertical: "top",
                horizontal: "right",
            }}
            open={isMenuOpen}
            onClose={handleMenuClose}
        >
            <MenuItem onClick={handleMenuClose}>Profile</MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout}>Logout</MenuItem>
        </Menu>
    );

    return (
        <Box>
            <AppBar
                elevation={0}
                position="relative"
                sx={{
                    backgroundColor: colors.solitude,
                    width: { sm: "100%" },
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
                    {isLoading ? (
                        <SkeletonRoundedCard
                            variant="rectangular"
                            height={30}
                            width={82}
                        />
                    ) : (
                        <Typography
                            noWrap
                            variant="h5"
                            component="div"
                            color="black"
                        >
                            {pageName}
                        </Typography>
                    )}
                    <Box sx={{ flexGrow: 1 }} />
                    {isLoading ? (
                        <SkeletonRoundedCard
                            variant="rectangular"
                            height={30}
                            width={120}
                        />
                    ) : (
                        <Box sx={{ display: { xs: "none", md: "flex" } }}>
                            {HEADER_MENU.map(
                                ({ id, Icon, title }: HeaderMenuItemType) => (
                                    <IconButton
                                        key={id}
                                        size="medium"
                                        color="inherit"
                                        sx={{ padding: "0px 18px" }}
                                        onClick={e =>
                                            handleHeaderMenu(e, title)
                                        }
                                    >
                                        <Icon />
                                    </IconButton>
                                )
                            )}
                        </Box>
                    )}
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
                {showDivider && <Divider sx={{ m: "0px" }} />}
            </AppBar>
            {renderMenu}
        </Box>
    );
};

export default Header;
