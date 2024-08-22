'use client';
import Welcome from '@/components/Welcome';
import { useAppContext } from '@/context';
import { CenterContainer } from '@/styles/global';
import GradientWrapper from '@/wrappers/gradiantWrapper';
import { useRouter } from 'next/navigation';

const Page = () => {
  const { env, user } = useAppContext();
  const router = useRouter();
  return (
    <CenterContainer>
      <GradientWrapper>
        <Welcome
          handleNext={() => router.push('/configure')}
          handleBack={() => router.push(`${env.AUTH_APP_URL}/user/logout`)}
          orgName={user.orgName}
          operatingCountry={'Dominican Republic of Congo'}
        />
      </GradientWrapper>
    </CenterContainer>
  );
};

export default Page;
