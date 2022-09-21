import { useEffect, useState } from "react";
import { OnBoardingFlow } from "../../components";
import { useHistory } from "react-router-dom";
import {
    useGetAccountDetailsQuery,
    useAddUserMutation,
    useUpdateFirstVisitMutation,
} from "../../generated";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { snackbarMessage, user } from "../../recoil";
import { CircularProgress } from "@mui/material";
import { CenterContainer } from "../../styles";

const OnBoarding = () => {
    const history = useHistory();
    const setNodeToastNotification = useSetRecoilState(snackbarMessage);
    const [userData, setUserData] = useState<any>();
    const [simAdded, setSimAdded] = useState<boolean>();
    const setUser = useSetRecoilState(user);
    const getUser = useRecoilValue(user);
    const { data: _user, loading } = useGetAccountDetailsQuery();
    useEffect(() => {
        if (!loading)
            setUser({
                ...getUser,
                has_logged_once: _user?.getAccountDetails.isFirstVisit,
            });
    }, [loading]);
    const [updateFirstVisit] = useUpdateFirstVisitMutation({
        onCompleted: res => {
            if (res) {
                setUser({
                    ...getUser,
                    has_logged_once: res.updateFirstVisit.firstVisit,
                });
            }
        },
    });
    const [addUser] = useAddUserMutation({
        onCompleted: res => {
            setUserData(res?.addUser);
            setSimAdded(true);
        },
        onError: err => {
            if (err?.message) {
                setNodeToastNotification({
                    id: "error-add-user-success",
                    message: `${err?.message}`,
                    type: "error",
                    show: true,
                });
            }
        },
    });

    useEffect(() => {
        if (getUser.has_logged_once == false) {
            history.push("/home");
        }
    }, [getUser.has_logged_once]);

    const handleEsimInstallation = (values: any) => {
        addUser({
            variables: {
                data: {
                    email: values.email,
                    name: values.name,
                    status: values.status,
                },
            },
        });
    };
    const handleNetworkSetup = () => {
        //handle network installation
    };
    const handleSkip = () => {
        updateFirstVisit({
            variables: {
                data: {
                    firstVisit: false,
                },
            },
        });
    };
    const goToConsole = () => {
        updateFirstVisit({
            variables: {
                data: {
                    firstVisit: false,
                },
            },
        });
    };

    return (
        <>
            {loading ||
                (getUser.has_logged_once == undefined && (
                    <CenterContainer>
                        <CircularProgress />
                    </CenterContainer>
                ))}
            {getUser.has_logged_once == true && (
                <OnBoardingFlow
                    handleEsimInstallation={handleEsimInstallation}
                    handleNetworkSetup={handleNetworkSetup}
                    goToConsole={goToConsole}
                    qrCodeId={userData?.iccid || ""}
                    name={userData?.name || ""}
                    simAdded={simAdded}
                    currentUser={getUser}
                    handleSkip={handleSkip}
                />
            )}
        </>
    );
};
export default OnBoarding;
