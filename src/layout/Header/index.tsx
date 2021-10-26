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
import { MoreVert } from "@mui/icons-material";
import { useSetRecoilState } from "recoil";
import { isLoginAtom } from "../../recoil";
import { DRAWER_WIDTH, HEADER_MENU } from "../../constants";
import LogoutIcon from "@mui/icons-material/Logout";
import PersonOutlineIcon from "@mui/icons-material/PersonOutline";
import { HeaderMenuItemType } from "../../types";
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

    const menuId = "primary-search-account-menu";
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
                color="transparent"
                elevation={0}
            >
                <Toolbar style={{ flexGrow: 1, padding: "0px" }}>
                    <Typography variant="h6" noWrap component="div">
                        {pageName}
                    </Typography>

                    <Box sx={{ flexGrow: 1 }} />
                    <Box sx={{ display: { xs: "none", md: "flex" } }}>
                        {HEADER_MENU.map(({ id, Icon }: HeaderMenuItemType) => (
                            <IconButton
                                key={id}
                                size="large"
                                aria-label="show 4 new mails"
                                color="inherit"
                            >
                                <Icon />
                            </IconButton>
                        ))}
                        );
                    </Box>
                    <Box sx={{ display: { xs: "flex", md: "none" } }}>
                        <IconButton
                            size="large"
                            aria-label="show more"
                            aria-controls={mobileMenuId}
                            aria-haspopup="true"
                            onClick={handleMobileMenuOpen}
                            color="inherit"
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
