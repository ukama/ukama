import "../../i18n/i18n";
import { useState } from "react";
import withOnBoardingFlowWrapperHOC from "../withOnBoardingFlowWrapperHOC";
import { NetworkSetup } from "../../components";
import Userform from "../../components/AddUser/Userform";
import ESimQR from "../../components/AddUser/ESimQR";
import { CircularProgress, Grid } from "@mui/material";

interface IOnBoardingFlow {
    handleNetworkSetup: Function;
    handleEsimInstallation: Function;
    goToConsole: Function;
    qrCodeId: string | undefined;
    name: string | undefined;
    simAdded: boolean;
    currentUser: any;
    handleSkip: Function;
}

const OnBoardingFlow = ({
    handleEsimInstallation,
    handleNetworkSetup,
    qrCodeId,
    goToConsole,
    name,
    handleSkip,
    currentUser,
    simAdded = false,
}: IOnBoardingFlow) => {
    const [step, setstep] = useState(1);
    const nextStep = () => {
        setstep(step + 1);
    };

    const prevStep = () => {
        setstep(step - 1);
    };
    const handleSimData = (data: any) => {
        setstep(step + 1);
        handleEsimInstallation(data);
    };

    switch (step) {
        case 1:
            return (
                <NetworkSetup
                    nextStep={nextStep}
                    networkData={handleNetworkSetup}
                />
            );
        case 2:
            return (
                <Userform
                    title={`Connect to the network`}
                    currentUser={currentUser}
                    handleSkip={handleSkip}
                    handleGoBack={prevStep}
                    description={
                        "Start accessing high quality and fast data now. You’ll be able to add more users to the network later."
                    }
                    handleSimInstallation={handleSimData}
                    isAddUser={false}
                />
            );
        case 3:
            return simAdded ? (
                <ESimQR
                    title={`Add user successful `}
                    description={`You have successfully added ${name} as a user to your network, and an eSIM installation invitation has been sent out to them. To install now, scan the QR code below.`}
                    qrCodeId={qrCodeId}
                    isOnBoarding={true}
                    goToConsole={goToConsole}
                />
            ) : (
                <Grid container justifyContent={"center"}>
                    <CircularProgress />
                </Grid>
            );

        default:
            return (
                <NetworkSetup
                    nextStep={nextStep}
                    networkData={handleNetworkSetup}
                />
            );
    }
};

export default withOnBoardingFlowWrapperHOC(OnBoardingFlow);
