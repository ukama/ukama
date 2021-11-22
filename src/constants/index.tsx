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
    { id: 1, label: "January ", value: "JANUARY " },
    { id: 2, label: "February", value: "FEBRUARY" },
    { id: 3, label: "March", value: "MARCH " },
    { id: 4, label: "April", value: "APRIL" },
    { id: 5, label: "May", value: "MAY" },
    { id: 6, label: "June", value: "JUNE" },
    { id: 7, label: "July", value: "JULY" },
    { id: 8, label: "August", value: "AUGUST" },
    { id: 9, label: "September", value: "SEPTEMBER" },
    { id: 10, label: "October", value: "OCTOBER" },
    { id: 11, label: "November", value: "NOVEMBER" },
    { id: 12, label: "December", value: "DECEMBER" },
];

const TIME_FILTER = [
    { id: 1, label: "Today", value: "TODAY" },
    { id: 2, label: "This week", value: "WEEK" },
    { id: 3, label: "This month", value: "MONTH" },
    { id: 4, label: "Total", value: "TOTAL" },
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

export {
    NETWORKS,
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
