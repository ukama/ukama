import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
} from "chart.js";
import "chartjs-adapter-luxon";
import { Line } from "react-chartjs-2";
import zoomPlugin from "chartjs-plugin-zoom";
import ChartStreaming from "chartjs-plugin-streaming";
import { GraphTitleWrapper } from "..";

ChartJS.register(
    zoomPlugin,
    ChartStreaming,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

const data: any = {
    datasets: [
        {
            label: "Current Temperature",
            borderColor: "rgb(53, 162, 235)",
            backgroundColor: "rgba(53, 162, 235, 0.5)",
            cubicInterpolationMode: "monotone",
            data: [],
        },
        {
            label: "Threshold ~85",
            borderColor: "rgb(255, 99, 132)",
            backgroundColor: "rgba(255, 99, 132, 0.5)",
            borderDash: [8, 4],
            data: [],
        },
    ],
};

const onRefresh = (chart: any) => {
    const now = Date.now();
    chart.data.datasets.forEach((dataset: any) => {
        dataset.data.push({
            x: now,
            y: Math.random() * (20 - -50) + -50,
        });
    });
};

const zoomOptions = {
    pan: {
        enabled: true,
        mode: "x",
    },
    zoom: {
        pinch: {
            enabled: true,
        },
        wheel: {
            enabled: true,
        },
        mode: "x",
    },
    limits: {
        x: {
            minDelay: 0,
            maxDelay: 4000,
            minDuration: 10000, //ZoomIn
            maxDuration: 80000, //ZoomOut
        },
    },
};

const config: any = {
    type: "line",
    data: data,
    options: {
        scales: {
            x: {
                type: "realtime",
                realtime: {
                    duration: 20000,
                    refresh: 1000,
                    delay: 2000,
                    onRefresh: onRefresh,
                },
            },
            y: {
                title: {
                    display: false,
                    text: "Celsius",
                },
            },
        },
        interaction: {
            intersect: false,
        },
        plugins: {
            annotation: false,
            datalabels: false,
            zoom: zoomOptions,
            legend: { position: "bottom" },
        },
    },
};

interface ILineChart {
    title: string;
    filter?: string;
    hasData?: boolean;
    showFilter?: boolean;
    handleFilterChange?: Function;
}

const LineChart = ({
    title,
    filter,
    hasData,
    showFilter = true,
    handleFilterChange,
}: ILineChart) => (
    <GraphTitleWrapper
        title={title}
        filter={filter}
        hasData={hasData}
        variant="subtitle1"
        showFilter={showFilter}
        handleFilterChange={handleFilterChange}
    >
        <Line
            key={title}
            id={title}
            height={80}
            data={data}
            options={{ ...config.options, id: title }}
        />
    </GraphTitleWrapper>
);

export default LineChart;
