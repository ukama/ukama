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
import { UpgradeNavFooter } from "../../components";
import { SkeletonRoundedCard } from "../../styles";

const useStyles = makeStyles(() => ({
    listItem: {
        opacity: 1,
        height: "40px",
        marginTop: "6px",
        padding: "8px 12px",
        borderRadius: "4px",
        fontFamily: "Work Sans",
        backgroundColor: colors.white,
    },
    listItemText: {
        color: colors.black,
    },
    listItemDone: {
        opacity: 1,
        height: "40px",
        marginTop: "8px",
        padding: "8px 12px",
        borderRadius: "4px",
        color: `${colors.white} !important`,
        backgroundColor: `${colors.primary} !important`,
    },
    listItemDoneText: {
        color: colors.white,
    },
}));

type SidebarProps = {
    page: string;
    isOpen: boolean;
    isLoading: boolean;
    handlePageChange: Function;
    handleDrawerToggle: Function;
};

const Sidebar = (
    {
        page,
        isOpen,
        isLoading,
        handlePageChange,
        handleDrawerToggle,
    }: SidebarProps,
    props: any
) => {
    const { window } = props;
    const classes = useStyles();
    const history = useHistory();

    const MenuList = (list: any) => (
        <List sx={{ padding: "8px 20px" }}>
            {list.map(({ id, title, Icon, route }: MenuItemType) => (
                <ListItem
                    button
                    key={id}
                    href={route}
                    onClick={() => {
                        handlePageChange(title);
                        history.push(route);
                    }}
                    selected={title === page}
                    className={
                        title === page ? classes.listItemDone : classes.listItem
                    }
                >
                    <ListItemIcon sx={{ minWidth: "44px" }}>
                        <Icon
                            color={
                                title === page ? colors.white : colors.vulcan
                            }
                        />
                    </ListItemIcon>
                    <ListItemText>
                        <Typography
                            variant="body1"
                            fontWeight={title === page ? "bold" : "normal"}
                            className={
                                title === page
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
            <Toolbar sx={{ padding: "33px 0px 12px 0px" }}>
                <Logo width={"100%"} height={"36px"} />
            </Toolbar>
            {MenuList(SIDEBAR_MENU1)}
            <div
                style={{
                    width: "100%",
                    display: "flex",
                    justifyContent: "center",
                }}
            >
                <Divider sx={{ width: 160, mt: "8px", mb: "0px !important" }} />
            </div>
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
                width: { xs: 0, sm: DRAWER_WIDTH },
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
            {isLoading ? (
                <SkeletonRoundedCard
                    variant="rectangular"
                    sx={{ borderRadius: 0 }}
                />
            ) : (
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
            )}
        </Box>
    );
};

export default Sidebar;
