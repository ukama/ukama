import {
    PackageAPIResDto,
    PackageDto,
    PackagesAPIResDto,
    PackagesResDto,
} from "./types";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface IPackageService {}

export interface IPackageMapper {
    dtoToPackagesDto(res: PackagesAPIResDto): PackagesResDto;
    dtoToPackageDto(res: PackageAPIResDto): PackageDto;
}
