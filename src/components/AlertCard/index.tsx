import { colors } from "../../theme";
import { SkeletonRoundedCard } from "../../styles";
import { Typography, Card, CardHeader } from "@mui/material";
import Avatar from "@mui/material/Avatar";
import IconButton from "@mui/material/IconButton";
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
                    }}
                    elevation={1}
                >
                    <CardHeader
                        avatar={
                            <Avatar
                                sx={{
                                    bgcolor: "white",

                                    m: "0px 0px 24px 0px !important",
                                }}
                            >
                                <Icon />
                            </Avatar>
                        }
                        action={
                            <IconButton aria-label="error_date">
                                <Typography
                                    variant="caption"
                                    color={colors.empress}
                                >
                                    {date}
                                </Typography>
                            </IconButton>
                        }
                        title={
                            <Typography variant="body1" color="initial">
                                {title}
                            </Typography>
                        }
                        subheader={
                            <Typography variant="body2" color="initial">
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
