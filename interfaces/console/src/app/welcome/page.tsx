'use client';
import Welcome from '@/components/Welcome';
import { useAppContext } from '@/context';
import { CenterContainer } from '@/styles/global';
import GradientWrapper from '@/wrappers/gradiantWrapper';
import { useRouter } from 'next/navigation';

const Page = () => {
  const { env } = useAppContext();
  const router = useRouter();
  return (
    <CenterContainer>
      <GradientWrapper>
        <Welcome
          handleNext={() => router.push('/configure')}
          handleBack={() => router.push(`${env.AUTH_APP_URL}/user/logout`)}
          orgName={'salman-org'}
          operatingCountry={'Dominican Republic of Congo'}
        />
      </GradientWrapper>
    </CenterContainer>
  );
};

export default Page;
