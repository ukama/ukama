import {
    HomeIcon,
    UsersIcon,
    RouterIcon,
    BillingIcon,
    ModuleStoreIcon,
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
    { id: "1", title: "Home", Icon: HomeIcon, route: "/dashboard" },
    { id: "2", title: "Nodes", Icon: RouterIcon, route: "/nodes" },
    { id: "3", title: "User", Icon: UsersIcon, route: "/user" },
    { id: "4", title: "Billing", Icon: BillingIcon, route: "/billing" },
];

const SIDEBAR_MENU2 = [
    { id: "5", title: "Module Store", Icon: ModuleStoreIcon, route: "/store" },
];

export { passwordRules, DRAWER_WIDTH, SIDEBAR_MENU1, SIDEBAR_MENU2 };
