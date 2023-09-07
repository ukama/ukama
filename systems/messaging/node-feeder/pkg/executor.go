package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Device5xxServerError struct {
	error
}

type Device4xxServerError struct {
	error
}

type RequestExecutor interface {
	Execute(req *DevicesUpdateRequest) error
}

type requestExecutor struct {
	nodeResolver      NodeIpResolver
	deviceNetworkConf *DeviceNetworkConfig
}

func NewRequestExecutor(deviceNet NodeIpResolver, deviceNetworkConf *DeviceNetworkConfig) *requestExecutor {
	return &requestExecutor{nodeResolver: deviceNet, deviceNetworkConf: deviceNetworkConf}
}

func (e *requestExecutor) Execute(req *DevicesUpdateRequest) error {
	segs := strings.Split(req.Target, ".")
	if len(segs) != 2 {
		return fmt.Errorf("invalid target format")
	}
	nodeId, err := ukama.ValidateNodeId(segs[1])
	if err != nil {
		return errors.Wrap(err, "invalid node id format")
	}
	ip, err := e.nodeResolver.Resolve(nodeId)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			logrus.Warningf("node %s not found", nodeId)
			return nil
		}

		return errors.Wrap(err, "error resolving device ip for device "+segs[1])
	}
	logrus.Infof("resolved device ip: %s", ip)

	strUrl := fmt.Sprintf("http://%s/%s", ip, strings.Trim(req.Path, "/"))
	if e.deviceNetworkConf.Port != 0 {
		strUrl = fmt.Sprintf("http://%s:%d/%s", ip, e.deviceNetworkConf.Port, strings.Trim(req.Path, "/"))
	}
	u, err := url.Parse(strUrl)
	if err != nil {
		return errors.Wrap(err, "malformed url")
	}

	c := http.Client{
		Timeout: time.Duration(e.deviceNetworkConf.TimeoutSeconds) * time.Second,
	}

	logrus.Infof("sending request to %s", u.String())
	resp, err := c.Do(&http.Request{
		Body:   io.NopCloser(strings.NewReader(req.Body)),
		Method: req.HttpMethod,
		URL:    u,
	})
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}

	logrus.Infof("Response status: %d", resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Warning("error reading response body ", err)
	} else {
		logrus.Debugf("Response body: %s", string(b))
	}

	// request with >500 status code considered as filed
	if resp.StatusCode >= 500 {
		return Device5xxServerError{
			fmt.Errorf("server error: %d", resp.StatusCode),
		}
	} else if resp.StatusCode >= 400 {
		return Device4xxServerError{
			fmt.Errorf("server error: %d", resp.StatusCode),
		}
	}

	if err != nil {
		return errors.Wrap(err, "error sending request")
	}

	return nil
}
