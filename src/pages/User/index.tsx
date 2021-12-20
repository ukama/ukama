import { ContainerHeader, UserCard, UserDetailsDialog } from "../../components";
import { UserData } from "../../constants/stubData";
import { Box } from "@mui/material";
import { useState } from "react";

const User = () => {
    const [showSimDialog, setShowSimDialog] = useState(false);

    const handleSimDialog = () => {
        setShowSimDialog(false);
    };

    const handleSimInstallation = () => {
        //console.log(value);
    };

    /* eslint-disable no-unused-vars */
    const getSearchValue = (searchValue: any) => {
        //console.log(searchValue);
    };
    const getUseDetails = (id: any) => {
        setShowSimDialog(true);
    };
    const getSimData = (simData: any) => {
        //console.log(simdata);
    };
    return (
        <Box sx={{ flexGrow: 1, mt: 3 }}>
            <UserCard
                userDetails={UserData}
                handleMoreUserdetails={getUseDetails}
            >
                <ContainerHeader
                    title="My Users"
                    stats={"78/2"}
                    handleButtonAction={handleSimInstallation}
                    buttonTitle="INSTALL SIMS"
                    handleSearchChange={getSearchValue}
                    withSearchBox
                />
                <UserDetailsDialog
                    userName="John Doe"
                    data="- 1.5 GB data used, 0.5 free GB left"
                    isOpen={showSimDialog}
                    userDetailsTitle="User Details"
                    btnLabel="Submit"
                    handleClose={handleSimDialog}
                    simDetailsTitle="SIM Details"
                    saveBtnLabel="save"
                    closeBtnLabel="close"
                    handleSaveSimUser={getSimData}
                />
            </UserCard>
        </Box>
    );
};

export default User;
