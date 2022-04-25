import { useRecoilValue } from "recoil";
import { isSkeltonLoading } from "../../recoil";
import { LoadingWrapper } from "../../components";

const Store = () => {
    const _isSkeltonLoading = useRecoilValue(isSkeltonLoading);
    return (
        <LoadingWrapper isLoading={_isSkeltonLoading} height={"90%"}>
            <div>
                <h1>Store Page:</h1>
            </div>
        </LoadingWrapper>
    );
};

export default Store;
