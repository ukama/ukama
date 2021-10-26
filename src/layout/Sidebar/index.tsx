import Box from "@mui/material/Box";
import List from "@mui/material/List";
import { Logo } from "../../assets/svg";
import Drawer from "@mui/material/Drawer";
import { MenuItemType } from "../../types";
import Divider from "@mui/material/Divider";
import Toolbar from "@mui/material/Toolbar";
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import { DRAWER_WIDTH, SIDEBAR_MENU1, SIDEBAR_MENU2 } from "../../constants";
import { useHistory } from "react-router-dom";
import { Dispatch, SetStateAction } from "react";
import { UpgradeNavFooter } from "../../components";
import { colors } from "../../theme";
import { makeStyles } from "@mui/styles";

const useStyles = makeStyles(() => ({
    listItem: {
        opacity: 1,
        borderRadius: "4px",
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
                    selected={route === path}
                    onClick={() => {
                        setPath(route);
                        history.push(route);
                    }}
                    className={
                        route === path ? classes.listItemDone : classes.listItem
                    }
                >
                    <ListItemIcon>
                        <Icon
                            color={
                                route === path ? colors.white : colors.vulcan
                            }
                        />
                    </ListItemIcon>
                    <ListItemText
                        primary={title}
                        className={
                            route === path
                                ? classes.listItemDoneText
                                : classes.listItemText
                        }
                    />
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
                <Logo width={"100%"} height={"40px"} />
            </Toolbar>
            {MenuList(SIDEBAR_MENU1)}
            <Divider />
            {MenuList(SIDEBAR_MENU2)}

            <div
                style={{
                    position: "inherit",
                    top: "calc(100vh - 64%)",
                }}
            >
                <UpgradeNavFooter />
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
