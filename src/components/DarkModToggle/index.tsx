import { IconButton } from "@mui/material";
import { useRecoilState } from "recoil";
import { isDarkMod } from "../../recoil";
import { Brightness3, Brightness7 } from "@mui/icons-material";

const DarkModToggle = () => {
    const [_isDarkMod, _setIsDarkMod] = useRecoilState(isDarkMod);
    const icon = _isDarkMod ? <Brightness7 /> : <Brightness3 />;
    const handleToggle = () => _setIsDarkMod(!_isDarkMod);

    return (
        <IconButton
            edge="end"
            color="inherit"
            aria-label="mode"
            onClick={handleToggle}
            sx={{ p: "8px" }}
        >
            {icon}
        </IconButton>
    );
};

export default DarkModToggle;
