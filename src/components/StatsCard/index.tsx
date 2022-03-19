import ApexLineChartIntegration from "../ApexLineChart";
import { RoundedCard, SkeletonRoundedCard } from "../../styles";

type StatsCardProps = {
    loading: boolean;
    metricData: any;
    hasMetricData: boolean;
};

const StatsCard = ({ loading, metricData, hasMetricData }: StatsCardProps) => {
    return (
        <>
            {loading ? (
                <SkeletonRoundedCard variant="rectangular" height={337} />
            ) : (
                <RoundedCard sx={{ minHeight: 337, display: "flex" }}>
                    <ApexLineChartIntegration
                        hasData={hasMetricData}
                        data={[metricData]}
                        name={metricData.name}
                    />
                </RoundedCard>
            )}
        </>
    );
};
export default StatsCard;
