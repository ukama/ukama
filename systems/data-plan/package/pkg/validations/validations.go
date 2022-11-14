package validationErrors

import "github.com/ukama/ukama/systems/data-plan/base-rate/pb"

func IsRequestEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func ReqStrTopb(e string) pb.SimType {
	switch e {
	case "inter_none":
		return pb.SimType_INTER_NONE
	case "inter_mno_data":
		return pb.SimType_INTER_MNO_DATA
	case "inter_ukama_all":
		return pb.SimType_INTER_UKAMA_ALL
	case "inter_mno_all":
		return pb.SimType_INTER_MNO_ALL
	default:
		return pb.SimType_INTER_NONE
	}
}
func ReqPbToStr(s pb.SimType) string {
	switch s {
	case pb.SimType_INTER_NONE:
		return "inter_none"
	case pb.SimType_INTER_MNO_DATA:
		return "inter_mno_data"
	case pb.SimType_INTER_MNO_ALL:
		return "inter_mno_all"
	case pb.SimType_INTER_UKAMA_ALL:
		return "inter_ukama_all"
	default:
		return "inter_none"
	}
}
