import {
    HomeIcon,
    UsersIcon,
    RouterIcon,
    BillingIcon,
    AccountIcon,
    SettingsIcon,
    ModuleStoreIcon,
    NotificationIcon,
} from "../assets/svg";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import { DataTableWithOptionColumns } from "./tableColumns";
import {
    checkPasswordLength,
    checkPasswordSpecialCharacter,
    combineUppercaseNLowercaseValidator,
} from "../utils";
import { MenuItemType } from "../types";

const DRAWER_WIDTH = 211;
const APP_VERSION = "v0.0.1";
const COPY_RIGHTS = "Copyright © Ukama Inc. 2021";
const PasswordRules = [
    {
        id: 1,
        idLabel: "checkLength",
        label: "8 characters minimum",
        validator: checkPasswordLength,
    },
    {
        id: 2,
        idLabel: "checkSpecialCaracter",
        label: "At least one special character",
        validator: checkPasswordSpecialCharacter,
    },
    {
        id: 3,
        idLabel: "checkLowerCaseNUppercase",
        label: "Upper & lowercase letters ",
        validator: combineUppercaseNLowercaseValidator,
    },
];

const SIDEBAR_MENU1 = [
    { id: "1", title: "Home", Icon: HomeIcon, route: "/home" },
    { id: "2", title: "Nodes", Icon: RouterIcon, route: "/nodes" },
    { id: "3", title: "User", Icon: UsersIcon, route: "/user" },
    { id: "4", title: "Billing", Icon: BillingIcon, route: "/billing" },
];
const STATS_OPTIONS = [
    { id: 1, label: "Connected", value: "Connected" },
    { id: 2, label: "Device uptime", value: "Device uptime" },
    { id: 3, label: "Throughput", value: "Throughput" },
];
const STATS_PERIOD = [
    { id: 1, label: "DAY" },
    { id: 2, label: "WEEK " },
    { id: 3, label: "MONTH " },
];

const HEADER_MENU = [
    { id: "1", Icon: SettingsIcon, title: "Setting" },
    { id: "2", Icon: NotificationIcon, title: "Notification" },
    { id: "3", Icon: AccountIcon, title: "Account" },
];

const SIDEBAR_MENU2 = [
    { id: "5", title: "Integrations", Icon: ModuleStoreIcon, route: "/store" },
];
const MONTH_FILTER = [
    { id: 1, label: "January ", value: "january " },
    { id: 2, label: "February", value: "february" },
    { id: 3, label: "March", value: "march " },
    { id: 4, label: "April", value: "april" },
    { id: 5, label: "May", value: "may" },
    { id: 6, label: "June", value: "june" },
    { id: 7, label: "July", value: "july" },
    { id: 8, label: "August", value: "august" },
    { id: 9, label: "September", value: "september" },
    { id: 10, label: "October", value: "october" },
    { id: 11, label: "November", value: "november" },
    { id: 12, label: "December", value: "december" },
];

const TIME_FILTER = [
    { id: 1, label: "Today", value: "today" },
    { id: 2, label: "This week", value: "week" },
    { id: 3, label: "Month", value: "month" },
    { id: 4, label: "Total", value: "total" },
    { id: 5, label: "Current", value: "current" },
];

const NETWORKS = [
    { id: 1, label: "Public Network", value: "public" },
    { id: 2, label: "Private Network", value: "private" },
];

const BASIC_MENU_ACTIONS: MenuItemType[] = [
    { id: 1, Icon: EditIcon, title: "Edit", route: "edit" },
    {
        id: 2,
        Icon: DeleteIcon,
        title: "Delete",
        route: "delete",
    },
];

const DEACTIVATE_EDIT_ACTION_MENU: MenuItemType[] = [
    {
        id: 1,
        Icon: DeleteIcon,
        title: "Deactivate",
        route: "deactivate",
    },
    { id: 2, Icon: EditIcon, title: "Edit", route: "edit" },
];

const UserActivation = {
    title: "Activate User",
    subTitle:
        "Choose all the eSIM(s) you want to assign to a user at this time. Once you pair an eSIM with a user ....................... policy details ................................",
};

const BillingTabs = [
    { id: 1, label: "CURRENT BILL", value: "1" },
    { id: 2, label: "BILLING HISTORY", value: "2" },
];

const UsersTabs = [
    { id: 1, label: "OVERVIEW", value: "1" },
    { id: 2, label: "CURRENTLY CONNECTED", value: "2" },
];
export {
    NETWORKS,
    UsersTabs,
    COPY_RIGHTS,
    APP_VERSION,
    TIME_FILTER,
    MONTH_FILTER,
    BillingTabs,
    HEADER_MENU,
    DRAWER_WIDTH,
    SIDEBAR_MENU1,
    SIDEBAR_MENU2,
    PasswordRules,
    STATS_OPTIONS,
    STATS_PERIOD,
    UserActivation,
    BASIC_MENU_ACTIONS,
    DataTableWithOptionColumns,
    DEACTIVATE_EDIT_ACTION_MENU,
};
