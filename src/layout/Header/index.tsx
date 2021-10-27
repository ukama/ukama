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
                            <IconButton key={id} size="large" color="inherit">
                                <Icon />
                            </IconButton>
                        ))}
                    </Box>
                </Toolbar>
                <Divider />
            </AppBar>
        </Box>
    );
};

export default Header;
