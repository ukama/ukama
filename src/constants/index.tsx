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
const STATS_OPTIONS = ["Connected users", "Device uptime", "Close"];
const STATS_PERIOD = ["DAY", "WEEK", "MONTH"];
const HEADER_MENU = [
    { id: "1", Icon: SettingsIcon, title: "Setting" },
    { id: "2", Icon: NotificationIcon, title: "Notification" },
    { id: "3", Icon: AccountIcon, title: "Account" },
];

const SIDEBAR_MENU2 = [
    { id: "5", title: "Module Store", Icon: ModuleStoreIcon, route: "/store" },
];

export {
    passwordRules,
    HEADER_MENU,
    DRAWER_WIDTH,
    SIDEBAR_MENU1,
    SIDEBAR_MENU2,
    STATS_OPTIONS,
    STATS_PERIOD,
};
