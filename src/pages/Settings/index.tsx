import {
    UserSettings,
    NodeSettings,
    AlertSettings,
    NetworkSettings,
} from "../../components";
import { colors } from "../../theme";
import React, { useState } from "react";
import { SETTING_MENU } from "../../constants";
import { SettingsMenuTypes } from "../../types";
import { HorizontalContainer, RoundedCard } from "../../styles";
import { Box, Grid, Button, Divider, MenuList, MenuItem } from "@mui/material";

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
            style={{ padding: "22px 0px" }}
            id={`currentTab-indexpanel-${index}`}
            aria-labelledby={`currentTab-index-${index}`}
        >
            {value === index && <Box>{children}</Box>}
        </div>
    );
};

const ActionButtons = ({
    handleCancelAction,
    handleSaveAction,
}: ActionButtonsProps) => (
    <HorizontalContainer
        sx={{
            mt: "12px",
            width: "90%",
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
);

const Settings = () => {
    const [menuId, setMenuId] = useState(1);
    const handleItemClick = (id: number) => setMenuId(id);
    const handleLogout = () => {
        /* TODO: Handle Logout */
    };
    const handleSave = () => {
        /* TODO: Handle Save Action */
    };
    const handleCancel = () => {
        /* TODO: Handle Cancel Action */
    };

    return (
        <>
            <Box mt={2}>
                <Grid container spacing={2}>
                    <Grid item xs={2}>
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
                    </Grid>
                    <Grid item xs={10}>
                        <RoundedCard>
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
                    </Grid>
                </Grid>
            </Box>
        </>
    );
};

export default Settings;
