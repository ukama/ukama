import {
    SubscriberAPIResDto,
    SubscriberDto,
    SubscribersAPIResDto,
    SubscribersResDto,
} from "./types";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface ISubscriberService {}

export interface ISubscriberMapper {
    dtoToSubscriberResDto(res: SubscriberAPIResDto): SubscriberDto;
    dtoToSubscribersResDto(res: SubscribersAPIResDto): SubscribersResDto;
}
