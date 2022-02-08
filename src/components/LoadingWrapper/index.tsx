import React from "react";
import { SkeletonRoundedCard } from "../../styles";

interface ILoadingWrapper {
    width?: string | number;
    height?: string | number;
    children: React.ReactNode;
    radius?: "small" | "medium";
    isLoading: boolean | undefined;
    variant?: "text" | "rectangular" | "circular";
}

const LoadingWrapper = ({
    children,
    width = "100%",
    height = "100%",
    radius = "medium",
    variant = "rectangular",
    isLoading = false,
}: ILoadingWrapper) => {
    const borderRadius = radius === "medium" ? "10px" : "4px";
    if (isLoading)
        return (
            <SkeletonRoundedCard
                width={width}
                height={height}
                variant={variant}
                sx={{ borderRadius: borderRadius }}
            />
        );

    return (
        <div style={{ height: "inherit", width: "inherit" }}>{children}</div>
    );
};

export default LoadingWrapper;
