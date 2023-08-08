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
                to: p.to,
                apn: p.apn,
                dlbr: p.dlbr,
                from: p.from,
                name: p.name,
                rate: p.rate,
                type: p.type,
                ulbr: p.ulbr,
                uuid: p.uuid,
                orgId: p.org_id,
                active: p.active,
                amount: p.amount,
                markup: p.markup,
                country: p.country,
                ownerId: p.owner_id,
                simType: p.sim_type,
                currency: p.currency,
                dataUnit: p.data_unit,
                duration: p.duration,
                flatrate: p.flatrate,
                provider: p.provider,
                createdAt: p.created_at,
                deletedAt: p.deleted_at,
                smsVolume: p.sms_volume,
                updatedAt: p.updated_at,
                voiceUnit: p.voice_unit,
                dataVolume: p.data_volume,
                messageUnit: p.message_unit,
                voiceVolume: p.voice_volume,
            });
        });
        return {
            packages: packages,
        };
    }
    dtoToPackageDto(res: PackageAPIResDto): PackageDto {
        return {
            to: res.package.to,
            apn: res.package.apn,
            dlbr: res.package.dlbr,
            from: res.package.from,
            name: res.package.name,
            rate: res.package.rate,
            type: res.package.type,
            ulbr: res.package.ulbr,
            uuid: res.package.uuid,
            orgId: res.package.org_id,
            active: res.package.active,
            amount: res.package.amount,
            markup: res.package.markup,
            country: res.package.country,
            ownerId: res.package.owner_id,
            simType: res.package.sim_type,
            currency: res.package.currency,
            dataUnit: res.package.data_unit,
            duration: res.package.duration,
            flatrate: res.package.flatrate,
            provider: res.package.provider,
            createdAt: res.package.created_at,
            deletedAt: res.package.deleted_at,
            smsVolume: res.package.sms_volume,
            updatedAt: res.package.updated_at,
            voiceUnit: res.package.voice_unit,
            dataVolume: res.package.data_volume,
            messageUnit: res.package.message_unit,
            voiceVolume: res.package.voice_volume,
        };
    }
}
export default <IPackageMapper>new PackageMapper();
