package server

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/pb"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/net/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strings"
)

type DnsServer struct {
	nns    pkg.NnsReader
	config *pkg.DnsConfig
}

func NewDnsServer(nns pkg.NnsReader, config *pkg.DnsConfig) *DnsServer {
	return &DnsServer{
		nns:    nns,
		config: config,
	}
}

func (d *DnsServer) Query(ctx context.Context, p *pb.DnsPacket) (*pb.DnsPacket, error) {
	m := new(dns.Msg)
	if err := m.Unpack(p.Msg); err != nil {
		return nil, fmt.Errorf("failed to unpack msg: %v", err)
	}
	r := new(dns.Msg)
	r.SetReply(m)
	r.Authoritative = true

	for _, q := range r.Question {
		hdr := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass}
		switch q.Qtype {
		case dns.TypeA:
			ip, err := d.resolveIp(ctx, q.Name)
			if err != nil {
				return nil, err
			}
			r.Answer = append(r.Answer, &dns.A{
				Hdr: hdr,
				A:   ip})
		case dns.TypeAAAA:
			return nil, status.Error(codes.NotFound, "No AAAA record found")
		default:
			return nil, fmt.Errorf("only A and AAAA supported, got qtype=%d", q.Qtype)
		}
	}

	if len(r.Answer) == 0 {
		r.Rcode = dns.RcodeNameError
	}

	out, err := r.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack msg: %v", err)
	}
	return &pb.DnsPacket{Msg: out}, nil
}

func (d *DnsServer) resolveIp(ctx context.Context, name string) (net.IP, error) {
	baseDomain := strings.ToLower(d.config.NodeDomain) + "."
	nodeId := strings.ToLower(name)

	if !strings.HasSuffix(nodeId, baseDomain) {
		logrus.Errorf("no subdomain")
		return nil, status.Error(codes.NotFound, "no subdomain")
	}

	nodeId = strings.TrimSuffix(nodeId, baseDomain)
	nodeId = strings.TrimSuffix(nodeId, ".")
	logrus.Debugln("Resolving: ", name)
	if len(nodeId) == 0 {
		logrus.Errorf("nodeId is empty")
		return nil, status.Error(codes.NotFound, "nodeId is empty")
	}

	ndIpResp, err := d.nns.Get(ctx, nodeId)
	if err != nil {
		code := status.Code(err)
		if code == codes.NotFound || code == codes.InvalidArgument {
			logrus.Info("NodeID not found. Error:", err.Error())
			return nil, status.Error(codes.NotFound, err.Error())
		}

		logrus.Errorf("failed to get node: %v", err)
		return nil, err
	}
	ip := net.ParseIP(ndIpResp)
	if ip == nil {
		logrus.Errorf("failed to parse ip from DB: %v", err)
		return nil, err
	}
	return ip, nil
}
