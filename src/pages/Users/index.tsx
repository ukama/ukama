import {
    UserCard,
    ContainerHeader,
    UserDetailsDialog,
    PagePlaceholder,
    LoadingWrapper,
} from "../../components";
import {
    GetUserDto,
    Get_User_Status_Type,
    useGetUserLazyQuery,
    useMyUsersLazyQuery,
} from "../../generated";
import { useRecoilValue } from "recoil";
import { UserData } from "../../constants";
import { RoundedCard } from "../../styles";
import { useState, useEffect } from "react";
import { Box, Card, Grid } from "@mui/material";
import { isSkeltonLoading, organizationId } from "../../recoil";

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
        <Box component="div" sx={{ height: "calc(100% - 3%)" }}>
            <LoadingWrapper
                width="100%"
                height="inherit"
                isLoading={isSkeltonLoad || userLoading || usersByOrgLoading}
            >
                {(usersRes && usersRes?.myUsers?.users?.length > 0) ||
                UserData?.length > 0 ? (
                    <RoundedCard sx={{ borderRadius: "4px", overflow: "auto" }}>
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
                        <Grid container spacing={2} mt={2}>
                            {UserData.map(
                                ({
                                    id,
                                    name,
                                    dataPlan,
                                    dataUsage,
                                    eSimNumber,
                                }: any) => (
                                    <Grid key={id} item xs={12} md={6} lg={3}>
                                        <Card
                                            variant="outlined"
                                            sx={{
                                                padding: "15px 18px 8px 18px",
                                            }}
                                        >
                                            <UserCard
                                                id={id}
                                                name={name}
                                                dataPlan={dataPlan}
                                                dataUsage={dataUsage}
                                                eSimNumber={eSimNumber}
                                                handleMoreUserdetails={
                                                    getUseDetails
                                                }
                                            />
                                        </Card>
                                    </Grid>
                                )
                            )}
                        </Grid>
                    </RoundedCard>
                ) : (
                    <PagePlaceholder
                        hyperlink=""
                        showActionButton={true}
                        buttonTitle="Install sims"
                        handleAction={() => setShowSimDialog(true)}
                        description="No users on network. Install SIMs to get started."
                    />
                )}

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
