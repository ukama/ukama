import { RoundedCard } from "../../styles";
import { Typography, Divider } from "@mui/material";

type CurrentBillProps = {
    title: string;
    amount: string;
    periodOf: string;
};

const CurrentBill = ({ title, amount, periodOf }: CurrentBillProps) => {
    return (
        <RoundedCard>
            <Typography variant="h6" sx={{ textTransform: "capitalize" }}>
                {title}
            </Typography>
            <Typography variant="body2">{periodOf}</Typography>

            <Divider />
            <Typography variant="h3" sx={{ m: "18px 0px" }}>
                {amount}
            </Typography>
        </RoundedCard>
    );
};
export default CurrentBill;
