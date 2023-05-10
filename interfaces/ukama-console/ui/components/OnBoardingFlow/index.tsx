import '@/i18n/i18n';
import { useEffect, useState } from 'react';
import { NetworkSetup } from '../../components';
import ESimQR from '../../components/AddUser/ESimQR';
import Userform from '../../components/AddUser/Userform';
import withOnBoardingFlowWrapperHOC from '../withOnBoardingFlowWrapperHOC';

interface IOnBoardingFlow {
  handleNetworkSetup: Function;
  handleEsimInstallation: Function;
  goToConsole: Function;
  qrCodeId: string | undefined;
  name: string | undefined;
  simAdded: boolean;
  currentUser: any;
  handleSkip: Function;
  loading: boolean;
}

const OnBoardingFlow = ({
  handleEsimInstallation,
  handleNetworkSetup,
  qrCodeId,
  goToConsole,
  name,
  handleSkip,
  currentUser,
  loading = false,
  simAdded = false,
}: IOnBoardingFlow) => {
  const [step, setstep] = useState(1);
  const nextStep = () => {
    setstep(step + 1);
  };
  const handleSimType = () => {
    //getSim type
  };
  const prevStep = () => {
    setstep(step - 1);
  };
  const handleSimData = (data: any) => {
    handleEsimInstallation(data);
  };

  useEffect(() => {
    if (simAdded) {
      setstep(step + 1);
    }
  }, [simAdded]);
  switch (step) {
    case 1:
      return (
        <NetworkSetup nextStep={nextStep} networkData={handleNetworkSetup} />
      );
    case 2:
      return (
        <Userform
          getSimType={handleSimType}
          title={`Connect to the network`}
          currentUser={currentUser}
          handleSkip={handleSkip}
          handleGoBack={prevStep}
          loading={loading}
          description={
            'Start accessing high quality and fast data now. You’ll be able to add more users to the network later.'
          }
          handleSimInstallation={handleSimData}
          isAddUser={false}
        />
      );
    case 3:
      return (
        <ESimQR
          title={`Add user successful `}
          description={`You have successfully added ${name} as a user to your network, and an eSIM installation invitation has been sent out to them. To install now, scan the QR code below.`}
          qrCodeId={qrCodeId}
          isOnBoarding={true}
          goToConsole={goToConsole}
        />
      );

    default:
      return (
        <NetworkSetup nextStep={nextStep} networkData={handleNetworkSetup} />
      );
  }
};

export default withOnBoardingFlowWrapperHOC(OnBoardingFlow);
