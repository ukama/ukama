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

const OnBoarding = () => {
    const history = useHistory();
    const setNodeToastNotification = useSetRecoilState(snackbarMessage);

    const getUser = useRecoilValue(user);
    const [userData, setUserData] = useState<any>();
    const [simAdded, setSimAdded] = useState<boolean>();
    const { data: account } = useGetAccountDetailsQuery();
    const [updateFirstVisit] = useUpdateFirstVisitMutation({
        onCompleted: res => {
            if (res) {
                window.location.reload();
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
        if (account?.getAccountDetails.isFirstVisit == false) {
            history.push("/home");
        }
    }, [account?.getAccountDetails.isFirstVisit]);
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
    );
};
export default OnBoarding;
