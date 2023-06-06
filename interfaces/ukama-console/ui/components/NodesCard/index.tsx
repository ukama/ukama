import { BASIC_MENU_ACTIONS, NODE_IMAGES } from '@/constants';
import { Node_Type } from '@/generated';
import UsersIcon from '@mui/icons-material/PeopleAlt';
import { Card, Divider, Grid, Typography, styled } from '@mui/material';
import React from 'react';
import { LoadingWrapper } from '..';
import OptionsPopover from '../../molecules/OptionsPopover';

const CpuIcon = React.lazy(() =>
  import('@/public/svg').then((module) => ({
    default: module.CpuIcon,
  })),
);
const BatteryIcon = React.lazy(() =>
  import('@/public/svg').then((module) => ({
    default: module.BatteryIcon,
  })),
);
const ThermometerIcon = React.lazy(() =>
  import('@/public/svg').then((module) => ({
    default: module.ThermometerIcon,
  })),
);

const Line = styled(Divider)(() => ({
  margin: '18px -18px 4px -18px',
  background: 'rgba(255, 255, 255, 0.12)',
}));

const IconStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
};

type NodeCardProps = {
  id: string;
  type: Node_Type;
  title: string;
  users?: number;
  loading?: boolean;
  subTitle: string;
  isConfigure?: boolean;
  updateShortNote: string;
  isUpdateAvailable: boolean;
  handleOptionItemClick?: Function;
};

const NodeCard = ({
  type,
  title,
  users,
  subTitle,
  loading,
  isUpdateAvailable = false,
  handleOptionItemClick = () => {
    /* Default empty function */
  },
}: NodeCardProps) => {
  return (
    <LoadingWrapper
      isLoading={loading}
      variant="rectangular"
      width={218}
      height={216}
    >
      <Card
        variant="outlined"
        sx={{
          width: '218px',
          height: '216px',
          borderRadius: '10px',
          justifyContent: 'center',
          padding: '15px 18px 8px 18px',
        }}
      >
        <Grid container>
          <Grid item xs={10}>
            <Typography
              variant="subtitle1"
              sx={{
                height: '20px',
                fontWeight: 500,
                overflow: 'hidden',
                lineHeight: '19px',
                display: '-webkit-box',
                letterSpacing: '-0.02em',
                textOverflow: 'ellipsis',
              }}
            >
              {title}
            </Typography>
          </Grid>
          <Grid
            item
            xs={2}
            display={'flex'}
            alignItems="center"
            justifyContent={'flex-end'}
          >
            <OptionsPopover
              style={{ cursor: 'pointer' }}
              cid={'node-card-options'}
              isShowUpdate={isUpdateAvailable}
              menuOptions={BASIC_MENU_ACTIONS}
              handleItemClick={(type: string) => handleOptionItemClick(type)}
            />
          </Grid>
          <Grid item xs={12} textAlign="initial">
            <Typography variant="caption">{subTitle}</Typography>
          </Grid>

          <Grid item xs={12} minHeight={'92px'} sx={{ ...IconStyle, py: 1 }}>
            <img
              src={NODE_IMAGES[type]}
              alt="node-img"
              style={{ maxWidth: '180px', maxHeight: '78px' }}
            />
          </Grid>
          <Grid item xs={12} mb={0.8}>
            <Line />
          </Grid>
          <Grid
            item
            xs={12}
            container
            spacing={1}
            mb={'2px'}
            sx={{ alignItems: 'center' }}
          >
            <Grid
              item
              xs={6}
              container
              display="flex"
              alignSelf="end"
              pt="0px !important"
              alignItems={'flex-end'}
            >
              <UsersIcon sx={{ width: '16px', height: '16px' }} />
              <Typography pl="8px" variant="caption" lineHeight={'14px'}>
                {users}
              </Typography>
            </Grid>
            <Grid xs={2} item sx={{ ...IconStyle }}>
              <ThermometerIcon />
            </Grid>
            <Grid xs={2} item sx={{ ...IconStyle }}>
              <BatteryIcon />
            </Grid>
            <Grid xs={2} item sx={{ ...IconStyle }}>
              <CpuIcon />
            </Grid>
          </Grid>
        </Grid>
      </Card>
    </LoadingWrapper>
  );
};

export default NodeCard;
