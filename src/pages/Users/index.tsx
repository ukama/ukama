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
    Get_User_Status_Type,
    useGetUserLazyQuery,
    useGetUsersByOrgQuery,
} from "../../generated";
import { useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { Box, Card, Grid } from "@mui/material";
import { isSkeltonLoading, user } from "../../recoil";

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
    status: Get_User_Status_Type.Active,
};

const User = () => {
    const { id: orgId = "" } = useRecoilValue(user);
    const [users, setUsers] = useState<GetUsersDto[]>([]);
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [showSimDialog, setShowSimDialog] = useState(false);
    const [selectedUser, setSelectedUser] = useState<GetUserDto>(userInit);

    const { data: usersRes, loading: usersByOrgLoading } =
        useGetUsersByOrgQuery({
            variables: { orgId: orgId },
            onCompleted: res => setUsers(res.getUsersByOrg),
        });

    const [getUser, { loading: userLoading }] = useGetUserLazyQuery({
        onCompleted: res => {
            if (res.getUser) setSelectedUser(res.getUser);
        },
    });

    const handleSimDialogClose = () => setShowSimDialog(false);

    const onViewMoreClick = (_user: GetUsersDto) => {
        setShowSimDialog(true);
        getUser({
            variables: {
                userInput: { orgId: orgId, userId: _user.id },
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
                            buttonTitle="INSTALL SIMS"
                            handleSearchChange={getSearchValue}
                            handleButtonAction={handleSimInstallation}
                            stats={`${users.length}`}
                        />
                        <Grid container spacing={2} mt={4}>
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
                />
            </LoadingWrapper>
        </Box>
    );
};

export default User;
