import { BillingTabs } from "../../constants";
import TabLayoutHeader from "../../components/TabLayout";
import { Box } from "@mui/material";
import { useState } from "react";

const Billing = () => {
    const [tab, setTab] = useState("1");

    const handleTabChange = (value: string) => setTab(value);

    return (
        <Box sx={{ mt: "24px" }}>
            <TabLayoutHeader
                tab={tab}
                tabs={BillingTabs}
                onTabChange={handleTabChange}
            />
            {tab === "1" ? <h4>Current Bill</h4> : <h4>Billing History</h4>}
        </Box>
    );
};

export default Billing;
