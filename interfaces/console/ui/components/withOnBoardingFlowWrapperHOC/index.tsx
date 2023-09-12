import { Box } from '@mui/material';
import { FunctionComponent } from 'react';
import { RootContainer, GradiantBar, ComponentContainer } from './style';

const withOnBoardingFlowWrapperHOC = (
  WrappedComponent: FunctionComponent<any>,
) => {
  return function HOC(props: any) {
    return (
      <RootContainer maxWidth="sm" disableGutters>
        <GradiantBar />

        <Box sx={ComponentContainer}>
          <WrappedComponent {...props} />
        </Box>
      </RootContainer>
    );
  };
};

export default withOnBoardingFlowWrapperHOC;
