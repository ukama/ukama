import {
    UserSettings,
    NodeSettings,
    AlertSettings,
    NetworkSettings,
    LoadingWrapper,
} from "../../components";
import {
    RoundedCard,
    VerticalContainer,
    HorizontalContainer,
} from "../../styles";
import {
    Box,
    Grid,
    Button,
    Divider,
    MenuList,
    MenuItem,
    Typography,
} from "@mui/material";
import { colors } from "../../theme";
import { useHistory } from "react-router-dom";
import { SettingsMenuTypes } from "../../types";
import React, { useEffect, useState } from "react";
import { isSkeltonLoading, organizationId, pageName } from "../../recoil";
import { useRecoilState, useResetRecoilState, useSetRecoilState } from "recoil";
import { APP_VERSION, COPY_RIGHTS, SETTING_MENU } from "../../constants";

interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
}

interface ActionButtonsProps {
    handleCancelAction: Function;
    handleSaveAction: Function;
}

const SettingMenuItem = ({ label, isSelected, handleItemClick }: any) => (
    <MenuItem
        selected={isSelected}
        onClick={handleItemClick}
        sx={{ color: `${label === "Log out" ? colors.red900 : "black"}` }}
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
                height: "100%",
                position: "relative",
            }}
            id={`currentTab-indexpanel-${index}`}
            aria-labelledby={`currentTab-index-${index}`}
        >
            {value === index && children}
        </div>
    );
};

const ActionButtons = ({
    handleCancelAction,
    handleSaveAction,
}: ActionButtonsProps) => (
    <VerticalContainer
        sx={{
            right: 0,
            bottom: 0,
            position: "absolute",
        }}
    >
        <Divider sx={{ width: "100%" }} />
        <HorizontalContainer
            sx={{
                mt: "4px",
                justifyContent: "flex-end",
            }}
        >
            <Button
                variant="outlined"
                sx={{ mr: "20px" }}
                onClick={() => handleCancelAction()}
            >
                CANCEL
            </Button>
            <Button variant="contained" onClick={() => handleSaveAction()}>
                SAVE SETTINGS
            </Button>
        </HorizontalContainer>
    </VerticalContainer>
);

const Settings = () => {
    const history = useHistory();
    const [menuId, setMenuId] = useState(1);
    const setPage = useSetRecoilState(pageName);
    const resetPageName = useResetRecoilState(pageName);
    const resetData = useResetRecoilState(organizationId);
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

    const handleSave = () => {
        /* TODO: Handle Save Action */
    };
    const handleCancel = () => {
        history.push("/");
    };

    return (
        <Box p="40px 85px" height="100%">
            <Box height="98%">
                <Grid container spacing={2} height="80%">
                    <Grid item xs={12} md={3}>
                        <LoadingWrapper height={237} isLoading={skeltonLoading}>
                            <RoundedCard
                                sx={{
                                    p: "12px 20px",
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
                                    <Divider />
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
                            <RoundedCard sx={{ p: "40px 44px 24px" }}>
                                <TabPanel index={1} value={menuId}>
                                    <UserSettings />
                                    <ActionButtons
                                        handleSaveAction={handleSave}
                                        handleCancelAction={handleCancel}
                                    />
                                </TabPanel>
                                <TabPanel value={menuId} index={2}>
                                    <NetworkSettings />
                                    <ActionButtons
                                        handleSaveAction={handleSave}
                                        handleCancelAction={handleCancel}
                                    />
                                </TabPanel>
                                <TabPanel value={menuId} index={3}>
                                    <AlertSettings />
                                    <ActionButtons
                                        handleSaveAction={handleSave}
                                        handleCancelAction={handleCancel}
                                    />
                                </TabPanel>
                                <TabPanel value={menuId} index={4}>
                                    <NodeSettings />
                                    <ActionButtons
                                        handleSaveAction={handleSave}
                                        handleCancelAction={handleCancel}
                                    />
                                </TabPanel>
                            </RoundedCard>
                        </LoadingWrapper>
                    </Grid>
                </Grid>
            </Box>
            {!skeltonLoading && (
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
            )}
        </Box>
    );
};

export default Settings;
