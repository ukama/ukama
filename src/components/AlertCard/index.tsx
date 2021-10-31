import { RoundedCard } from "../../styles";
import { Box, Typography, CardActions, Card } from "@mui/material";
import { useTranslation } from "react-i18next";
import { colors } from "../../theme";
import { AlertItemType } from "../../types";
import "../../i18n/i18n";
type AlertCardProps = {
    alertCardItems: AlertItemType[];
};
const AlertCard = ({ alertCardItems }: AlertCardProps) => {
    const { t } = useTranslation();
    return (
        <>
            <RoundedCard>
                <Box mb={2}>
                    <Typography variant="h6">{t("ALERT.Title")}</Typography>
                </Box>
                {alertCardItems.map(
                    ({ id, date, description, title, Icon }: AlertItemType) => (
                        <Card elevation={2} key={id} sx={{ marginBottom: 1 }}>
                            <CardActions>
                                <Icon />
                                <div style={{ padding: 7 }}>{title}</div>
                                <Typography
                                    variant="caption"
                                    color={colors.empress}
                                >
                                    {date}
                                </Typography>
                            </CardActions>
                            <Box pl={5.5}>
                                <Typography variant="subtitle2" color="initial">
                                    {description}
                                </Typography>
                            </Box>
                        </Card>
                    )
                )}
            </RoundedCard>
        </>
    );
};
export default AlertCard;
