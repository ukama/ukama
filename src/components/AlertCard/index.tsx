import { colors } from "../../theme";
import { Box, Typography, Card, Grid } from "@mui/material";
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
            sx={{
                marginBottom: 1,
                width: "100%",
                padding: 1,
            }}
        >
            <Grid container spacing={1} p={1}>
                <Grid item xs>
                    <div
                        style={{
                            display: "flex",
                            alignItems: "center",
                            flexWrap: "wrap",
                        }}
                    >
                        <Icon />
                        <Typography variant="body1" color="initial">
                            {title}
                        </Typography>
                    </div>
                </Grid>
                <Grid item>
                    <Typography variant="caption" color={colors.empress}>
                        {date}
                    </Typography>
                </Grid>
                <Grid item xs={12}>
                    <Box pl={3.7}>
                        <Typography variant="subtitle2" color="initial">
                            {description}
                        </Typography>
                    </Box>
                </Grid>
            </Grid>
        </Card>
    );
};
export default AlertCard;
