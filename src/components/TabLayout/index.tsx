import { Box, Tabs, Tab } from "@mui/material";
import { SelectItemType } from "../../types";
type TabLayoutProps = {
    tab: string;
    tabs: SelectItemType[];
    onTabChange: Function;
};

const TabLayoutHeader = ({ tab, tabs, onTabChange }: TabLayoutProps) => {
    return (
        <Box sx={{ width: "100%", typography: "body1" }}>
            <Box
                sx={{
                    borderBottom: 2,
                    borderColor: "rgba(196, 196, 196, 0.1)",
                }}
            >
                <Tabs
                    value={tab}
                    onChange={(_, b) => onTabChange(b)}
                    aria-label="wrapped label tabs example"
                >
                    {tabs.map(({ id, label, value }: SelectItemType) => (
                        <Tab
                            key={id}
                            value={value}
                            label={label}
                            wrapped
                            sx={{
                                p: "0px 36px",
                                fontSize: "15px",
                                fontWeight: 500,
                            }}
                        />
                    ))}
                </Tabs>
            </Box>
        </Box>
    );
};

export default TabLayoutHeader;
