import React from "react";
import Chart from "react-apexcharts";
import { format } from "date-fns";
import { GraphTitleWrapper } from "..";

const TIME_RANGE_IN_MILLISECONDS = 100;

interface IApexLineChartIntegration {
    data: any;
    name: string;
    filter: string;
    hasData: boolean;
    refreshInterval?: number;
    onRefreshData?: Function;
    onFilterChange?: Function;
}

const ApexLineChart = (props: any) => {
    const options: any = {
        stroke: {
            lineCap: "butt",
            curve: "smooth",
            width: 5,
        },
        chart: {
            minHeight: "200px",
            height: "100%",
            width: "100%",
            zoom: {
                type: "x",
                enabled: false,
                autoScaleYaxis: true,
            },
            animations: {
                enabled: true,
                easing: "linear",
                dynamicAnimation: {
                    speed: 1000,
                },
            },
            dropShadow: {
                enabled: true,
                enabledSeries: [0],
                top: -2,
                left: 2,
                blur: 5,
                opacity: 0.06,
            },
        },
        xaxis: {
            type: "datetime",
            range: props.range,
            labels: {
                formatter: (val: any) =>
                    val ? format(new Date(val * 1000), "mm:ss") : "",
            },
            tooltip: {
                enabled: false,
                offsetX: 0,
            },
        },
        yaxis: {
            labels: {
                formatter: (val: any) => val.toFixed(2),
            },
            // min: 0,
            // max: 100,
            // tooltip: {
            //     enabled: true,
            // },
            tickAmount: 8,
        },
    };
    return (
        <Chart
            type="line"
            key={props.name}
            height={"300px"}
            options={options}
            series={props.dataList}
        />
    );
};

const ApexLineChartIntegration = ({
    name,
    data = [],
    onRefreshData,
    filter = "DAY",
    hasData = false,
    refreshInterval = 10000,
    onFilterChange = () => {
        /*DEFAULT FUNCTION*/
    },
}: IApexLineChartIntegration) => {
    React.useEffect(() => {
        const interval = setInterval(() => {
            onRefreshData && onRefreshData();
        }, refreshInterval);

        return () => clearInterval(interval);
    });

    return (
        <GraphTitleWrapper
            key={name}
            title={name}
            filter={filter}
            hasData={hasData}
            variant="subtitle1"
            handleFilterChange={onFilterChange}
        >
            <ApexLineChart
                key={name}
                name={name}
                dataList={data}
                range={TIME_RANGE_IN_MILLISECONDS}
            />
        </GraphTitleWrapper>
    );
};

export default ApexLineChartIntegration;
