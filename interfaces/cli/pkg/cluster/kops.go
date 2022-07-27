package cluster

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	goerr "errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/services/common/errors"
)

const KOPS_BIN = "kops"

type KopsWrapper struct {
	log             pkg.Logger
	stateBucketName string // s3 bucket name without s3://
	debugMode       bool
}

func NewKopsWrapper(log pkg.Logger, stateBucketName string, debugMode bool) *KopsWrapper {
	return &KopsWrapper{
		log:             log,
		stateBucketName: stateBucketName,
		debugMode:       debugMode,
	}
}

func (k *KopsWrapper) createBucket(name string, region string) error {
	k.log.Printf("Creating bucket %s", name)
	if len(name) == 0 {
		return fmt.Errorf("bucket name is empty")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to load aws config")
	}

	// use us-east-1 for all cluster states
	cfg.Region = "us-east-1"

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})

	var bne *types.BucketAlreadyExists
	if goerr.As(err, &bne) {
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "failed to create bucket")
	}

	k.log.Debugf("Bucket %s created", name)
	return nil
}

func (k *KopsWrapper) ProvisionAwsCluster(dnsName string, region string) error {

	// kops create cluster --dns public  kopscluster.k8s.local   --zones "us-east-1a,us-east-1b" --cloud aws  --state s3://kops-bucket-test-denis
	if len(dnsName) == 0 {
		return fmt.Errorf("dnsName is required")
	}
	if len(region) == 0 {
		return fmt.Errorf("zones is required")
	}

	if len(k.stateBucketName) == 0 {
		return fmt.Errorf("kops state bucket name is required")
	}

	zones := strings.Join([]string{region + "a", region + "b"}, ",")

	err := k.createBucket(k.stateBucketName, region)
	if err != nil {
		return errors.Wrap(err, "failed to create or access bucket")
	}

	k.log.Printf("Provisioning cluster %s\n", dnsName)

	err = k.runKopsCmd("create", "cluster", "--yes", "--dns", "public", dnsName, "--zones", zones, "--cloud", "aws", "--state", k.s3BucketFullUrl())

	if err != nil {
		return errors.Wrap(err, "failed to create cluster")
	}
	k.log.Printf("Waiting for cluster %s to come up", dnsName)

	err = k.Validate(dnsName, true)
	if err != nil {
		return errors.Wrap(err, "failed to validate cluster")
	}

	return nil
}

func (k *KopsWrapper) runKopsCmd(args ...string) error {
	if k.debugMode {
		args = append(args, "-v", "9")
	}
	cmd := exec.Command(KOPS_BIN, args...)
	cmd.Stdout = k.log.Stdout()
	cmd.Stderr = k.log.Stderr()
	k.log.Debugf("Running command: %s", cmd.String())
	return cmd.Run()
}

func (k *KopsWrapper) s3BucketFullUrl() string {
	return "s3://" + k.stateBucketName
}

func (k *KopsWrapper) Validate(name string, wait bool) error {
	k.log.Printf("Validating cluster %s", name)
	// kops validate cluster --wait 10m --name kopscluster.dev.ukama.com --state s3://kops-bucket-test-denis
	args := []string{"validate", "cluster", "--name", name, "--state", k.s3BucketFullUrl()}
	if wait {
		args = append(args, "--wait", "15m")
	}
	k.log.Printf("Validating cluster %s", name)
	err := k.runKopsCmd(args...)
	if err != nil {
		return errors.Wrap(err, "failed to validate cluster")
	}
	return nil
}

func (k *KopsWrapper) Delete() {
	// kops delete cluster  --name myfirstcluster.dev.ukama.com --yes --state s3://kops-bucket-test-denis
}

func (k *KopsWrapper) UpdateKubeConfig() {
	//  kops update cluster --name kopscluster.k8s.local --yes --admin --state s3://kops-bucket-test-denis
}
