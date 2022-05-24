import {
    UserCard,
    ContainerHeader,
    UserDetailsDialog,
    PagePlaceholder,
    LoadingWrapper,
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
    useGetUsersDataUsageLazyQuery,
    useGetUsersDataUsageSSubscription,
    UserInputDto,
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
                        phone: userRes.phone,
                        dataPlan: userRes.dataPlan,
                        dataUsage: userRes.dataUsage,
                    },
                    ...users.slice(index + 1),
                ]);
            }
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

    const handleSimInstallationClose = () => setShowInstallSim(false);

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

    const handleEsimInstallation = (eSimData: UserInputDto) => {
        if (eSimData) {
            addUser({
                variables: {
                    data: {
                        email: eSimData.email,
                        name: eSimData.name,
                        status: eSimData.status || false,
                        phone: "",
                    },
                },
            });
        }
    };
    const handlePhysicalSimInstallation = (physicalSimData: UserInputDto) => {
        if (physicalSimData) {
            addUser({
                variables: {
                    data: {
                        email: physicalSimData.email,
                        name: physicalSimData.name,
                        status: physicalSimData.status || false,
                        phone: "",
                    },
                },
            });
        }
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
                        phone: selectedUser.phone,
                        status: selectedUser.status,
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
                        userDetailsTitle="User Details"
                        handleClose={handleSimDialogClose}
                        userStatusLoading={updateUserStatusLoading}
                        handleServiceAction={handleUpdateUserStatus}
                        handleSubmitAction={handleUserSubmitAction}
                    />
                )}

                {showInstallSim && (
                    <AddUser
                        iSeSimAdded={isEsimAdded}
                        handlePhysicalSimInstallation={
                            handlePhysicalSimInstallation
                        }
                        loading={addUserLoading}
                        handleEsimInstallation={handleEsimInstallation}
                        addedUserName={newAddedUserName}
                        qrCodeId={qrCodeId}
                        isOpen={showInstallSim}
                        handleClose={handleSimInstallationClose}
                    />
                )}
            </LoadingWrapper>
        </Box>
    );
};

export default User;
