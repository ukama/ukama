/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Icon registry — maps the design's material icon names (used in nav config
 * and datasets) to MUI rounded icons, so configs stay serializable strings.
 */
import type { SvgIconProps } from '@mui/material/SvgIcon';
import AccountTreeRounded from '@mui/icons-material/AccountTreeRounded';
import AppsRounded from '@mui/icons-material/AppsRounded';
import BadgeRounded from '@mui/icons-material/BadgeRounded';
import BatteryAlertRounded from '@mui/icons-material/BatteryAlertRounded';
import BatteryChargingFullRounded from '@mui/icons-material/BatteryChargingFullRounded';
import BoltRounded from '@mui/icons-material/BoltRounded';
import LightModeRounded from '@mui/icons-material/LightModeRounded';
import SettingsInputAntennaRounded from '@mui/icons-material/SettingsInputAntennaRounded';
import CellTowerRounded from '@mui/icons-material/CellTowerRounded';
import DonutSmallRounded from '@mui/icons-material/DonutSmallRounded';
import ErrorRounded from '@mui/icons-material/ErrorRounded';
import GroupRounded from '@mui/icons-material/GroupRounded';
import HomeRounded from '@mui/icons-material/HomeRounded';
import HubRounded from '@mui/icons-material/HubRounded';
import InfoRounded from '@mui/icons-material/InfoRounded';
import InsightsRounded from '@mui/icons-material/InsightsRounded';
import LocationOnRounded from '@mui/icons-material/LocationOnRounded';
import ManageAccountsRounded from '@mui/icons-material/ManageAccountsRounded';
import MonetizationOnRounded from '@mui/icons-material/MonetizationOnRounded';
import NetworkCheckRounded from '@mui/icons-material/NetworkCheckRounded';
import NotificationsRounded from '@mui/icons-material/NotificationsRounded';
import PaymentsRounded from '@mui/icons-material/PaymentsRounded';
import PersonRounded from '@mui/icons-material/PersonRounded';
import RouterRounded from '@mui/icons-material/RouterRounded';
import SettingsRounded from '@mui/icons-material/SettingsRounded';
import SimCardRounded from '@mui/icons-material/SimCardRounded';
import SupportAgentRounded from '@mui/icons-material/SupportAgentRounded';
import SyncRounded from '@mui/icons-material/SyncRounded';
import WarningRounded from '@mui/icons-material/WarningRounded';

const REGISTRY: Record<string, React.ComponentType<SvgIconProps>> = {
  home: HomeRounded,
  payments: PaymentsRounded,
  group: GroupRounded,
  donut_small: DonutSmallRounded,
  apps: AppsRounded,
  monetization_on: MonetizationOnRounded,
  manage_accounts: ManageAccountsRounded,
  sim_card: SimCardRounded,
  location_on: LocationOnRounded,
  router: RouterRounded,
  account_tree: AccountTreeRounded,
  insights: InsightsRounded,
  hub: HubRounded,
  badge: BadgeRounded,
  support_agent: SupportAgentRounded,
  settings: SettingsRounded,
  notifications: NotificationsRounded,
  person: PersonRounded,
  cell_tower: CellTowerRounded,
  battery_alert: BatteryAlertRounded,
  battery_charging_full: BatteryChargingFullRounded,
  bolt: BoltRounded,
  light_mode: LightModeRounded,
  settings_input_antenna: SettingsInputAntennaRounded,
  network_check: NetworkCheckRounded,
  sync: SyncRounded,
  error: ErrorRounded,
  warning: WarningRounded,
  info: InfoRounded,
};

/** Render a registered icon by its design name. */
export function Ic({ name, ...props }: { name: string } & SvgIconProps) {
  const Cmp = REGISTRY[name] ?? InfoRounded;
  return <Cmp {...props} />;
}
