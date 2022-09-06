import "../../i18n/i18n";
import { useState } from "react";
import withOnBoardingFlowWrapperHOC from "../withOnBoardingFlowWrapperHOC";
import { NetworkSetup } from "../../components";
import Userform from "../../components/AddUser/Userform";
import ESimQR from "../../components/AddUser/ESimQR";
interface IOnBoardingFlow {
    handleNetworkSetup: Function;
    handleEsimInstallation: Function;
}
const OnBoardingFlow = ({
    handleEsimInstallation,
    handleNetworkSetup,
}: IOnBoardingFlow) => {
    const [step, setstep] = useState(1);
    const nextStep = () => {
        setstep(step + 1);
    };

    // function for going to previous step by decreasing step state by 1
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
                    handleSkip={nextStep}
                    handleGoBack={prevStep}
                    description={
                        "Start accessing high quality and fast data now. You’ll be able to add more users to the network later."
                    }
                    handleSimInstallation={handleSimData}
                />
            );
        case 3:
            return (
                <ESimQR
                    title={`Add user successful `}
                    description={`You have successfully added [Name] as a user to your network, and an eSIM installation invitation has been sent out to them. To install now, scan the QR code below.`}
                    qrCodeId={`123`}
                    isOnBoarding={true}
                />
            );

        default:
            return <div className="App"></div>;
    }
};

export default withOnBoardingFlowWrapperHOC(OnBoardingFlow);
