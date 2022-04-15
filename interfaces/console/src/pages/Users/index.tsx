import {
    UserCard,
    ContainerHeader,
    UserDetailsDialog,
    PagePlaceholder,
    LoadingWrapper,
} from "../../components";
import {
    GetUserDto,
    GetUsersDto,
    useGetUserLazyQuery,
    useGetUsersByOrgQuery,
    useUpdateUserStatusMutation,
} from "../../generated";
import { useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { Box, Card, Grid } from "@mui/material";
import { isSkeltonLoading } from "../../recoil";

const userInit = {
    id: "",
    name: "",
    iccid: "",
    email: "",
    phone: "",
    dataPlan: 0,
    dataUsage: 0,
    roaming: true,
    eSimNumber: "",
    status: false,
};

const User = () => {
    const [users, setUsers] = useState<GetUsersDto[]>([]);
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [showSimDialog, setShowSimDialog] = useState(false);
    const [selectedUser, setSelectedUser] = useState<GetUserDto>(userInit);

    const { data: usersRes, loading: usersByOrgLoading } =
        useGetUsersByOrgQuery({
            onCompleted: res => setUsers(res.getUsersByOrg),
        });

    const [getUser, { loading: userLoading }] = useGetUserLazyQuery({
        onCompleted: res => {
            if (res.getUser) setSelectedUser(res.getUser);
        },
    });

    const [updateUserStatus, { loading: updateUserLoading }] =
        useUpdateUserStatusMutation();

    const handleSimDialogClose = () => setShowSimDialog(false);

    const onViewMoreClick = (_user: GetUsersDto) => {
        setShowSimDialog(true);
        getUser({
            variables: {
                userId: _user.id,
            },
        });
    };

    const handleSimInstallation = () => {
        setSelectedUser(userInit);
        setShowSimDialog(true);
    };

    const getSearchValue = (search: string) => {
        if (search.length > 2) {
            setUsers(
                users.filter((_user: GetUsersDto) =>
                    _user.name.toLocaleLowerCase().includes(search)
                )
            );
        } else {
            setUsers(usersRes?.getUsersByOrg || []);
        }
    };

    const handleUpdateUserStatus = (
        id: string,
        iccid: string,
        status: boolean
    ) => {
        updateUserStatus({
            variables: {
                data: {
                    userId: id,
                    simId: iccid,
                    status: status,
                },
            },
        });
    };

    const handleSave = () => {
        setShowSimDialog(false);
    };

    return (
        <Box component="div" sx={{ height: "calc(100% - 3%)" }}>
            <LoadingWrapper
                width="100%"
                height="inherit"
                isLoading={isSkeltonLoad || usersByOrgLoading}
            >
                {usersRes && usersRes?.getUsersByOrg?.length > 0 ? (
                    <RoundedCard sx={{ borderRadius: "4px", overflow: "auto" }}>
                        <ContainerHeader
                            title="My Users"
                            showButton={true}
                            showSearchBox={true}
                            buttonSize="medium"
                            buttonTitle="INSTALL SIMS"
                            handleSearchChange={getSearchValue}
                            handleButtonAction={handleSimInstallation}
                            stats={`${users.length}`}
                        />
                        <Grid container spacing={2} mt={{ xs: 2, md: 4 }}>
                            {users.map((item: GetUsersDto) => (
                                <Grid key={item.id} item xs={12} md={6} lg={3}>
                                    <Card
                                        variant="outlined"
                                        sx={{
                                            padding: "15px 18px 8px 18px",
                                        }}
                                    >
                                        <UserCard
                                            user={item}
                                            handleMoreUserdetails={
                                                onViewMoreClick
                                            }
                                        />
                                    </Card>
                                </Grid>
                            ))}
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
                    user={selectedUser}
                    saveBtnLabel="save"
                    closeBtnLabel="close"
                    loading={userLoading}
                    isOpen={showSimDialog}
                    setUserForm={setSelectedUser}
                    simDetailsTitle="SIM Details"
                    handleSaveSimUser={handleSave}
                    userDetailsTitle="User Details"
                    handleClose={handleSimDialogClose}
                    userStatusLoading={updateUserLoading}
                    handleServiceAction={handleUpdateUserStatus}
                />
            </LoadingWrapper>
        </Box>
    );
};

export default User;
