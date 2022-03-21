import {
    ListItemIcon,
    ListItemText,
    MenuItem,
    Popover,
    IconButton,
} from "@mui/material";
import { useState } from "react";
import { MenuItemType } from "../../types";
import MenuDots from "@mui/icons-material/MoreHoriz";

type ItemProps = {
    Icon: any;
    type: string;
    title: string;
    handleItemClick: Function;
};

const OptionItem = ({ type, Icon, title, handleItemClick }: ItemProps) => (
    <MenuItem onClick={() => handleItemClick(type)}>
        <ListItemIcon>
            <Icon fontSize="small" />
        </ListItemIcon>
        <ListItemText>{title}</ListItemText>
    </MenuItem>
);

type OptionsPopoverProps = {
    cid: string;
    menuOptions: MenuItemType[];
    handleItemClick: Function;
    style?: any;
};

const OptionsPopover = ({
    cid,
    menuOptions,
    handleItemClick,
    style,
}: OptionsPopoverProps) => {
    const [anchorEl, setAnchorEl] = useState(null);
    const handlePopoverClose = () => setAnchorEl(null);
    const handlePopoverOpen = (event: any) => setAnchorEl(event.currentTarget);

    const open = Boolean(anchorEl);
    const id = open ? cid : undefined;
    return (
        <>
            <IconButton
                onClick={handlePopoverOpen}
                aria-describedby={id}
                style={style}
                sx={{ px: 0 }}
            >
                <MenuDots />
            </IconButton>
            <Popover
                id={id}
                open={open}
                anchorEl={anchorEl}
                onClose={handlePopoverClose}
                anchorOrigin={{
                    vertical: "bottom",
                    horizontal: "right",
                }}
                transformOrigin={{
                    vertical: "top",
                    horizontal: "right",
                }}
            >
                {menuOptions.map(({ id: optId, Icon, title, route }: any) => (
                    <OptionItem
                        key={`${cid}-${optId}`}
                        type={route}
                        Icon={Icon}
                        title={title}
                        handleItemClick={(type: string) => {
                            handleItemClick(type);
                            handlePopoverClose();
                        }}
                    />
                ))}
            </Popover>
        </>
    );
};

export default OptionsPopover;
