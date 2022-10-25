package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/models"
)

var region = "us-east-1"
var bucket = "telco-rates"
var destinationFileName = "ratessss.csv"

type Server struct {
	H db.Handler
	pb.UnimplementedRatesServiceServer
}

type Result struct {
	Id           int64
	Country      string
	Network      string
	Vpmn         string
	Imsi         string
	Sms_mo       string
	Sms_mt       string
	Data         string
	X2g          string
	X3g          string
	X5g          string
	Lte          string
	Lte_m        string
	Apn          string
	Created_at   string
	Effective_at string
	End_at       string
}

func (s *Server) GetRates(ctx context.Context, req *pb.RatesRequest) (*pb.RatesResponse, error) {
	logrus.Infof("Get all rates  %v", req.GetCountry())

	var rate_list *pb.RatesResponse = &pb.RatesResponse{}

	if result := s.H.DB.Find(&rate_list.Rates); result.Error != nil {
		fmt.Println(result.Error)
	}

	rates := pb.Rate{}
	rate_list.Rates = append(rate_list.Rates, &rates)

	return rate_list, nil

}

func (s *Server) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {
	logrus.Infof("Get rate  %v", req.GetRateId())

	var rate models.Rate

	if result := s.H.DB.First(&rate, req.RateId); result.Error != nil {
		fmt.Println(result.Error)
	}

	data := &pb.Rate{
		Country:     rate.Country,
		Network:     rate.Network,
		Vpmn:        rate.Vpmn,
		Imsi:        rate.Imsi,
		SmsMo:       rate.Sms_mo,
		SmsMt:       rate.Sms_mt,
		Data:        rate.Data,
		X2G:         rate.X2g,
		X3G:         rate.X3g,
		Lte:         rate.Lte,
		LteM:        rate.Lte_m,
		Apn:         rate.Apn,
		CreatedAt:   rate.Created_at,
		EffectiveAt: rate.Effective_at,
		EndAt:       rate.End_at,
	}
	return &pb.RateResponse{
		Rate: data,
	}, nil
}

func (s *Server) UploadRates(ctx context.Context, req *pb.UploadRatesRequest) (*pb.UploadRatesResponse, error) {
	s3filePath := req.FilePath
	ratesApplicableFrom := req.EffectiveAt

	createFile(destinationFileName)
	err := retrieveFile(s3filePath, bucket, region, destinationFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Open(destinationFileName)
	check(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	check(err)

	query := createQuery(data, ratesApplicableFrom)

	deleteFile(destinationFileName)

	s.H.DB.Exec(query)

	var rate_list *pb.UploadRatesResponse = &pb.UploadRatesResponse{}

	if result := s.H.DB.Find(&rate_list.Rate); result.Error != nil {
		fmt.Println(result.Error)
	}

	rates := pb.Rate{}
	rate_list.Rate = append(rate_list.Rate, &rates)

	return rate_list, nil
}

//w TODO: Move these func to util

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createFile(fileName string) {
	f, err := os.Create(fileName)
	check(err)

	defer f.Close()
}

func deleteFile(fileName string) {
	e := os.Remove("ratessss.csv")
	check(e)
}

func retrieveFile(key string, bucket string, region string, destPath string) error {
	sess, err := session.NewSession(
		&aws.Config{Region: aws.String(region)},
	)
	if err != nil {
		return err
	}
	svc := s3.New(sess)
	params := &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}
	res, err := svc.GetObject(params)
	check(err)

	defer res.Body.Close()

	if destPath == "" {
		io.Copy(os.Stdout, res.Body)
		return nil
	}

	outFile, er := os.Create(destPath)
	check(er)
	defer outFile.Close()
	io.Copy(outFile, res.Body)

	return nil
}

func createQuery(rows [][]string, effective_at string) string {
	headerStr := ""
	valueStrings := make([]string, 0, len(rows))
	for i, row := range rows {
		if i == 0 {
			headerStr = "(" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
				strings.ReplaceAll(strings.ToLower(strings.Join(row[:], ",")), " ", "_"),
				"-", "_"), "2g", "X2g"), "3g", "X3g") + ",effective_at,end_at,X5g" + ")"
			headerStr = strings.ReplaceAll(strings.ReplaceAll(headerStr, "country_on_cronus,", ""), "network_id_on_cronus,", "")
			continue
		}
		values := row
		str := ""
		for j, value := range values {
			if j == 1 || j == 3 {
				continue
			} else {
				if j == len(values)-1 {
					str = str + "'" + strings.ReplaceAll(strings.ReplaceAll(value, "'", ""), ",", " ") + "'"
				} else {
					str = str + "'" + strings.ReplaceAll(strings.ReplaceAll(value, "'", ""), ",", "") + "', "
				}
			}
		}
		str = str + ", '" + effective_at + "', '', NULL"
		valueStrings = append(valueStrings, "("+strings.ReplaceAll(str, "''", "NULL")+")")
	}
	stmt := fmt.Sprintf("INSERT INTO rates %s VALUES %s", headerStr, strings.Join(valueStrings, ","))
	return stmt
}

//
