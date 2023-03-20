import { ISimMapper } from "./interface";
import { SimAPIResDto, SimResDto, SimDetailsDto } from "./types";

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
    dtoToSimDetailsDto(response: any): SimDetailsDto {
        const {
            id,
            subscriberId,
            networkId,
            orgId,
            Package,
            iccid,
            msisdn,
            imsi,
            type,
            status,
            isPhysical,
            firstActivatedOn,
            lastActivatedOn,
            activationsCount,
            deactivationsCount,
            allocatedAt,
        } = response;

        return {
            id,
            subscriberId,
            networkId,
            orgId,
            Package,
            iccid,
            msisdn,
            imsi,
            type,
            status,
            isPhysical,
            firstActivatedOn: firstActivatedOn?.toDate(),
            lastActivatedOn: lastActivatedOn?.toDate(),
            activationsCount,
            deactivationsCount,
            allocatedAt: allocatedAt?.toDate(),
        };
    }
}
export default <ISimMapper>new SimMapper();
