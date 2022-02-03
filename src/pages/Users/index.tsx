import {
    UserCard,
    ContainerHeader,
    UserDetailsDialog,
    PagePlaceholder,
    LoadingWrapper,
} from "../../components";
import { Box } from "@mui/material";
import { useRecoilValue } from "recoil";
import { useState, useEffect } from "react";
import { isSkeltonLoading, organizationId } from "../../recoil";
import {
    GetUserDto,
    Get_User_Status_Type,
    useGetUserLazyQuery,
    useMyUsersLazyQuery,
} from "../../generated";
import { RoundedCard } from "../../styles";
import { UserData } from "../../constants/stubData";

const User = () => {
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [showSimDialog, setShowSimDialog] = useState(false);
    const [userForm, setUserForm] = useState<GetUserDto>({
        //TODO: Remove these static initilization after UI review
        id: "123",
        status: Get_User_Status_Type.Active,
        name: "Test123",
        eSimNumber: "123456789",
        iccid: "0987654321",
        email: "test@123.com",
        phone: "1231231",
        roaming: true,
        dataPlan: 100.0,
        dataUsage: 12.3,
    });
    const orgId = useRecoilValue(organizationId);
    const handleSimDialog = () => {
        setShowSimDialog(false);
    };

    const [getUsersByOrgId, { data: usersRes, loading: usersByOrgLoading }] =
        useMyUsersLazyQuery();

    const [getUser, { data: userRes, loading: userLoading }] =
        useGetUserLazyQuery();

    useEffect(() => {
        if (orgId) {
            getUsersByOrgId({
                variables: {
                    orgId: orgId,
                },
            });
        }
    }, [orgId]);

    useEffect(() => {
        if (userRes) {
            setUserForm({ ...userRes.getUser });
        }
    }, [userRes]);

    const getUseDetails = (simDetailsId: string) => {
        setShowSimDialog(true);
        getUser({
            variables: {
                id: simDetailsId,
            },
        });
    };

    const handleSimInstallation = () => {
        setShowSimDialog(true);
    };

    const getSearchValue = () => {
        /* TODO: Handle HandleSearch */
    };

    const handleSave = () => {
        /* TODO: CALL UPDATE USER API HERE */
    };

    return (
        <Box sx={{ mt: 3, height: "calc(100% - 15%)" }}>
            <LoadingWrapper
                width="100%"
                height="100%"
                isLoading={isSkeltonLoad || userLoading || usersByOrgLoading}
            >
                <RoundedCard sx={{ borderRadius: "4px" }}>
                    {usersRes?.myUsers?.users || UserData ? (
                        <>
                            <ContainerHeader
                                title="My Users"
                                showButton={true}
                                showSearchBox={true}
                                buttonTitle="INSTALL SIMS"
                                handleSearchChange={getSearchValue}
                                handleButtonAction={handleSimInstallation}
                                stats={`${
                                    usersRes?.myUsers?.users?.length ||
                                    UserData.length
                                }`}
                            />
                            <UserCard
                                userDetails={
                                    usersRes?.myUsers?.users || UserData
                                }
                                handleMoreUserdetails={getUseDetails}
                            />
                        </>
                    ) : (
                        <PagePlaceholder
                            hyperlink=""
                            showActionButton={true}
                            buttonTitle="Install sims"
                            handleAction={() => setShowSimDialog(true)}
                            description="No users on network. Install SIMs to get started."
                        />
                    )}
                </RoundedCard>
                <UserDetailsDialog
                    user={userForm}
                    saveBtnLabel="save"
                    closeBtnLabel="close"
                    isOpen={showSimDialog}
                    setUserForm={setUserForm}
                    handleClose={handleSimDialog}
                    simDetailsTitle="SIM Details"
                    handleSaveSimUser={handleSave}
                    userDetailsTitle="User Details"
                />
            </LoadingWrapper>
        </Box>
    );
};

export default User;
