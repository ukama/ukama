import {
    Box,
    ListItemIcon,
    ListItemText,
    MenuItem,
    Popover,
} from "@mui/material";
import { useState } from "react";
import { MenuDots } from "../../assets/svg";
import { MenuItemType } from "../../types";

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
};

const OptionsPopover = ({
    cid,
    menuOptions,
    handleItemClick,
}: OptionsPopoverProps) => {
    const [anchorEl, setAnchorEl] = useState(null);
    const handlePopoverClose = () => setAnchorEl(null);
    const handlePopoverOpen = (event: any) => setAnchorEl(event.currentTarget);

    const open = Boolean(anchorEl);
    const id = open ? cid : undefined;
    return (
        <>
            <Box
                aria-describedby={id}
                onClick={handlePopoverOpen}
                sx={{ cursor: "pointer" }}
            >
                <MenuDots />
            </Box>
            <Popover
                id={id}
                open={open}
                anchorEl={anchorEl}
                onClose={handlePopoverClose}
                anchorOrigin={{
                    vertical: "bottom",
                    horizontal: "left",
                }}
                transformOrigin={{
                    vertical: "top",
                    horizontal: "left",
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
