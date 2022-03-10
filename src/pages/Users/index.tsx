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
    useGetUsersByOrgQuery,
} from "../../generated";
import { useState } from "react";
import { useRecoilValue } from "recoil";
import { RoundedCard } from "../../styles";
import { Box, Card, Grid } from "@mui/material";
import { isSkeltonLoading, organizationId } from "../../recoil";

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
    const orgId = useRecoilValue(organizationId) || "";
    const [users, setUsers] = useState<GetUserDto[]>([]);
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [showSimDialog, setShowSimDialog] = useState(false);
    const [selectedUser, setSelectedUser] = useState<GetUserDto>(userInit);
    // eslint-disable-next-line no-unused-vars
    const [userForm, setUserForm] = useState<GetUserDto>(userInit);

    const { data: usersRes, loading: usersByOrgLoading } =
        useGetUsersByOrgQuery({
            variables: { orgId: orgId },
            onCompleted: res => setUsers(res.getUsersByOrg.users),
        });

    const handleSimDialogClose = () => setShowSimDialog(false);

    const onViewMoreClick = (user: GetUserDto) => {
        setShowSimDialog(true);
        setSelectedUser(user);
    };

    const handleSimInstallation = () => {
        setSelectedUser(userInit);
        setShowSimDialog(true);
    };

    const getSearchValue = (search: string) => {
        if (search.length > 2) {
            setUsers(
                users.filter(user =>
                    user.name.toLocaleLowerCase().includes(search)
                )
            );
        } else {
            setUsers(usersRes?.getUsersByOrg?.users || []);
        }
    };

    const handleSave = () => {
        /* TODO: CALL UPDATE USER API HERE */
    };

    return (
        <Box component="div" sx={{ height: "calc(100% - 3%)" }}>
            <LoadingWrapper
                width="100%"
                height="inherit"
                isLoading={isSkeltonLoad || usersByOrgLoading}
            >
                {usersRes && usersRes?.getUsersByOrg?.users?.length > 0 ? (
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
                        <Grid container spacing={2} mt={2}>
                            {users.map((item: GetUserDto) => (
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
                    isOpen={showSimDialog}
                    setUserForm={setUserForm}
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
