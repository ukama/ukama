import { useState } from "react";
import { Box } from "@mui/material";
import { UsersTabs } from "../../constants";
import TabLayoutHeader from "../../components/TabLayout";

const User = () => {
    const [tab, setTab] = useState("1");

    const handleTabChange = (value: string) => setTab(value);

    const handleActivateUser = () => {
        /* TODO: Handle activate user action */
    };

    return (
        <Box sx={{ mt: "24px" }}>
            <TabLayoutHeader
                tab={tab}
                tabs={UsersTabs}
                withActionButton={true}
                onTabChange={handleTabChange}
                handleAction={handleActivateUser}
            />
            {tab === "1" ? <h4>Overview</h4> : <h4>Currently connected</h4>}
        </Box>
    );
};

export default User;
