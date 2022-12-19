package utils

import (
	pb "github.com/ukama/ukama/systems/subscriber/simPool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/db"
)

// type RawRates struct {
// 	Country string `csv:"Country"`
// 	Network string `csv:"Network"`
// 	Vpmn    string `csv:"VPMN"`
// 	Imsi    string `csv:"IMSI"`
// 	Sms_mo  string `csv:"SMS MO"`
// 	Sms_mt  string `csv:"SMS MT"`
// 	Data    string `csv:"Data"`
// 	X2g     string `csv:"2G"`
// 	X3g     string `csv:"3G"`
// 	X5g     string `csv:"5G"`
// 	Lte     string `csv:"LTE"`
// 	Lte_m   string `csv:"LTE-M"`
// 	Apn     string `csv:"APN"`
// }

// func FetchData(url string) ([]RawRates, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	content, _ := io.ReadAll(resp.Body)

// 	var r []RawRates
// 	errorStr := "invalid CSV file data"
// 	err = csvutil.Unmarshal(content, &r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(r) == 0 || validations.IsEmpty(r[0].Country) ||
// 		validations.IsEmpty(r[0].Network) ||
// 		validations.IsEmpty(r[0].Data) {
// 		return nil, errors.New(errorStr)
// 	}

// 	return r, nil
// }

func PbParseToModel(slice []*pb.AddSimPool) []db.SimPool {
	var simPool []db.SimPool
	// for _, value := range slice {
	// 	simPool = append(simPool, db.SimPool{

	// 	})
	// }
	return simPool
}

func ParseModelToPb(slice []*db.SimPool) []pb.SimPool {
	var simPool []pb.SimPool
	// for _, value := range slice {
	// 	simPool = append(simPool, db.SimPool{

	// 	})
	// }
	return simPool
}
