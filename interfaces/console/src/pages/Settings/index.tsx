import {
    UserSettings,
    AlertSettings,
    ConsoleSettings,
    NetworkSettings,
    LoadingWrapper,
} from "../../components";
import { RoundedCard } from "../../styles";
import {
    Divider,
    MenuList,
    MenuItem,
    Typography,
    Card,
    CardActions,
    CardContent,
    IconButton,
    Box,
    Grid,
} from "@mui/material";
import {
    useRecoilState,
    useRecoilValue,
    useResetRecoilState,
    useSetRecoilState,
} from "recoil";
import { colors } from "../../theme";
import { routes } from "../../router/config";
import { useHistory } from "react-router-dom";
import { SettingsMenuTypes } from "../../types";
import CloseIcon from "@mui/icons-material/Close";
import React, { useEffect, useState } from "react";
import { APP_VERSION, COPY_RIGHTS, SETTING_MENU } from "../../constants";
import { user, pageName, isDarkmode, isSkeltonLoading } from "../../recoil";

interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
}

const SettingMenuItem = ({ label, isSelected, handleItemClick }: any) => (
    <MenuItem
        selected={isSelected}
        onClick={handleItemClick}
        sx={{ color: `${label === "Log out" ? colors.error : "primary"}` }}
    >
        {label}
    </MenuItem>
);

const TabPanel = ({ children, index, value }: TabPanelProps) => {
    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            style={{
                width: "100%",
                height: "540px",
                overflowY: "scroll",
                overflowX: "hidden",
            }}
            id={`currentTab-indexpanel-${index}`}
            aria-labelledby={`currentTab-index-${index}`}
        >
            {value === index && children}
        </div>
    );
};

const Settings = () => {
    const history = useHistory();
    const [menuId, setMenuId] = useState(1);
    const _isDarkMod = useRecoilValue(isDarkmode);
    const setPage = useSetRecoilState(pageName);
    const resetPageName = useResetRecoilState(pageName);
    const resetData = useResetRecoilState(user);
    const [skeltonLoading, setSkeltonLoading] =
        useRecoilState(isSkeltonLoading);

    useEffect(() => {
        return () => setPage("Home");
    }, []);

    const handleItemClick = (id: number) => setMenuId(id);

    const handleLogout = () => {
        resetData();
        resetPageName();
        setSkeltonLoading(true);
        window.location.replace(`${process.env.REACT_APP_AUTH_URL}/logout`);
    };

    const handleCancel = () => {
        history.push(routes.Root.path);
    };

    return (
        <Card
            sx={{
                height: "100%",
                display: "flex",
                overflow: "scroll",
                flexDirection: "column",
                p: { xs: 0, md: "40px 85px 20px 85px" },
                backgroundColor: _isDarkMod ? colors.black : colors.solitude,
                justifyContent: "space-between",
            }}
        >
            <CardContent>
                <Grid container spacing={3} height="80%">
                    <Grid item xs={12} md={3}>
                        <LoadingWrapper height={237} isLoading={skeltonLoading}>
                            <RoundedCard
                                sx={{
                                    p: { xs: 1, md: 2 },
                                    height: "fit-content",
                                }}
                            >
                                <MenuList>
                                    {SETTING_MENU.map(
                                        ({ id, title }: SettingsMenuTypes) => (
                                            <SettingMenuItem
                                                key={id}
                                                label={title}
                                                isSelected={id === menuId}
                                                handleItemClick={() =>
                                                    handleItemClick(id)
                                                }
                                            />
                                        )
                                    )}
                                    <Divider
                                        sx={{
                                            mt: "18px !important",
                                            mb: "18px !important",
                                        }}
                                    />
                                    <SettingMenuItem
                                        label={"Log out"}
                                        isSelected={false}
                                        handleItemClick={handleLogout}
                                    />
                                </MenuList>
                            </RoundedCard>
                        </LoadingWrapper>
                    </Grid>
                    <Grid
                        item
                        xs={12}
                        md={9}
                        sx={{ mb: "18px", height: "90%" }}
                    >
                        <LoadingWrapper height={364} isLoading={skeltonLoading}>
                            <Card
                                sx={{
                                    px: { xs: 3, md: 4 },
                                    py: { xs: 2, md: 3 },
                                    borderRadius: "10px",
                                    boxShadow:
                                        "2px 2px 6px rgba(0, 0, 0, 0.05)",
                                }}
                            >
                                <CardContent
                                    sx={{
                                        p: 0,
                                        position: "relative",
                                    }}
                                >
                                    <Grid container xs={12}>
                                        <TabPanel index={1} value={menuId}>
                                            <UserSettings />
                                        </TabPanel>
                                        <TabPanel value={menuId} index={2}>
                                            <NetworkSettings />
                                        </TabPanel>
                                        <TabPanel value={menuId} index={3}>
                                            <AlertSettings />
                                        </TabPanel>
                                        <TabPanel value={menuId} index={4}>
                                            <ConsoleSettings />
                                        </TabPanel>
                                    </Grid>

                                    <Box
                                        component="div"
                                        sx={{
                                            top: 0,
                                            right: 0,
                                            position: "absolute",
                                        }}
                                    >
                                        <IconButton
                                            sx={{ p: 0 }}
                                            onClick={handleCancel}
                                        >
                                            <CloseIcon />
                                        </IconButton>
                                    </Box>
                                </CardContent>
                            </Card>
                        </LoadingWrapper>
                    </Grid>
                </Grid>
            </CardContent>
            <CardActions sx={{ alignItems: "center", display: "contents" }}>
                {!skeltonLoading && (
                    <Typography
                        variant={"caption"}
                        color="textSecondary"
                        sx={{
                            display: "block",
                            textAlign: "center",
                        }}
                    >
                        {`${APP_VERSION}`} <br /> {`${COPY_RIGHTS}`}
                    </Typography>
                )}
            </CardActions>
        </Card>
    );
};

export default Settings;
