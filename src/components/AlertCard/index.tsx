import { colors } from "../../theme";
import { SkeletonRoundedCard } from "../../styles";
import { Typography, Card, CardHeader } from "@mui/material";
import Avatar from "@mui/material/Avatar";
type AlertCardProps = {
    Icon: any;
    id: number;
    date: string;
    title: string;
    loading?: boolean;
    description: string;
};

const AlertCard = ({
    date,
    description,
    title,
    Icon,
    loading,
}: AlertCardProps) => {
    return (
        <>
            {loading ? (
                <SkeletonRoundedCard variant="rectangular" height={64} />
            ) : (
                <Card
                    sx={{
                        width: "100%",
                        marginBottom: "3px",
                    }}
                    elevation={1}
                >
                    <CardHeader
                        sx={{ p: "10px 20px 10px 0px !important" }}
                        avatar={
                            <Avatar
                                sx={{
                                    bgcolor: "white",
                                    position: "relative",
                                    bottom: "12px",
                                }}
                            >
                                <Icon />
                            </Avatar>
                        }
                        action={
                            <Typography
                                variant="caption"
                                color={colors.empress}
                            >
                                {date}
                            </Typography>
                        }
                        title={
                            <Typography
                                variant="body1"
                                color="initial"
                                style={{ position: "relative", right: "17px" }}
                            >
                                {title}
                            </Typography>
                        }
                        subheader={
                            <Typography
                                variant="body2"
                                color="initial"
                                style={{ position: "relative", right: "17px" }}
                            >
                                {description}
                            </Typography>
                        }
                    />
                </Card>
            )}
        </>
    );
};
export default AlertCard;
