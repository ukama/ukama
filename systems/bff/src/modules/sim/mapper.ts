import { ISimMapper } from "./interface";
import { SimAPIResDto, SimResDto } from "./types";

class SimMapper implements ISimMapper {
    dtoToSimResDto = (res: SimAPIResDto): SimResDto => {
        return {
            activationCode: res.sim.activationCode,
            createdAt: res.sim.createdAt,
            iccid: res.sim.iccid,
            id: res.sim.id,
            isAllocated: res.sim.isAllocated,
            isPhysical: res.sim.isPhysical,
            msisdn: res.sim.msisdn,
            qrCode: res.sim.qrCode,
            simType: res.sim.simType,
            smDpAddress: res.sim.smDpAddress,
        };
    };
}
export default <ISimMapper>new SimMapper();
