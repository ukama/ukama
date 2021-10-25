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
    const history = useHistory();
    const drawer = (
        <div>
            <Toolbar>
                <Logo width={"100%"} height={"40px"} />
            </Toolbar>

            <List sx={{ padding: "8px 14px" }}>
                {SIDEBAR_MENU1.map(
                    ({ id, title, Icon, route }: MenuItemType) => (
                        <ListItem
                            button
                            key={id}
                            href={route}
                            selected={route === path}
                            onClick={() => {
                                setPath(route);
                                history.push(route);
                            }}
                            sx={{
                                borderRadius: "4px",
                                opacity: 1,
                            }}
                        >
                            <ListItemIcon>
                                <Icon style={{ fill: "white" }} />
                            </ListItemIcon>
                            <ListItemText primary={title} />
                        </ListItem>
                    )
                )}
            </List>
            <Divider />
            <List>
                {SIDEBAR_MENU2.map(
                    ({ id, title, Icon, route }: MenuItemType) => (
                        <ListItem button key={id} href={route}>
                            <ListItemIcon>
                                <Icon />
                            </ListItemIcon>
                            <ListItemText primary={title} />
                        </ListItem>
                    )
                )}
            </List>
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
                variant="permanent"
                sx={{
                    display: { xs: "none", sm: "block" },
                    "& .MuiDrawer-paper": {
                        boxSizing: "border-box",
                        width: DRAWER_WIDTH,
                    },
                }}
                open
            >
                {drawer}
            </Drawer>
        </Box>
    );
};

export default Sidebar;
