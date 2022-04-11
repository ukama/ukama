import ApexLineChart from "../ApexLineChart";
import { RoundedCard, SkeletonRoundedCard } from "../../styles";

type StatsCardProps = {
    loading: boolean;
    metricData: any;
};

const StatsCard = ({ loading, metricData }: StatsCardProps) => {
    return (
        <>
            {loading ? (
                <SkeletonRoundedCard variant="rectangular" height={337} />
            ) : (
                <RoundedCard sx={{ minHeight: 337, display: "flex" }}>
                    <ApexLineChart data={metricData["temperaturetrx"]} />
                </RoundedCard>
            )}
        </>
    );
};
export default StatsCard;
