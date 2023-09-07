import {
    Card,
    FormControlLabel,
    Grid,
    Radio,
    RadioGroup,
    Typography,
} from "@mui/material";
import { useState } from "react";
import { useRecoilState } from "recoil";
import { isDarkmode } from "../../recoil";

const ConsoleSettings = () => {
    const [_isDarkMod, _setIsDarkMod] = useRecoilState(isDarkmode);
    const [apperanceMode, setAppearanceMode] = useState<string>(
        _isDarkMod ? "dark" : "light"
    );

    const handleToggle = (e: any) => {
        setAppearanceMode(e.target.value);
        _setIsDarkMod(e.target.value === "dark");
    };

    return (
        <Grid container spacing={2}>
            <Grid item container spacing={2}>
                <Grid item xs={12} md={4}>
                    <Typography variant="h6">Appearance</Typography>
                </Grid>
                <Grid item container xs={12} md={8}>
                    <RadioGroup
                        onChange={handleToggle}
                        value={apperanceMode}
                        name="appearance-mode"
                        sx={{ width: { xs: "100%", md: "70%" } }}
                    >
                        {[
                            { id: 1, label: "Light", value: "light" },
                            { id: 2, label: "Dark", value: "dark" },
                        ].map(({ value, label }) => (
                            <Card
                                key={value}
                                variant="outlined"
                                sx={{ px: 3, py: 1, mb: 2, width: "100%" }}
                            >
                                <FormControlLabel
                                    value={value}
                                    label={label}
                                    control={<Radio />}
                                    sx={{ width: "100%", typography: "body1" }}
                                />
                            </Card>
                        ))}
                    </RadioGroup>
                </Grid>
            </Grid>
        </Grid>
    );
};

export default ConsoleSettings;
