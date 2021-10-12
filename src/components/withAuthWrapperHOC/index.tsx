import { Box } from "@mui/material";
import { FunctionComponent } from "react";
import { RootContainer, GradiantBar, ComponentContainer } from "./style";

const withAuthWrapperHOC = (WrappedComponent: FunctionComponent<any>) => {
    const HOC = () => {
        return (
            <RootContainer maxWidth="sm">
                <GradiantBar />
                <Box sx={ComponentContainer}>
                    <WrappedComponent />
                </Box>
            </RootContainer>
        );
    };
    return HOC;
};

export default withAuthWrapperHOC;
