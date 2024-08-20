'use client';
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
    </CenterContainer>
  );
};

export default ConfigureLayout;
