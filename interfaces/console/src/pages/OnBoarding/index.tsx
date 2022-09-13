import { useEffect, useState } from "react";
import { OnBoardingFlow } from "../../components";
import { useHistory } from "react-router-dom";
import {
    useGetAccountDetailsQuery,
    useAddUserMutation,
    useUpdateFirstVisitMutation,
} from "../../generated";
import { useSetRecoilState } from "recoil";
import { snackbarMessage } from "../../recoil";

const OnBoarding = () => {
    const history = useHistory();
    const setNodeToastNotification = useSetRecoilState(snackbarMessage);
    const [userData, setUserData] = useState<any>();
    const [simAdded, setSimAdded] = useState<boolean>();
    const { data: account } = useGetAccountDetailsQuery();
    const [updateFirstVisit, { loading: updateVisitLoading }] =
        useUpdateFirstVisitMutation({
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
        if (!account?.getAccountDetails?.isFirstVisit) {
            history.push("/");
        }
    }, [updateVisitLoading]);

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
        />
    );
};
export default OnBoarding;
