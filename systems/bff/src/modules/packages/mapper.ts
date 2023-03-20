import { IPackageMapper } from "./interface";
import {
    PackageAPIResDto,
    PackageDto,
    PackagesAPIResDto,
    PackagesResDto,
} from "./types";

class PackageMapper implements IPackageMapper {
    dtoToPackagesDto(res: PackagesAPIResDto): PackagesResDto {
        const packages: PackageDto[] = [];
        res.packages.forEach(p => {
            packages.push({
                uuid: p.uuid,
                name: p.name,
                orgId: p.org_id,
                active: p.active,
                duration: p.duration,
                simType: p.sim_type,
                createdAt: p.created_at,
                deletedAt: p.deleted_at,
                updatedAt: p.updated_at,
                smsVolume: p.sms_volume,
                dataVolume: p.data_volume,
                voiceVolume: p.voice_volume,
                orgRatesId: p.org_rates_id,
            });
        });
        return {
            packages: packages,
        };
    }
    dtoToPackageDto(res: PackageAPIResDto): PackageDto {
        return {
            uuid: res.package.uuid,
            name: res.package.name,
            orgId: res.package.org_id,
            active: res.package.active,
            duration: res.package.duration,
            simType: res.package.sim_type,
            createdAt: res.package.created_at,
            deletedAt: res.package.deleted_at,
            updatedAt: res.package.updated_at,
            smsVolume: res.package.sms_volume,
            dataVolume: res.package.data_volume,
            voiceVolume: res.package.voice_volume,
            orgRatesId: res.package.org_rates_id,
        };
    }
}
export default <IPackageMapper>new PackageMapper();
