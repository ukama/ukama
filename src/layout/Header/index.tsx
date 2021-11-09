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
import { colors } from "../../theme";
import { useSetRecoilState } from "recoil";
import { isLoginAtom } from "../../recoil";
import { HEADER_MENU } from "../../constants";
import { MoreVert } from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import { HeaderMenuItemType } from "../../types";

type HeaderProps = {
    pageName: string;
    handleDrawerToggle: Function;
};

const Header = ({ pageName, handleDrawerToggle }: HeaderProps) => {
    const setIsLogin = useSetRecoilState(isLoginAtom);
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
                <Toolbar style={{ padding: "0px" }}>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={() => handleDrawerToggle()}
                        sx={{ mr: 2, display: { sm: "none" } }}
                    >
                        <MenuIcon color={"primary"} />
                    </IconButton>
                    <Typography
                        noWrap
                        variant="h5"
                        component="div"
                        color="black"
                    >
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
                            aria-controls={menuId}
                            onClick={handleMobileMenuOpen}
                        >
                            <MoreVert color={"primary"} />
                        </IconButton>
                    </Box>
                </Toolbar>
                <Divider sx={{ m: "0px" }} />
            </AppBar>
            {renderMenu}
        </Box>
    );
};

export default Header;
