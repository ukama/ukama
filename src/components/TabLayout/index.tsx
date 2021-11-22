import { Box, Tabs, Tab, Button } from "@mui/material";
import { SelectItemType } from "../../types";
type TabLayoutProps = {
    tab: string;
    tabs: SelectItemType[];
    onTabChange: Function;
    withActionButton?: boolean;
    handleAction?: Function;
};

const TabLayoutHeader = ({
    tab,
    tabs,
    onTabChange,
    withActionButton = false,
    handleAction = () => {},
}: TabLayoutProps) => {
    return (
        <Box
            sx={{
                width: "100%",
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
            {withActionButton && (
                <Button
                    size="medium"
                    color="primary"
                    variant="contained"
                    sx={{ height: "42px" }}
                    onClick={() => handleAction()}
                >
                    ACTIVATE USER
                </Button>
            )}
        </Box>
    );
};

export default TabLayoutHeader;
