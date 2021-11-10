import {
    Box,
    List,
    Drawer,
    Divider,
    Toolbar,
    ListItem,
    ListItemIcon,
    ListItemText,
    Typography,
} from "@mui/material";
import {
    APP_VERSION,
    COPY_RIGHTS,
    DRAWER_WIDTH,
    SIDEBAR_MENU1,
    SIDEBAR_MENU2,
} from "../../constants";
import { colors } from "../../theme";
import { Logo } from "../../assets/svg";
import { makeStyles } from "@mui/styles";
import { MenuItemType } from "../../types";
import { useHistory } from "react-router-dom";
import { Dispatch, SetStateAction } from "react";
import { UpgradeNavFooter } from "../../components";

const useStyles = makeStyles(() => ({
    listItem: {
        opacity: 1,
        borderRadius: "4px",
        fontFamily: "Work Sans",
        backgroundColor: colors.white,
    },
    listItemText: {
        color: colors.black,
    },
    listItemDone: {
        opacity: 1,
        borderRadius: "4px",
        color: `${colors.white} !important`,
        backgroundColor: `${colors.primary} !important`,
    },
    listItemDoneText: {
        color: colors.white,
    },
}));

type SidebarProps = {
    path: string;
    isOpen: boolean;
    handleDrawerToggle: Function;
    setPath: Dispatch<SetStateAction<string>>;
};

const Sidebar = (
    { isOpen, handleDrawerToggle, path, setPath }: SidebarProps,
    props: any
) => {
    const { window } = props;
    const classes = useStyles();
    const history = useHistory();

    const MenuList = (list: any) => (
        <List sx={{ padding: "8px 14px" }}>
            {list.map(({ id, title, Icon, route }: MenuItemType) => (
                <ListItem
                    button
                    key={id}
                    href={route}
                    onClick={() => {
                        setPath(title);
                        history.push(route);
                    }}
                    selected={title === path}
                    className={
                        title === path ? classes.listItemDone : classes.listItem
                    }
                >
                    <ListItemIcon sx={{ minWidth: "44px" }}>
                        <Icon
                            color={
                                title === path ? colors.white : colors.vulcan
                            }
                        />
                    </ListItemIcon>
                    <ListItemText>
                        <Typography
                            variant="body1"
                            fontWeight={title === path ? "bold" : "normal"}
                            className={
                                title === path
                                    ? classes.listItemDoneText
                                    : classes.listItemText
                            }
                        >
                            {title}
                        </Typography>
                    </ListItemText>
                </ListItem>
            ))}
        </List>
    );

    const drawer = (
        <div
            style={{
                height: "100%",
                overflowY: "auto",
                position: "relative",
            }}
        >
            <Toolbar sx={{ padding: "16px 0px" }}>
                <Logo width={"100%"} height={"44px"} />
            </Toolbar>
            {MenuList(SIDEBAR_MENU1)}
            <Divider sx={{ m: "16px 14px 0px 14px" }} />
            {MenuList(SIDEBAR_MENU2)}

            <div
                style={{
                    position: "absolute",
                    bottom: "10px",
                }}
            >
                <UpgradeNavFooter />
                <Typography
                    variant={"caption"}
                    sx={{
                        display: "block",
                        textAlign: "center",
                        color: colors.empress,
                    }}
                >
                    {`${APP_VERSION}`} <br /> {`${COPY_RIGHTS}`}
                </Typography>
            </div>
        </div>
    );

    const container =
        window !== undefined ? () => window().document.body : undefined;

    return (
        <Box
            component="nav"
            sx={{
                width: { sm: DRAWER_WIDTH },
                flexShrink: { sm: 0 },
                boxShadow: "6px 0px 18px rgba(0, 0, 0, 0.06)",
            }}
            aria-label="mailbox folders"
        >
            <Drawer
                open={isOpen}
                variant="temporary"
                container={container}
                onClose={() => handleDrawerToggle()}
                ModalProps={{
                    keepMounted: true,
                }}
                sx={{
                    display: { xs: "block", sm: "none" },
                    "& .MuiDrawer-paper": {
                        boxSizing: "border-box",
                        width: DRAWER_WIDTH,
                    },
                    borderRight: "0px",
                }}
            >
                {drawer}
            </Drawer>
            <Drawer
                open
                variant="permanent"
                sx={{
                    display: { xs: "none", sm: "block" },
                    "& .MuiDrawer-paper": {
                        boxSizing: "border-box",
                        width: DRAWER_WIDTH,
                    },
                }}
            >
                {drawer}
            </Drawer>
        </Box>
    );
};

export default Sidebar;
