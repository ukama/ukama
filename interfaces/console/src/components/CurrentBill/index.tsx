import { RoundedCard } from "../../styles";
import { Typography, Divider, Stack } from "@mui/material";
import PictureAsPdfIcon from "@mui/icons-material/PictureAsPdf";
import colors from "../../theme/colors";
type CurrentBillProps = {
    title: string;
    amount: string;
    periodOf: string;
};

const CurrentBill = ({ title, amount, periodOf }: CurrentBillProps) => {
    return (
        <RoundedCard>
            <Stack direction="row" spacing={1} alignItems="center">
                <Typography variant="h6" sx={{ textTransform: "capitalize" }}>
                    {title}
                </Typography>
                <PictureAsPdfIcon sx={{ color: colors.primaryMain }} />
            </Stack>

            <Typography variant="body2">{periodOf}</Typography>

            <Divider />
            <Typography variant="h3" sx={{ m: "18px 0px" }}>
                {amount}
            </Typography>
        </RoundedCard>
    );
};
export default CurrentBill;
