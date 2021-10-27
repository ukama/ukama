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
import { useSetRecoilState } from "recoil";
import { isLoginAtom } from "../../recoil";
import { MoreVert } from "@mui/icons-material";
import { HeaderMenuItemType } from "../../types";
import LogoutIcon from "@mui/icons-material/Logout";
import { DRAWER_WIDTH, HEADER_MENU } from "../../constants";
import PersonOutlineIcon from "@mui/icons-material/PersonOutline";

type HeaderProps = {
    pageName: string;
};

const Header = ({ pageName }: HeaderProps) => {
    const setIsLogin = useSetRecoilState(isLoginAtom);
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const [mobileMoreAnchorEl, setMobileMoreAnchorEl] =
        React.useState<null | HTMLElement>(null);

    const isMenuOpen = Boolean(anchorEl);
    const isMobileMenuOpen = Boolean(mobileMoreAnchorEl);

    const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMobileMenuClose = () => {
        setMobileMoreAnchorEl(null);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
        handleMobileMenuClose();
    };

    const handleMobileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
        setMobileMoreAnchorEl(event.currentTarget);
    };

    const handleLogout = () => {
        handleMenuClose();
        setIsLogin(false);
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

    const mobileMenuId = "primary-search-account-menu-mobile";
    const renderMobileMenu = (
        <Menu
            anchorEl={mobileMoreAnchorEl}
            anchorOrigin={{
                vertical: "top",
                horizontal: "right",
            }}
            id={mobileMenuId}
            keepMounted
            transformOrigin={{
                vertical: "top",
                horizontal: "right",
            }}
            open={isMobileMenuOpen}
            onClose={handleMobileMenuClose}
        >
            <MenuItem>
                <IconButton size="large" color="inherit">
                    <PersonOutlineIcon />
                </IconButton>
                <p>Profile</p>
            </MenuItem>

            <MenuItem onClick={handleProfileMenuOpen}>
                <IconButton size="large" color="inherit">
                    <LogoutIcon />
                </IconButton>
                <p>Logout</p>
            </MenuItem>
        </Menu>
    );

    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar
                position="fixed"
                sx={{
                    padding: "4px 30px",
                    ml: `${DRAWER_WIDTH}px`,
                    width: `calc(100% - ${DRAWER_WIDTH}px)`,
                }}
                elevation={0}
                color="transparent"
            >
                <Toolbar style={{ flexGrow: 1, padding: "0px" }}>
                    <Typography variant="h6" noWrap component="div">
                        {pageName}
                    </Typography>

                    <Box sx={{ flexGrow: 1 }} />
                    <Box sx={{ display: { xs: "none", md: "flex" } }}>
                        {HEADER_MENU.map(
                            ({ id, Icon, title }: HeaderMenuItemType) => (
                                <IconButton
                                    key={id}
                                    size="large"
                                    color="inherit"
                                    onClick={e => handleHeaderMenu(e, title)}
                                >
                                    <Icon />
                                </IconButton>
                            )
                        )}
                    </Box>
                    <Box sx={{ display: { xs: "flex", md: "none" } }}>
                        <IconButton
                            size="large"
                            color="inherit"
                            aria-controls={mobileMenuId}
                            onClick={handleMobileMenuOpen}
                        >
                            <MoreVert />
                        </IconButton>
                    </Box>
                </Toolbar>
                <Divider />
            </AppBar>

            {renderMobileMenu}
            {renderMenu}
        </Box>
    );
};

export default Header;
