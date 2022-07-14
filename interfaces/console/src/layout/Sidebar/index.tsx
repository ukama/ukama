import {
    Box,
    List,
    Paper,
    Stack,
    Drawer,
    Toolbar,
    ListItem,
    Typography,
    Chip,
    ListItemIcon,
    ListItemText,
} from "@mui/material";
import React, { useEffect } from "react";
import { colors } from "../../theme";
import { useRecoilValue } from "recoil";
import { makeStyles } from "@mui/styles";
import { isDarkmode } from "../../recoil";
import { MenuItemType } from "../../types";
import { useHistory } from "react-router-dom";
import { LoadingWrapper } from "../../components";
import { DRAWER_WIDTH, SIDEBAR_MENU1 } from "../../constants";
import { useGetNodesByOrgLazyQuery } from "../../generated";

const Logo = React.lazy(() =>
    import("../../assets/svg").then(module => ({
        default: module.Logo,
    }))
);

const useStyles = makeStyles(() => ({
    listItem: {
        opacity: 1,
        height: "40px",
        marginTop: "6px",
        padding: "8px 12px",
        borderRadius: "4px",
    },
    listItemText: {},
    listItemDone: {
        opacity: 1,
        height: "40px",
        marginTop: "8px",
        padding: "8px 12px",
        borderRadius: "4px",
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
    const _isDarkMod = useRecoilValue(isDarkmode);
    const [getNodesByOrg, { data: nodesRes, loading: nodesLoading }] =
        useGetNodesByOrgLazyQuery({
            fetchPolicy: "cache-and-network",
        });
    useEffect(() => {
        getNodesByOrg();
    }, []);
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
                            fontSize="medium"
                            sx={{
                                fill: _isDarkMod ? colors.white : colors.vulcan,
                            }}
                        />
                    </ListItemIcon>
                    <ListItemText>
                        <Stack
                            direction="row"
                            justifyContent="space-between"
                            alignItems="center"
                            spacing={1}
                        >
                            <Typography
                                variant="body1"
                                fontWeight={title === page ? 600 : "normal"}
                                className={
                                    title !== page ? classes.listItemText : ""
                                }
                            >
                                {title}
                            </Typography>
                            {title == "Nodes" &&
                                nodesRes?.getNodesByOrg.totalNodes && (
                                    <LoadingWrapper
                                        isLoading={nodesLoading}
                                        width={"120px"}
                                        height={"20px"}
                                    >
                                        <Chip
                                            label={
                                                nodesRes?.getNodesByOrg
                                                    .totalNodes
                                            }
                                            size="small"
                                            color="primary"
                                        />
                                    </LoadingWrapper>
                                )}
                        </Stack>
                    </ListItemText>
                </ListItem>
            ))}
        </List>
    );

    const drawer = (
        <Paper
            sx={{
                height: "100%",
                overflowY: "auto",
            }}
        >
            <Toolbar sx={{ padding: "33px 0px 12px 0px" }}>
                <Logo
                    width={"100%"}
                    height={"36px"}
                    color={_isDarkMod ? colors.white : colors.primaryMain}
                />
            </Toolbar>
            <Stack
                sx={{
                    display: "flex",
                    minHeight: "200px",
                    height: `calc(100% - 400px)`,
                }}
            >
                {MenuList(SIDEBAR_MENU1)}
            </Stack>
        </Paper>
    );

    const container =
        window !== undefined ? () => window().document.body : undefined;

    return (
        <Box
            component="nav"
            sx={{
                flexShrink: { sm: 0 },
                width: { xs: 0, sm: DRAWER_WIDTH },
                boxShadow: "6px 0px 18px rgba(0, 0, 0, 0.06)",
            }}
            aria-label="mailbox folders"
        >
            <LoadingWrapper isLoading={isLoading}>
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
            </LoadingWrapper>
        </Box>
    );
};

export default Sidebar;
