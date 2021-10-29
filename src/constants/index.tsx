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

const DRAWER_WIDTH = 240;
import {
    checkPasswordSpecialCharacter,
    checkPasswordLength,
    checkPasswordHasLetters,
} from "../utils";

const passwordRules = [
    {
        id: 1,
        label: "Be a minimum of 8 characters",
        validator: checkPasswordLength,
    },
    {
        id: 2,
        label: "At least one special character",
        validator: checkPasswordSpecialCharacter,
    },
    {
        id: 3,
        label: "Upper & lowercase letters ",
        validator: checkPasswordHasLetters,
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
    { id: "5", title: "Module Store", Icon: ModuleStoreIcon, route: "/store" },
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
    { id: 3, label: "This month", value: "month" },
    { id: 4, label: "Total", value: "total" },
];

const NETWORKS = [
    { id: 1, label: "Public Network", value: "public" },
    { id: 2, label: "Private Network", value: "private" },
];
export {
    NETWORKS,
    TIME_FILTER,
    MONTH_FILTER,
    passwordRules,
    HEADER_MENU,
    DRAWER_WIDTH,
    SIDEBAR_MENU1,
    SIDEBAR_MENU2,
    STATS_OPTIONS,
    STATS_PERIOD,
};
