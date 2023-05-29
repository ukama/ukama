import { ISimMapper } from "./interface";
import {
    SimAPIResDto,
    SimDetailsDto,
    SimDto,
    SimsAPIResDto,
    SimsResDto,
} from "./types";

class SimMapper implements ISimMapper {
    dtoToSimResDto = (res: SimAPIResDto): SimDto => {
        return {
            activationCode: res.sim.activation_code,
            createdAt: res.sim.created_at,
            iccid: res.sim.iccid,
            id: res.sim.id,
            isAllocated: res.sim.is_allocated,
            isPhysical: res.sim.is_physical,
            msisdn: res.sim.msisdn,
            qrCode: res.sim.qr_code,
            simType: res.sim.sim_type,
            smapAddress: res.sim.sm_ap_address,
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
    dtoToSimsDto = (res: SimsAPIResDto): SimsResDto => {
        const sims: SimDto[] = [];
        for (const sim of res.sims) {
            sims.push(this.dtoToSimResDto({ sim: sim }));
        }
        return {
            sim: sims,
        };
    };
}
export default <ISimMapper>new SimMapper();
