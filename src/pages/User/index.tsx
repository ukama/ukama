import { ContainerHeader, UserCard } from "../../components";
import { UserData } from "../../constants/stubData";
import { Box } from "@mui/material";
const User = () => {
    const handleSimInstallation = () => {
        /* TODO: Handle activate user action */
    };
    return (
        <Box sx={{ flexGrow: 1, mt: 3 }}>
            <UserCard userDetails={UserData}>
                <ContainerHeader
                    title="My Users"
                    stats={"78/2"}
                    handleButtonAction={handleSimInstallation}
                    buttonTitle="INSTALL SIMS"
                    withSearchBox
                />
            </UserCard>
        </Box>
    );
};

export default User;
