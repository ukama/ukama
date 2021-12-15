import { ContainerHeader, UserCard } from "../../components";
import { RoundedCard } from "../../styles";
import { Box } from "@mui/material";
const User = () => {
    const handleActivateUser = () => {
        /* TODO: Handle activate user action */
    };
    return (
        <Box sx={{ flexGrow: 1, pb: "18px", mt: 2 }}>
            <RoundedCard>
                <ContainerHeader
                    title="My Users"
                    stats={"78/2"}
                    handleButtonAction={handleActivateUser}
                    buttonTitle="INSTALL SIMS"
                    withSearchBox
                />

                <UserCard />
            </RoundedCard>
        </Box>
    );
};

export default User;
