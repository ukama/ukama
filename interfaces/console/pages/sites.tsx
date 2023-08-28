import { commonData } from '@/app-recoil';
import { MONTH_FILTER, TIME_FILTER } from '@/constants';
import { useGetSitesQuery } from '@/generated';
import { DataBilling, DataUsage, UsersWithBG } from '@/public/svg';
import { TCommonData } from '@/types';
import StatusCard from '@/ui/components/StatusCard';
import EmptyView from '@/ui/molecules/EmptyView';
import {
  LabelOverlayUI,
  SitesSelection,
  SitesTree,
} from '@/ui/molecules/NetworkMap/OverlayUI';
import NetworkStatus from '@/ui/molecules/NetworkStatus';
import NetworkIcon from '@mui/icons-material/Hub';
import { Paper } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2';
import dynamic from 'next/dynamic';
import { useRecoilValue } from 'recoil';
const DynamicMap = dynamic(
  () => import('../ui/molecules/NetworkMap/DynamicMap'),
  {
    ssr: false,
  },
);

export default function Page() {
  const _commonData = useRecoilValue<TCommonData>(commonData);

  const { data: networkRes, loading: networkLoading } = useGetSitesQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      networkId: _commonData?.networkId,
    },
  });

  return (
    <>
      <Grid container spacing={2}>
        <Grid xs={12}>
          <NetworkStatus
            loading={false}
            availableNodes={4}
            statusType="ONLINE"
            tooltipInfo="Network is online"
          />
        </Grid>
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            Icon={UsersWithBG}
            title={'Connected Users'}
            options={TIME_FILTER}
            subtitle1={'0'}
            subtitle2={''}
            option={''}
            loading={false}
            handleSelect={(value: string) => {}}
          />
        </Grid>
        <Grid xs={12} md={6} lg={4}>
          <StatusCard
            title={'Data Usage'}
            subtitle1={`0`}
            subtitle2={`Package`}
            Icon={DataUsage}
            options={TIME_FILTER}
            option={'usage'}
            loading={false}
            handleSelect={(value: string) => {}}
          />
        </Grid>
      </Grid>
    </>
  );
}
