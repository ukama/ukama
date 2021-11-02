import { RoundedCard } from "../../styles";
import {
    Box,
    Typography,
    CardActions,
    Card,
    List,
    ListItem,
} from "@mui/material";
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
                <List
                    disablePadding
                    sx={{
                        position: "relative",
                        overflow: "auto",
                        maxHeight: 300,
                        width: "100%",
                        maxWidth: 360,
                    }}
                >
                    {alertCardItems.map(
                        ({
                            id,
                            date,
                            description,
                            title,
                            Icon,
                        }: AlertItemType) => (
                            <ListItem key={id} style={{ padding: 1 }}>
                                <Card
                                    elevation={2}
                                    sx={{
                                        marginBottom: 1,
                                        width: "100%",
                                    }}
                                >
                                    <CardActions>
                                        <Icon />
                                        <div style={{ padding: 7 }}>
                                            {title}
                                        </div>
                                        <Typography
                                            variant="caption"
                                            color={colors.empress}
                                        >
                                            {date}
                                        </Typography>
                                    </CardActions>
                                    <Box pl={5.5}>
                                        <Typography
                                            variant="subtitle2"
                                            color="initial"
                                        >
                                            {description}
                                        </Typography>
                                    </Box>
                                </Card>
                            </ListItem>
                        )
                    )}
                </List>
            </RoundedCard>
        </>
    );
};
export default AlertCard;
