const colors = [
    { border: "rgb(53, 162, 235)", background: "rgba(53, 162, 235, 0.5)" },
    { border: "rgb(255, 99, 132)", background: "rgba(255, 99, 132, 0.5)" },
];

const LineGraphtypes = [
    { cubicInterpolationMode: "monotone" },
    { borderDash: [8, 4] },
];

const getDatasetsWRTLabels = (labels: string[]) => {
    const dataset = labels.map((item, index) => {
        return {
            data: [],
            label: item,
            ...LineGraphtypes[index],
            borderColor: colors[index].border,
            backgroundColor: colors[index].background,
        };
    });

    return {
        datasets: dataset,
    };
};

export { getDatasetsWRTLabels };
