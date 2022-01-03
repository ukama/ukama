import { useState } from "react";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import { Box, Menu, MenuItem, Button, Typography } from "@mui/material";

interface TextSelectProps {
    value: number;
    options: string[];
    setValue: (_value: number) => void;
}
const TextSelect = ({ value, setValue, options }: TextSelectProps) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const open = Boolean(anchorEl);

    const handleClose = () => {
        setAnchorEl(null);
    };

    return (
        <Box sx={{ marginLeft: 1, marginRight: 1 }}>
            <Button
                aria-haspopup="true"
                id="custom-select-text-button"
                endIcon={<ArrowDropDownIcon color="action" />}
                aria-expanded={open ? "true" : undefined}
                onClick={event => setAnchorEl(event.currentTarget)}
                aria-controls={open ? "custom-select-text" : undefined}
            >
                <Typography variant={"h6"}>{options[value]}</Typography>
            </Button>
            <Menu
                open={open}
                anchorEl={anchorEl}
                onClose={handleClose}
                id="custom-select-text"
                MenuListProps={{
                    "aria-labelledby": "custom-select-text-button",
                }}
            >
                {options.map((option, index) => (
                    <MenuItem
                        key={index}
                        value={option}
                        onClick={() => {
                            handleClose();
                            setValue(index);
                        }}
                    >
                        {option}
                    </MenuItem>
                ))}
            </Menu>
        </Box>
    );
};
export default TextSelect;
