import { CenterContainer, RoundedCard } from "../../styles";
import { NodePlaceholder } from "../../assets/images";
import { LoadingWrapper } from "..";

interface INodeDetailsCard {
    loading: boolean;
}

const NodeDetailsCard = ({ loading }: INodeDetailsCard) => {
    return (
        <LoadingWrapper radius={"small"} height={450} isLoading={loading}>
            <RoundedCard
                sx={{
                    borderRadius: "4px",
                    height: "fit-content",
                }}
            >
                <CenterContainer>
                    <img src={NodePlaceholder} width="100%" />
                </CenterContainer>
            </RoundedCard>
        </LoadingWrapper>
    );
};

export default NodeDetailsCard;
