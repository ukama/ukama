/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
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
