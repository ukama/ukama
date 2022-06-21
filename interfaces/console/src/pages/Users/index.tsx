import {
    UserCard,
    ContainerHeader,
    UserDetailsDialog,
    PagePlaceholder,
    LoadingWrapper,
    DeactivateUser,
    AddUser,
} from "../../components";
import {
    GetUserDto,
    GetUsersDto,
    useAddUserMutation,
    useGetUserLazyQuery,
    useGetUsersByOrgQuery,
    useUpdateUserMutation,
    useGetEsimQrLazyQuery,
    useUpdateUserStatusMutation,
    useDeactivateUserMutation,
    useGetUsersDataUsageLazyQuery,
    useGetUsersDataUsageSSubscription,
    UserInputDto,
    useUpdateUserRoamingMutation,
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
    dataPlan: "0",
    dataUsage: "0",
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
    const [showInstallSim, setShowInstallSim] = useState(false);
    const [selectedUser, setSelectedUser] = useState<GetUserDto>(userInit);
    const setUserNotification = useSetRecoilState(snackbarMessage);
    const [qrCodeId, setqrCodeId] = useState<any>();
    const [isEsimAdded, setIsEsimAdded] = useState<boolean>(false);
    const [newAddedUserName, setNewAddedUserName] = useState<any>();
    const [isPsimAdded, setIsPsimAdded] = useState<boolean>(false);
    const [simFlow, setSimFlow] = useState<number>(1);
    const [deactivateUserDialog, setDeactivateUserDialog] = useState({
        isShow: false,
        userId: "",
        userName: "",
    });
    const [addUser, { loading: addUserLoading }] = useAddUserMutation({
        onCompleted: res => {
            if (res?.addUser) {
                setIsEsimAdded(true);
                setNewAddedUserName(res?.addUser?.name);
                handleGetSimQrCode(res?.addUser.id, res?.addUser?.iccid || "");
                refetchResidents();
            }
        },
        onError: err => {
            if (err?.message) {
                setUserNotification({
                    id: "error-add-user",
                    message: `${err?.message}`,
                    type: "error",
                    show: true,
                });
            }
        },
    });
    const [getEsimQrdcodeId, { data: getEsimQrCodeRes }] =
        useGetEsimQrLazyQuery();

    useEffect(() => {
        setqrCodeId(getEsimQrCodeRes?.getEsimQR?.qrCode);
    }, [getEsimQrCodeRes]);
    const [updateUser, { loading: updateUserLoading }] = useUpdateUserMutation({
        onCompleted: res => {
            if (res?.updateUser) {
                setUserNotification({
                    id: "updateUserNotification",
                    message: `The ${res?.updateUser?.name} has been updated successfully!`,
                    type: "success",
                    show: true,
                });
            }
        },
        onError: err => {
            if (err?.message) {
                setUserNotification({
                    id: "updateUserNotification",
                    message: `${err?.message}`,
                    type: "error",
                    show: true,
                });
            }
        },
    });

    const [getUsersDataUsage, { loading: usersDataUsageLoading }] =
        useGetUsersDataUsageLazyQuery();

    useGetUsersDataUsageSSubscription({
        fetchPolicy: "network-only",
        onSubscriptionData: res => {
            if (res.subscriptionData.data?.getUsersDataUsage?.id) {
                const userRes = res.subscriptionData.data?.getUsersDataUsage;
                const index = users.findIndex(item => item.id === userRes.id);
                setUsers([
                    ...users.slice(0, index),
                    {
                        id: userRes.id,
                        name: userRes.name,
                        email: userRes.email,
                        dataPlan: userRes.dataPlan,
                        dataUsage: userRes.dataUsage,
                    },
                    ...users.slice(index + 1),
                ]);
            }
        },
    });

    const [updateUserRoaming, { loading: updateUserRoamingLoading }] =
        useUpdateUserRoamingMutation({
            onCompleted: () => {
                setUserNotification({
                    id: "updateUserRoaming",
                    message: `User roaming status updated.`,
                    type: "success",
                    show: true,
                });
            },
        });

    const {
        data: usersRes,
        loading: usersByOrgLoading,
        refetch: refetchResidents,
    } = useGetUsersByOrgQuery({
        onCompleted: res => {
            setUsers(res.getUsersByOrg);
            getUsersDataUsage({
                variables: {
                    data: { ids: res.getUsersByOrg.map(u => u.id) },
                },
            });
        },
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

    const handleGetSimQrCode = async (userId: string, simId: string) => {
        await getEsimQrdcodeId({
            variables: {
                data: {
                    userId: userId,
                    simId: simId,
                },
            },
        });
    };
    const handleCloseDeactivateUser = () =>
        setDeactivateUserDialog({ ...deactivateUserDialog, isShow: false });

    const [deactivateUser] = useDeactivateUserMutation({
        onCompleted: res => {
            setUserNotification({
                id: "userDeactivated",
                message: `${res.deactivateUser.name} has been deactivated successfully!`,
                type: "success",
                show: true,
            });
            refetchResidents();
        },
        onError: err =>
            setUserNotification({
                id: "userDeactivated",
                message: `${err?.message}`,
                type: "error",
                show: true,
            }),
    });
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

    const handleSimInstallation = () => setShowInstallSim(true);

    const handleSimInstallationClose = () => {
        setShowInstallSim(false);
        setSimFlow(1);
        setIsEsimAdded(false);
    };
    const handleDeactivateUser = () => {
        handleCloseDeactivateUser();
        deactivateUser({
            variables: {
                id: deactivateUserDialog.userId,
            },
        });
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

    const handleUserRoamingAction = (status: boolean) => {
        updateUserRoaming({
            variables: {
                data: {
                    simId: selectedUser.iccid,
                    userId: selectedUser.id,
                    status: status,
                },
            },
        });
    };

    const handleEsimInstallation = (eSimData: UserInputDto) => {
        if (eSimData) {
            addUser({
                variables: {
                    data: {
                        email: eSimData.email,
                        name: eSimData.name,
                        status: eSimData.status || false,
                    },
                },
            });
        }
    };
    const handlePhysicalSimEmailFlow = () => {
        setSimFlow(simFlow + 1);
    };
    const handlePhysicalSimSecurityFlow = () => {
        setSimFlow(simFlow + 1);
        setIsPsimAdded(true);
    };
    const handleUserSubmitAction = () => {
        handleSimDialogClose();
        if (simDialog.type === "edit" && selectedUser.id) {
            updateUser({
                variables: {
                    userId: selectedUser.id,
                    data: {
                        email: selectedUser.email,
                        name: selectedUser.name,
                        status: null,
                    },
                },
            });
        }
    };

    const handleDeactivateAction = (userId: any) => {
        setSimDialog({ ...simDialog, isShow: false });
        setDeactivateUserDialog({
            isShow: true,
            userId: userId,
            userName: users?.find(item => item.id === userId)?.name || "",
        });
    };

    return (
        <Box component="div" sx={{ height: "calc(100% - 3%)" }}>
            <LoadingWrapper
                width="100%"
                height="inherit"
                isLoading={
                    isSkeltonLoad || usersByOrgLoading || updateUserLoading
                }
            >
                {usersRes && usersRes?.getUsersByOrg?.length > 0 ? (
                    <RoundedCard sx={{ borderRadius: "4px", overflow: "auto" }}>
                        <ContainerHeader
                            title="My Users"
                            showButton={true}
                            showSearchBox={true}
                            buttonSize="medium"
                            buttonTitle="ADD USERS"
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
                                            loading={usersDataUsageLoading}
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
                        handleAction={handleSimInstallation}
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
                        userDetailsTitle="User Settings"
                        handleClose={handleSimDialogClose}
                        roamingLoading={updateUserRoamingLoading}
                        userStatusLoading={updateUserStatusLoading}
                        handleServiceAction={handleUpdateUserStatus}
                        handleSubmitAction={handleUserSubmitAction}
                        handleDeactivateAction={handleDeactivateAction}
                        handleUserRoamingAction={handleUserRoamingAction}
                    />
                )}
                {deactivateUserDialog.isShow && (
                    <DeactivateUser
                        isClosable={true}
                        isOpen={deactivateUserDialog.isShow}
                        title={"Deactivate User Confirmation"}
                        description={`${deactivateUserDialog.userName} will be deactivated permanently. Other copy depends on surrounding policy.`}
                        labelSuccessBtn={"DEACTIVATE USER"}
                        labelNegativeBtn={"cancel"}
                        handleCloseAction={handleCloseDeactivateUser}
                        handleSuccessAction={handleDeactivateUser}
                    />
                )}

                {showInstallSim && (
                    <AddUser
                        isPsimAdded={isPsimAdded}
                        iSeSimAdded={isEsimAdded}
                        loading={addUserLoading}
                        handleEsimInstallation={handleEsimInstallation}
                        addedUserName={newAddedUserName}
                        step={simFlow}
                        qrCodeId={qrCodeId}
                        isOpen={showInstallSim}
                        handleClose={handleSimInstallationClose}
                        handlePhysicalSimInstallationFlow1={
                            handlePhysicalSimEmailFlow
                        }
                        handlePhysicalSimInstallationFlow2={
                            handlePhysicalSimSecurityFlow
                        }
                    />
                )}
            </LoadingWrapper>
        </Box>
    );
};

export default User;
