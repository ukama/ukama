'use client';
import AppSnackbar from '@/components/AppSnackbar/page';
import { CenterContainer } from '@/styles/global';
import GradientWrapper from '@/wrappers/gradiantWrapper';

const ConfigureLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  return (
    <CenterContainer>
      <GradientWrapper>{children}</GradientWrapper>
      <AppSnackbar />
    </CenterContainer>
  );
};

export default ConfigureLayout;
