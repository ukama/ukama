import React from 'react';
import { Tabs, Tab } from '@mui/material';

interface TabsProps {
  selectedTab: number;
  handleTabChange: (event: React.SyntheticEvent, newValue: number) => void;
}

const TabsComponent: React.FC<TabsProps> = ({
  selectedTab,
  handleTabChange,
}) => {
  return (
    <Tabs value={selectedTab} onChange={handleTabChange} aria-label="tabs menu">
      <Tab label="Information" />
      <Tab label="Data Usage" />
      <Tab label="SIMs" />
      <Tab label="History" />
    </Tabs>
  );
};

export default TabsComponent;
