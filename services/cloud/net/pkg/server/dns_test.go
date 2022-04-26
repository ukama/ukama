package server

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/pb"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/cloud/net/mocks"
	"github.com/ukama/ukamaX/cloud/net/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestDnsServer_Query(t *testing.T) {

	config := &pkg.DnsConfig{NodeDomain: "test.node"}

	const ip = "192.168.0.1"

	type nnsReturn struct {
		ip  string
		err error
	}
	tests := []struct {
		name    string
		request string
		nnsVals nnsReturn
		wantErr func(e error) bool
	}{
		{
			name:    "SuccessfulReposnse",
			nnsVals: nnsReturn{ip: ip, err: nil},
			request: "uk-sa2203-hnode-a1-0a16.test.node",
		},
		{
			name:    "NoBaseDomain",
			nnsVals: nnsReturn{ip: "", err: status.Error(codes.NotFound, "")},
			request: "uk-sa2203-hnode-a1-0a16",
			wantErr: func(e error) bool { return status.Code(e) == codes.NotFound },
		},
		{
			name:    "NonodeIdeDoesNotExist",
			nnsVals: nnsReturn{ip: "", err: status.Error(codes.NotFound, "")},
			request: "uk-sa2203-hnode-a1-0a16.test.node",
			wantErr: func(e error) bool { return status.Code(e) == codes.NotFound },
		},
		{
			name:    "InternalError",
			nnsVals: nnsReturn{ip: "", err: fmt.Errorf("internal error")},
			request: "uk-sa2203-hnode-a1-0a16.test.node",
			wantErr: func(e error) bool { return e != nil },
		},
		{
			name:    "EmptyNodeId",
			nnsVals: nnsReturn{ip: "", err: nil},
			request: "test.node",
			wantErr: func(e error) bool { return status.Code(e) == codes.NotFound },
		},
		{
			name:    "MixedCaseNodeId",
			nnsVals: nnsReturn{ip: ip, err: nil},
			request: "uk-SA2203-HNOde-a1-0a16.test.node",
		},
		{
			name:    "UpperCaseBaseDomain",
			nnsVals: nnsReturn{ip: ip, err: nil},
			request: "uk-sa2203-hnode-a1-0a16.TEST.NODE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nns := mocks.NnsReader{}
			nns.On("Get", mock.Anything, mock.Anything).Return(tt.nnsVals.ip, tt.nnsVals.err)

			d := &DnsServer{
				nns:    &nns,
				config: config,
			}

			m := new(dns.Msg)
			m.SetQuestion(tt.request+".", dns.TypeA)
			b, err := m.Pack()
			if err != nil {
				assert.FailNow(t, "Failed to pack message")
			}

			got, err := d.Query(context.TODO(), &pb.DnsPacket{Msg: b})

			if tt.wantErr != nil {
				assert.True(t, tt.wantErr(err), "Expected error")
			} else {
				m := new(dns.Msg)
				err = m.Unpack(got.Msg)

				assert.NoError(t, err)
				assert.Equal(t, len(m.Answer), 1)
				assert.Equal(t, m.Rcode, dns.RcodeSuccess)
				r := m.Answer[0].(*dns.A)
				assert.Equal(t, ip, r.A.String())
			}

		})
	}
}
