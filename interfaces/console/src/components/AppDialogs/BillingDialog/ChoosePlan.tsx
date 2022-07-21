import { useState } from "react";
import { colors } from "../../../theme";
import { BILLING_PLANS } from "../../../constants";
import { RadioGroup, Card, FormControlLabel, Radio } from "@mui/material";

const ChoosePlan = () => {
    const [plan, setPlan] = useState("default");
    const handleToggle = (e: any) => {
        setPlan(e.target.value);
    };
    return (
        <RadioGroup value={plan} onChange={handleToggle} sx={{ mt: 3 }}>
            {BILLING_PLANS.map(({ value, label }) => (
                <Card
                    key={value}
                    variant="outlined"
                    sx={{
                        mb: 2,
                        width: "100%",
                        cursor: "pointer",
                        ":hover": {
                            border: `1px solid ${colors.primaryMain}`,
                        },
                    }}
                >
                    <FormControlLabel
                        value={value}
                        label={label}
                        control={<Radio />}
                        sx={{
                            m: 0,
                            px: 3,
                            py: 1,
                            width: "100%",
                            typography: "body1",
                        }}
                    />
                </Card>
            ))}
        </RadioGroup>
    );
};

export default ChoosePlan;
