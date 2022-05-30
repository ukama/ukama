import PhysicalSimform from "./PhysicalSimform";
import PhysicalSimFlow2 from "./PhysicalSimFlow2";
import Success from "./Success";
interface IPhysicalSimformFlowPops {
    handleClose: Function;
    handlePhysicalSimInstallationFlow1: Function;
    handlePhysicalSimInstallationFlow2: Function;
    step: number;
}
const PhysicalSimFlow = ({
    step,
    handleClose,
    handlePhysicalSimInstallationFlow1,
    handlePhysicalSimInstallationFlow2,
}: IPhysicalSimformFlowPops) => {
    switch (step) {
        case 1:
            return (
                <PhysicalSimform
                    handleClose={handleClose}
                    handlePhysicalSimInstallation={
                        handlePhysicalSimInstallationFlow1
                    }
                    description="Add user xyz. Physical SIM instructions."
                />
            );
        case 2:
            return (
                <PhysicalSimFlow2
                    handleClose={handleClose}
                    handlePhysicalSimflow2={handlePhysicalSimInstallationFlow2}
                    description="Enter security code for Physical SIM lorem ipsum. Instructions for remembering to install SIM after?"
                />
            );
        case 3:
            return (
                <Success
                    description={
                        "You have successfully added [Name] as a user to your network. Instructions for installing physical SIM (might need more thinking if this process is complex)."
                    }
                />
            );

        default:
            return <></>;
    }
};

export default PhysicalSimFlow;
