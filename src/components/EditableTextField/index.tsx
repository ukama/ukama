import { useState } from "react";
import { colors } from "../../theme";
import EditIcon from "@mui/icons-material/Edit";
import { TextField, IconButton, InputAdornment } from "@mui/material";

type EditableTextFieldProps = {
    label: string;
    type?: string;
    value: string;
    isEditable?: boolean;
    handleOnChange?: Function;
};

const EditableTextField = ({
    label,
    value,
    type = "text",
    isEditable = true,
    // eslint-disable-next-line no-empty-function
    handleOnChange = () => {},
}: EditableTextFieldProps) => {
    const [iseditable, setIsEditable] = useState(false);
    return (
        <TextField
            fullWidth
            id={label}
            name={label}
            label={label}
            value={value}
            variant="standard"
            disabled={!iseditable}
            sx={{ width: "440px" }}
            InputLabelProps={{
                shrink: true,
                // HERE style: { color: colors.lightGrey },
            }}
            onChange={e => handleOnChange(e.target.value)}
            inputRef={input => iseditable && input?.focus()}
            InputProps={{
                type: type,
                disableUnderline: true,
                // HERE style: { color: colors.black },
                endAdornment: (
                    <InputAdornment
                        position="end"
                        sx={{
                            display: isEditable ? "flex" : "none",
                        }}
                    >
                        <IconButton
                            edge="end"
                            onClick={() => setIsEditable(!iseditable)}
                            sx={{
                                svg: {
                                    path: {
                                        fill: `${
                                            iseditable
                                                ? colors.primaryMain
                                                : colors.silver
                                        }`,
                                    },
                                },
                            }}
                        >
                            <EditIcon />
                        </IconButton>
                    </InputAdornment>
                ),
            }}
        />
    );
};

export default EditableTextField;
