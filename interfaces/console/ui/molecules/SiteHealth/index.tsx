import React from 'react';
import { styled, Box } from '@mui/material';
import { BatteryFull, Wifi, Settings } from '@mui/icons-material';

interface BaseStationSiteHealthProps {
  batteryLevel: number;
  internetSwitch: boolean;
  controllerSwitch: boolean;
}

const BaseStationSiteHealthContainer = styled(Box)`
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
`;

const BatteryIcon = styled(BatteryFull)`
  color: ${(props) => props.theme.palette.primary.main};
`;

const WifiIcon = styled(Wifi)`
  color: ${(props) => props.theme.palette.primary.main};
`;

const SettingsIcon = styled(Settings)`
  color: ${(props) => props.theme.palette.primary.main};
`;

const Line = styled(Box)`
  position: absolute;
  width: 2px;
  height: 50%;
  background-color: ${(props) => props.theme.palette.primary.main};
`;

const BatteryContainer = styled(Box)`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-right: 16px;
`;

const SwitchContainer = styled(Box)`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-left: 16px;
`;

const BaseStationSiteHealth: React.FC<BaseStationSiteHealthProps> = ({
  batteryLevel,
  internetSwitch,
  controllerSwitch,
}) => {
  return (
    <BaseStationSiteHealthContainer>
      <BatteryContainer>
        <BatteryIcon />
        <Line />
        <div>{batteryLevel}%</div>
      </BatteryContainer>
      <SwitchContainer>
        <WifiIcon />
        <div>{internetSwitch ? 'On' : 'Off'}</div>
        <Line />
        <SettingsIcon />
        <div>{controllerSwitch ? 'On' : 'Off'}</div>
      </SwitchContainer>
    </BaseStationSiteHealthContainer>
  );
};

export default BaseStationSiteHealth;
