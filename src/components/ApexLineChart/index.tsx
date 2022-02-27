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
        chart: {
            minHeight: "200px",
            height: "100%",
            width: "100%",
            zoom: {
                enabled: false,
            },
            animations: {
                enabled: true,
                easing: "linear",
                dynamicAnimation: {
                    speed: 1000,
                },
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
                    onRefreshData().then((res: any) => {
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
