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
    useAddUserMutation,
    useGetUserLazyQuery,
    useGetUsersByOrgQuery,
    useUpdateUserMutation,
    useUpdateUserStatusMutation,
} from "../../generated";
import { useEffect, useState } from "react";
import { RoundedCard } from "../../styles";
import { Box, Card, Grid } from "@mui/material";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { isSkeltonLoading, snackbarMessage } from "../../recoil";

const userInit = {
    id: "",
    name: "",
    iccid: "",
    email: "",
    phone: "",
    dataPlan: 0,
    dataUsage: 0,
    roaming: false,
    eSimNumber: "",
    status: false,
};

const User = () => {
    const [users, setUsers] = useState<GetUsersDto[]>([]);
    const isSkeltonLoad = useRecoilValue(isSkeltonLoading);
    const [simDialog, setSimDialog] = useState({
        isShow: false,
        type: "add",
    });
    const [selectedUser, setSelectedUser] = useState<GetUserDto>(userInit);
    const setUserNotification = useSetRecoilState(snackbarMessage);
    const [
        addUser,
        { loading: addUserLoading, data: addUserRes, error: addUserError },
    ] = useAddUserMutation();
    const [
        updateUser,
        {
            loading: updateUserLoading,
            data: updateUserRes,
            error: updateUserError,
        },
    ] = useUpdateUserMutation();
    const { data: usersRes, loading: usersByOrgLoading } =
        useGetUsersByOrgQuery({
            onCompleted: res => setUsers(res.getUsersByOrg),
        });

    const [getUser, { loading: userLoading }] = useGetUserLazyQuery({
        onCompleted: res => {
            if (res.getUser) setSelectedUser(res.getUser);
        },
    });

    const [updateUserStatus, { loading: updateUserStatusLoading }] =
        useUpdateUserStatusMutation({
            onCompleted: res => {
                if (res) {
                    setSelectedUser({
                        ...selectedUser,
                        status: res.updateUserStatus.carrier.services.data,
                        roaming: res.updateUserStatus.ukama.services.data,
                    });
                }
            },
        });

    useEffect(() => {
        if (addUserRes) {
            setUserNotification({
                id: "addUserNotification",
                message: `The user has been added successfully!`,
                type: "success",
                show: true,
            });
        }
    }, [addUserRes]);

    useEffect(() => {
        if (updateUserRes) {
            setUserNotification({
                id: "updateUserNotification",
                message: `The user has been updated successfully!`,
                type: "success",
                show: true,
            });
        }
    }, [updateUserRes]);

    useEffect(() => {
        if (addUserError) {
            setUserNotification({
                id: "addUserNotification",
                message: `${addUserError.message}`,
                type: "error",
                show: true,
            });
        }
    }, [addUserError]);

    useEffect(() => {
        if (updateUserError) {
            setUserNotification({
                id: "updateUserNotification",
                message: `${updateUserError.message}`,
                type: "error",
                show: true,
            });
        }
    }, [updateUserError]);

    const handleSimDialogClose = () =>
        setSimDialog({ ...simDialog, isShow: false });

    const onViewMoreClick = (_user: GetUsersDto) => {
        setSimDialog({ isShow: true, type: "edit" });
        getUser({
            variables: {
                userId: _user.id,
            },
        });
    };

    const handleSimInstallation = () => {
        setSelectedUser(userInit);
        setSimDialog({ isShow: true, type: "add" });
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

    const handleUserSubmitAction = () => {
        handleSimDialogClose();
        if (simDialog.type === "add") {
            addUser({
                variables: {
                    data: {
                        email: selectedUser.email,
                        name: selectedUser.name,
                        phone: selectedUser.phone,
                    },
                },
            });
        } else if (simDialog.type === "edit" && selectedUser.id) {
            updateUser({
                variables: {
                    userId: selectedUser.id,
                    data: {
                        email: selectedUser.email,
                        name: selectedUser.name,
                        phone: selectedUser.phone,
                    },
                },
            });
        }
    };

    return (
        <Box component="div" sx={{ height: "calc(100% - 3%)" }}>
            <LoadingWrapper
                width="100%"
                height="inherit"
                isLoading={
                    isSkeltonLoad ||
                    usersByOrgLoading ||
                    addUserLoading ||
                    updateUserLoading
                }
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
                        handleAction={() =>
                            setSimDialog({ isShow: true, type: "add" })
                        }
                        description="No users on network. Install SIMs to get started."
                    />
                )}

                {simDialog.isShow && (
                    <UserDetailsDialog
                        user={selectedUser}
                        type={simDialog.type}
                        saveBtnLabel={"Save"}
                        closeBtnLabel="close"
                        loading={
                            userLoading || addUserLoading || updateUserLoading
                        }
                        isOpen={simDialog.isShow}
                        setUserForm={setSelectedUser}
                        simDetailsTitle="SIM Details"
                        userDetailsTitle="User Details"
                        handleClose={handleSimDialogClose}
                        userStatusLoading={updateUserStatusLoading}
                        handleServiceAction={handleUpdateUserStatus}
                        handleSubmitAction={handleUserSubmitAction}
                    />
                )}
            </LoadingWrapper>
        </Box>
    );
};

export default User;
