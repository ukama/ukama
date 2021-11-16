import { INodeMapper } from "./interface";
import { NodeDto } from "./types";

class NodeMapper implements INodeMapper {
    dtoToDto = (data: NodeDto[]): NodeDto[] => {
        const nodes: NodeDto[] = [];

        for (let i = 0; i < data.length; i++) {
            if (data[i]) {
                const alert = {
                    id: data[i].id,
                    title: data[i].title,
                    description: data[i].description,
                    totalUser: data[i].totalUser,
                };
                nodes.push(alert);
            }
        }

        return nodes;
    };
}
export default <INodeMapper>new NodeMapper();
