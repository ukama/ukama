import { useState } from "react";
import { Box } from "@mui/material";
import { UsersTabs } from "../../constants";
import { ContainerHeader } from "../../components";
import { RoundedCard } from "../../styles";

const User = () => {
    const [tab, setTab] = useState("1");

    const handleTabChange = (value: string) => setTab(value);

    const handleActivateUser = () => {
        /* TODO: Handle activate user action */
    };
    return (
        <RoundedCard>
            <ContainerHeader
                title="My Users"
                stats={"78/2"}
                handleButtonAction={handleActivateUser}
                buttonTitle="INSTALL SIMS"
                withSearchBox
            />
        </RoundedCard>
    );
};

export default User;
