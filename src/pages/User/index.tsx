import { ContainerHeader, UserCard, UserDetailsDialog } from "../../components";
import { Box } from "@mui/material";
import { useState } from "react";
import { useMyUsersQuery } from "../../generated";
import { organizationId } from "../../recoil";
import { useRecoilValue } from "recoil";

const User = () => {
    const [showSimDialog, setShowSimDialog] = useState(false);
    const [userId, setUserId] = useState();
    const orgId = useRecoilValue(organizationId);
    const handleSimDialog = () => {
        setShowSimDialog(false);
    };
    const { data: usersRes } = useMyUsersQuery({
        variables: { orgId: orgId || "" },
    });

    const handleSimInstallation = () => {
        //console.log(value);
    };

    /* eslint-disable no-unused-vars */
    const getSearchValue = (searchValue: any) => {
        //console.log(searchValue);
    };
    const getUseDetails = (simDetailsId: any) => {
        setShowSimDialog(true);
        setUserId(simDetailsId);
    };
    /* eslint-disable no-unused-vars */
    const getSimData = (simData: any) => {
        //console.log(simData);
    };
    return (
        <Box sx={{ flexGrow: 1, mt: 3 }}>
            <UserCard
                userDetails={usersRes?.myUsers.users}
                handleMoreUserdetails={getUseDetails}
            >
                <ContainerHeader
                    title="My Users"
                    stats={`${usersRes?.myUsers.users.length}`}
                    handleButtonAction={handleSimInstallation}
                    buttonTitle="INSTALL SIMS"
                    handleSearchChange={getSearchValue}
                    showSearchBox={true}
                    showButton={true}
                />
                <UserDetailsDialog
                    id={userId}
                    userName="John Doe"
                    data="- 1.5 GB data used, 0.5 free GB left"
                    isOpen={showSimDialog}
                    userDetailsTitle="User Details"
                    btnLabel="Submit"
                    handleClose={handleSimDialog}
                    simDetailsTitle="SIM Details"
                    saveBtnLabel="save"
                    closeBtnLabel="close"
                    handleSaveSimUser={getSimData}
                />
            </UserCard>
        </Box>
    );
};

export default User;
