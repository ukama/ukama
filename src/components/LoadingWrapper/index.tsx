import { SkeletonRoundedCard } from "../../styles";

const LoadingWrapper = (props: any) => {
    const {
        children,
        isLoading,
        width = "100%",
        height = "100%",
        variant = "rectangular",
    } = props;
    if (isLoading)
        return (
            <SkeletonRoundedCard
                width={width}
                height={height}
                variant={variant}
            />
        );

    return <div style={{ height: "100%" }}>{children}</div>;
};

export default LoadingWrapper;
