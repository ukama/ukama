import { ResponseProps } from "../types";
import { DEFAULT_RESPONSE } from "../constants";
import { useCallback, useEffect, useRef, useState } from "react";

const useWhoami = () => {
    const isCurrent = useRef(true);
    const [responses, setResponses] = useState<ResponseProps>(DEFAULT_RESPONSE);

    useEffect(() => {
        return () => {
            isCurrent.current = false;
        };
    }, []);

    const runQuery = useCallback(async () => {
        setResponses(() => ({ ...responses, loading: true }));
        return fetch(
            `${process.env.REACT_APP_KRATOS_BASE_URL}/sessions/whoami`,
            {
                credentials: "include",
                headers: {
                    Accept: "application/json",
                    "Content-Type": "application/json",
                },
            }
        )
            .then(response => response.json())
            .then(res => {
                if (res?.identity?.id) {
                    setResponses((prev: ResponseProps) => ({
                        ...prev,
                        response: { isValid: true },
                        error: "",
                        loading: false,
                    }));
                } else {
                    setResponses((prev: ResponseProps) => ({
                        ...prev,
                        response: { isValid: false },
                        error: "Unauthorized",
                        loading: false,
                    }));
                }
            })
            .catch(() =>
                setResponses((prev: ResponseProps) => ({
                    ...prev,
                    response: { isValid: false },
                    error: "WhoamiError",
                    loading: false,
                }))
            );
    }, []);

    useEffect(() => {
        runQuery();
    }, [runQuery]);

    return { ...responses, refetch: runQuery };
};

export default useWhoami;
