import {
    ContainerHeader,
    UserCard,
    UserDetailsDialog,
    PagePlaceholder,
} from "../../components";
import { Box } from "@mui/material";
import { useState } from "react";
import { UserData } from "../../constants/stubData";
import { useMyUsersQuery, useGetUserQuery } from "../../generated";
import { organizationId } from "../../recoil";
import { useRecoilValue } from "recoil";
const User = () => {
    const [showSimDialog, setShowSimDialog] = useState(false);
    const [userId, setUserId] = useState<any>();
    const orgId = useRecoilValue(organizationId);
    const handleSimDialog = () => {
        setShowSimDialog(false);
    };
    const { data: usersRes } = useMyUsersQuery({
        variables: { orgId: orgId || "" },
    });
    const { data: userRes } = useGetUserQuery({
        variables: {
            id: userId,
        },
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
            {usersRes?.myUsers?.users ? (
                <UserCard
                    userDetails={UserData}
                    handleMoreUserdetails={getUseDetails}
                >
                    <ContainerHeader
                        title="My Users"
                        stats={`${UserData.length}`}
                        handleButtonAction={handleSimInstallation}
                        buttonTitle="INSTALL SIMS"
                        handleSearchChange={getSearchValue}
                        showSearchBox={true}
                        showButton={true}
                    />
                    <UserDetailsDialog
                        getUser={userRes?.getUser}
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
            ) : (
                <PagePlaceholder
                    description="No users on network. Install SIMs to get started."
                    hyperlink="Helll"
                    buttonTitle="Install sims"
                    showActionButton={true}
                />
            )}
        </Box>
    );
};

export default User;
