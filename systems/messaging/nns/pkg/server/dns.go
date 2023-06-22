package server

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/coredns/pb"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DnsServer struct {
	nns    pkg.NnsReader
	config *pkg.DnsConfig
	pb.UnimplementedDnsServiceServer
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
		log.Errorf("no subdomain")
		return nil, status.Error(codes.NotFound, "no subdomain")
	}

	nodeId = strings.TrimSuffix(nodeId, baseDomain)
	nodeId = strings.TrimSuffix(nodeId, ".")
	log.Debugln("Resolving: ", name)
	if len(nodeId) == 0 {
		log.Errorf("nodeId is empty")
		return nil, status.Error(codes.NotFound, "nodeId is empty")
	}

	ndIpResp, err := d.nns.Get(ctx, nodeId)
	if err != nil {
		code := status.Code(err)
		if code == codes.NotFound || code == codes.InvalidArgument {
			log.Info("NodeID not found. Error:", err.Error())
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.Errorf("failed to get node: %v", err)
		return nil, err
	}
	ip := net.ParseIP(ndIpResp)
	if ip == nil {
		log.Errorf("failed to parse ip from DB: %v", err)
		return nil, err
	}
	return ip, nil
}
