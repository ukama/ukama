import {
    ContainerHeader,
    UserCard,
    UserDetailsDialog,
    PagePlaceholder,
} from "../../components";
import { Box } from "@mui/material";
import { useState } from "react";
import { useMyUsersQuery, useGetUserQuery } from "../../generated";
import { organizationId } from "../../recoil";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
const User = () => {
    const [showSimDialog, setShowSimDialog] = useState(false);
    const [userId, setUserId] = useState<string>();
    const orgId = useRecoilValue(organizationId);
    const handleSimDialog = () => {
        setShowSimDialog(false);
    };
    const { data: usersRes } = useMyUsersQuery({
        variables: { orgId: orgId || "" },
    });
    const { data: userRes } = useGetUserQuery({
        variables: {
            id: userId || "",
        },
    });
    const handleSimInstallation = () => {
        /* TODO:handle sim activation */
    };

    const getSearchValue = () => {
        /* TODO:handle HandleSearch */
    };
    const getUseDetails = (simDetailsId: any) => {
        setShowSimDialog(true);
        setUserId(simDetailsId);
    };
    const getSimData = () => {
        /* TODO:Get sim Details */
    };
    return (
        <Box sx={{ flexGrow: 1, mt: 3 }}>
            <RoundedCard sx={{ height: "100%" }}>
                {usersRes?.myUsers?.users ? (
                    <>
                        <ContainerHeader
                            title="My Users"
                            stats={`${usersRes?.myUsers?.users.length || "0"}`}
                            handleButtonAction={handleSimInstallation}
                            buttonTitle="INSTALL SIMS"
                            handleSearchChange={getSearchValue}
                            showSearchBox={true}
                            showButton={true}
                        />
                        <UserCard
                            userDetails={usersRes?.myUsers?.users}
                            handleMoreUserdetails={getUseDetails}
                        />

                        <UserDetailsDialog
                            getUser={userRes && userRes?.getUser}
                            isOpen={showSimDialog}
                            userDetailsTitle="User Details"
                            btnLabel="Submit"
                            handleClose={handleSimDialog}
                            simDetailsTitle="SIM Details"
                            saveBtnLabel="save"
                            closeBtnLabel="close"
                            handleSaveSimUser={getSimData}
                        />
                    </>
                ) : (
                    <PagePlaceholder
                        description="No users on network. Install SIMs to get started."
                        hyperlink="Helll"
                        buttonTitle="Install sims"
                        showActionButton={true}
                    />
                )}
            </RoundedCard>
        </Box>
    );
};

export default User;
