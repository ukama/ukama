import { ContainerHeader, UserCard } from "../../components";
import { RoundedCard } from "../../styles";
import { UserData } from "../../constants/stubData";
const User = () => {
    const handleSimInstallation = () => {
        /* TODO: Handle activate user action */
    };
    return (
        <RoundedCard sx={{ mt: 3 }}>
            <ContainerHeader
                title="My Users"
                stats={"78/2"}
                handleButtonAction={handleSimInstallation}
                buttonTitle="INSTALL SIMS"
                withSearchBox
            />
            <UserCard userDetails={UserData} />
        </RoundedCard>
    );
};

export default User;
