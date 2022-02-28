import React from "react";
import Chart from "react-apexcharts";
import { format } from "date-fns";
import { GraphTitleWrapper } from "..";
import { getRandomData } from "../../utils";

const TIME_RANGE_IN_MILLISECONDS = 1 * 300;

interface IApexLineChartIntegration {
    name: string;
    initData: any;
    hasData: boolean;
    isStatic?: boolean;
    refreshInterval?: number;
    onRefreshData?: Function;
}

const ApexLineChart = (props: any) => {
    const options: any = {
        grid: {
            padding: {
                left: 30, // or whatever value that works
                right: 30, // or whatever value that works
            },
        },
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
                // offsetX: 14,
                // offsetY: -5,
                formatter: (val: any) => val.toFixed(2),
            },
            // tooltip: {
            //     enabled: true,
            // },
            tickAmount: 10,
        },
    };
    return (
        <Chart
            type="line"
            options={options}
            series={props.dataList}
            height={"300px"}
        />
    );
};

const ApexLineChartIntegration = ({
    name,
    initData = [],
    onRefreshData,
    isStatic = true,
    hasData = false,
    refreshInterval = 10000,
}: IApexLineChartIntegration) => {
    const nameList = [name];
    const defaultDataList = nameList.map(name => ({
        name: name,
        data: initData,
    }));
    const [dataList, setDataList] = React.useState(defaultDataList);

    React.useEffect(() => {
        const interval = setInterval(() => {
            if (isStatic) {
                setDataList(
                    dataList.map(val => {
                        return {
                            name: val.name,
                            data: [...val.data, ...getRandomData()],
                        };
                    })
                );
            } else
                onRefreshData &&
                    onRefreshData()?.then((res: any) => {
                        if (res && res?.data && res.data?.getMetricsCpuTRX)
                            setDataList(
                                dataList.map(val => {
                                    return {
                                        name: val.name,
                                        data: [
                                            ...val.data,
                                            ...res.data?.getMetricsCpuTRX,
                                        ],
                                    };
                                })
                            );
                    });
        }, refreshInterval);

        return () => clearInterval(interval);
    });

    return (
        <GraphTitleWrapper hasData={hasData} variant="subtitle1" title={name}>
            <ApexLineChart
                dataList={dataList}
                range={TIME_RANGE_IN_MILLISECONDS}
            />
        </GraphTitleWrapper>
    );
};

export default ApexLineChartIntegration;
