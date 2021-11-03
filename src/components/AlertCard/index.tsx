import { colors } from "../../theme";
import { Box, Typography, CardActions, Card } from "@mui/material";

type AlertCardProps = {
    Icon: any;
    id: number;
    date: string;
    title: string;
    description: string;
};
const AlertCard = ({ date, description, title, Icon }: AlertCardProps) => {
    return (
        <Card
            elevation={2}
            sx={{
                marginBottom: 1,
                width: "100%",
            }}
        >
            <CardActions>
                <Icon />
                <div style={{ padding: 7 }}>{title}</div>
                <Typography variant="caption" color={colors.empress}>
                    {date}
                </Typography>
            </CardActions>
            <Box pl={5.5}>
                <Typography variant="subtitle2" color="initial">
                    {description}
                </Typography>
            </Box>
        </Card>
    );
};
export default AlertCard;
