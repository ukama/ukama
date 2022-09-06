import { OnBoardingFlow } from "../../components";
const OnBoarding = () => {
    const handleEsimInstallation = () => {
        //handle sim installation
    };
    const handleNetworkSetup = () => {
        //handle network installation
    };
    return (
        <OnBoardingFlow
            handleEsimInstallation={handleEsimInstallation}
            handleNetworkSetup={handleNetworkSetup}
        />
    );
};
export default OnBoarding;
