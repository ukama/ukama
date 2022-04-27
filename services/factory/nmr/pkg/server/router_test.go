package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/services/factory/nmr/mocks"
	"github.com/ukama/ukama/services/factory/nmr/pkg"

	"github.com/ukama/ukama/services/factory/nmr/internal/db"
	"github.com/ukama/ukama/services/factory/nmr/pkg/router"
)

func init() {
	pkg.IsDebugMode = true
}

var defaultCongif = &pkg.Config{
	Server: rest.HttpConfig{
		Cors: cors.Config{
			AllowAllOrigins: true,
		},
	},
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func NewNode(id string) *db.Node {
	return &db.Node{
		NodeID:        id,
		Type:          "hnode",
		PartNumber:    "a1",
		Skew:          "s1",
		Mac:           "00:01:02:03:04:05",
		SwVersion:     "1.1",
		OemName:       "ukama",
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenrated",
	}
}

func NewModule(id string) *db.Module {
	return &db.Module{
		ModuleID:   id,
		Type:       "TRX",
		PartNumber: "a1",
		HwVersion:  "s1",
		Mac:        "00:01:02:03:04:05",
		SwVersion:  "1.1",
		MfgName:    "ukama",
		Status:     "StatusLabelGenrated",
	}
}

func Test_PutNode(t *testing.T) {
	// Arrange
	nodeId := "1001"
	node := NewNode(nodeId)

	body, _ := json.Marshal(node)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/node/?node=1001&looking_to=update", bytes.NewReader(body))

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	nodeRepo.On("AddOrUpdateNode", node).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_GetNode(t *testing.T) {
	t.Run("Read node", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		node := NewNode(nodeId)

		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/node/?node=1001&looking_for=*", nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("GetNode", nodeId).Return(node, nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), nodeId)
	})

}

func Test_DeleteNode(t *testing.T) {
	t.Run("Delete node", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/node/?node=1001&looking_to=*", nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("DeleteNode", nodeId).Return(nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

	})

}

func Test_PutNodeStatus(t *testing.T) {
	// Arrange
	nodeId := "1001"
	status := "StatusLabelGenrated"

	url := "/node/status?node=" + nodeId + "&looking_to=update&status=" + status
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	mfgStatus, err := db.MfgState(status)
	assert.NoError(t, err, nil)

	nodeRepo.On("UpdateNodeStatus", nodeId, *mfgStatus).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_GetNodeStatus(t *testing.T) {
	t.Run("Read node status", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		status := db.StatusLabelGenrated

		url := "/node/status?node=" + nodeId + "&looking_for=*"

		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("GetNodeStatus", nodeId).Return(&status, nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), status)

	})

}

func Test_PutNodeMfgTestStatus(t *testing.T) {
	// Arrange
	nodeId := "1001"

	jReq := "{ \"mfgTestStatus\" : \"MfgTestStatusUnderTest\", \"mfgReport\" : \"production test pass\", \"status\": \"StatusModuleTest\" }"

	url := "/node/mfgstatus?node=" + nodeId + "&looking_to=update"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, strings.NewReader(jReq))

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	nodeRepo.On("UpdateNodeMfgTestStatus", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_PutNodeMfgTestStatusFail(t *testing.T) {
	// Arrange
	nodeId := "1001"

	jReq := "{ \"mfgTestStatus\" : \"testing\", \"mfgReport\" : \"production test pass\", \"status\": \"StatusModuleTest\" }"

	url := "/node/mfgstatus?node=" + nodeId + "&looking_to=update"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, strings.NewReader(jReq))

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	nodeRepo.On("UpdateNodeMfgTestStatus", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "is invalid mfg test status")

}

func Test_PutNodeMfgTestStatusFail2(t *testing.T) {
	// Arrange
	nodeId := "1001"

	jReq := "{ \"mfgTestStatus\" : \"MfgTestStatusPending\", \"mfgReport\" : \"production test pass\", \"status\": \"StatusModuleUnkown\" }"

	url := "/node/mfgstatus?node=" + nodeId + "&looking_to=update"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, strings.NewReader(jReq))

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	nodeRepo.On("UpdateNodeMfgTestStatus", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "is invalid mfg status")

}

func Test_GetNodeMfgTestStatus(t *testing.T) {
	t.Run("Read node mfg status", func(t *testing.T) {
		// Arrange
		nodeId := "1001"
		status := db.MfgTestStatusPass

		mfg := []byte("\"report: passed\"")

		url := "/node/mfgstatus?node=" + nodeId + "&looking_for=*"

		//body, _ := json.Marshal(node)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		nodeRepo.On("GetNodeMfgTestStatus", nodeId).Return(&status, &mfg, nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), status.String())

	})

}

func Test_PutModule(t *testing.T) {
	// Arrange
	moduleId := "1001"
	module := ReqAddOrUpdateModule{
		ModuleID:   moduleId,
		Type:       "TRX",
		PartNumber: "a1",
		HwVersion:  "s1",
		Mac:        "00:01:02:03:04:05",
		SwVersion:  "1.1",
		MfgName:    "ukama",
		Status:     "StatusLabelGenrated",
		UnitID:     "0001",
	}
	url := "/module/?module=" + moduleId + "&looking_to=update"
	body, _ := json.Marshal(module)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	moduleRepo.On("UpsertModule", mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_GetModule(t *testing.T) {
	// Arrange
	moduleId := "1001"
	module := &db.Module{}
	url := "/module/?module=" + moduleId + "&looking_for=*"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	moduleRepo.On("GetModule", mock.Anything).Return(module, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_DeleteModule(t *testing.T) {
	t.Run("Delete Module", func(t *testing.T) {
		// Arrange
		moduleId := "1001"

		url := "/module/?module=" + moduleId + "&looking_to=*"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", url, nil)

		nodeRepo := mocks.NodeRepo{}
		moduleRepo := mocks.ModuleRepo{}
		rs := router.RouterServer{}

		r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

		moduleRepo.On("DeleteModule", moduleId).Return(nil)

		// act
		r.ServeHTTP(w, req)

		// assert
		assert.Equal(t, 200, w.Code)

	})

}

func Test_PutAssignModule(t *testing.T) {
	// Arrange
	moduleId := "M1001"
	nodeId := "N1001"

	url := "/module/assign?module=" + moduleId + "&looking_to=update" + "&node=" + nodeId

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	moduleRepo.On("UpdateNodeId", moduleId, nodeId).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_PutModuleMfgStatus(t *testing.T) {
	// Arrange
	moduleId := "1001"
	status := "StatusLabelGenrated"

	url := "/module/status?module=" + moduleId + "&looking_to=update&status=" + status

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	mfgtestStatus, err := db.MfgState(status)
	assert.NoError(t, err, nil)

	moduleRepo.On("UpdateModuleMfgStatus", moduleId, *mfgtestStatus).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_PutModuleMfgStatusFail(t *testing.T) {
	// Arrange
	moduleId := "1001"
	status := "StatusLabelUnkown"

	url := "/module/status?module=" + moduleId + "&looking_to=update&status=" + status

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	moduleRepo.On("UpdateModuleMfgStatus", moduleId, mock.Anything).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "is invalid mfg status")

}

func Test_GetModuleMfgStatus(t *testing.T) {
	// Arrange
	moduleId := "1001"
	var status db.MfgStatus
	url := "/module/status?module=" + moduleId + "&looking_for=*"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	moduleRepo.On("GetModuleMfgStatus", moduleId).Return(&status, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}

func Test_PutModuleMfgField(t *testing.T) {
	// Arrange
	moduleId := "1001"
	fieldList := [...]string{"bootstrap_cert", "user_config", "factory_config", "user_calibration", "factory_calibration", "inventory_data"}

	for _, field := range fieldList {
		t.Run(field, func(t *testing.T) {
			data := []byte(field)
			url := "/module/field?module=" + moduleId + "&looking_to=update&field=" + field

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))

			nodeRepo := mocks.NodeRepo{}
			moduleRepo := mocks.ModuleRepo{}
			rs := router.RouterServer{}

			r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

			arg, err := db.GetModuleDataFieldName(field)
			assert.NoError(t, err)

			moduleRepo.On("UpdateModuleMfgField", moduleId, *arg, mock.Anything).Return(nil)

			// act
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, 200, w.Code)
		})
	}

}

func Test_GetModuleMfgField(t *testing.T) {
	moduleId := "1001"
	fieldList := [...]string{"bootstrap_cert", "user_config", "factory_config", "user_calibration", "factory_calibration", "inventory_data"}

	for _, field := range fieldList {
		t.Run(field, func(t *testing.T) {
			module := db.Module{
				ModuleID: moduleId,
			}

			var columnName string
			byteBody := []byte(field)

			switch field {
			case "mfg_report":
				columnName = "mfg_report"
				module.MfgReport = &byteBody

			case "bootstrap_cert":
				columnName = "bootstrap_certs"
				module.BootstrapCerts = &byteBody

			case "user_config":
				columnName = "user_config"
				module.UserConfig = &byteBody
			case "factory_config":
				columnName = "factory_config"
				module.FactoryConfig = &byteBody

			case "user_calibration":
				columnName = "user_calibration"
				module.UserCalibration = &byteBody

			case "factory_calibration":
				columnName = "factory_calibration"
				module.FactoryCalibration = &byteBody

			// case "cloud_certs":
			// 	columnName = "cloud_certs"
			// 	module, err = r.moduleRepo.GetModuleMfgField(req.ModuleID, req.Field)
			// 	if err != nil {
			// 		return nil, rest.HttpError{
			// 			HttpCode: http.StatusNotFound,
			// 			Message:  err.Error(),
			// 		}
			// 	}
			// 	data = module.CloudCerts

			case "inventory_data":
				columnName = "inventory_data"
				module.InventoryData = &byteBody

			default:
				columnName = ""

			}

			assert.NotEqual(t, columnName, "")

			url := "/module/field?module=" + moduleId + "&looking_for=update&field=" + field

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)

			nodeRepo := mocks.NodeRepo{}
			moduleRepo := mocks.ModuleRepo{}
			rs := router.RouterServer{}

			r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

			moduleRepo.On("GetModuleMfgField", moduleId, columnName).Return(&module, nil)

			// act
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, 200, w.Code)
		})
	}

}

func Test_DeleteBootstrapCerts(t *testing.T) {
	// Arrange
	moduleId := "1001"

	url := "/module/bootstrapcerts?module=" + moduleId + "&looking_to=update"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)

	nodeRepo := mocks.NodeRepo{}
	moduleRepo := mocks.ModuleRepo{}
	rs := router.RouterServer{}

	r := NewRouter(defaultCongif, &rs, &nodeRepo, &moduleRepo).fizz.Engine()

	moduleRepo.On("DeleteBootstrapCert", moduleId).Return(nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

}
