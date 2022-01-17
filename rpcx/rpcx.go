package rpcx

import (
	"context"
	"fmt"
	"log"

	// "runtime/metrics"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

const _base_path string = "ee/rpc"

type HostPort struct {
	Host  string
	Port  int16
	Proto string // tcp/udp/http
}

func (hp HostPort) ProtoHostPort() string {
	return fmt.Sprintf("%s@%s:%d",
		hp.Proto, hp.Host, hp.Port)
}
func (hp HostPort) HostPort() string {
	return fmt.Sprintf("%s:%d",
		hp.Host, hp.Port)
}

type RPCXClient struct {
	cli client.XClient
}

func (cli *RPCXClient) Call(ctx context.Context, method string, req, resp interface{}) error {
	if cli == nil || cli.cli == nil {
		return nil
	}
	defer func() {
		cli.Close()
	}()

	return cli.cli.Call(ctx, method, req, resp)
}
func (cli *RPCXClient) Close() error {
	if cli == nil || cli.cli == nil {
		return nil
	}
	defer func() {
		cli.cli = nil
	}()
	return cli.cli.Close()
}

func NewClientByP2P(rpcServer HostPort, serviceName string) (*RPCXClient, error) {
	discovery, err := client.NewPeer2PeerDiscovery(rpcServer.ProtoHostPort(), "")
	if err != nil {
		return nil, err
	}
	xclient := client.NewXClient(serviceName, client.Failtry, client.RandomSelect, discovery, client.DefaultOption)

	return &RPCXClient{cli: xclient}, nil
}

func NewClientByConsul(consulServer HostPort, serviceName string) (*RPCXClient, error) {
	discovery, err := client.NewConsulDiscovery(_base_path, serviceName,
		[]string{fmt.Sprintf("%s:%d", consulServer.Host, consulServer.Port)}, nil)
	if err != nil {
		return nil, err
	}

	option := client.DefaultOption
	option.SerializeType = protocol.JSON
	option.Heartbeat = true
	option.HeartbeatInterval = time.Second * 1
	option.TCPKeepAlivePeriod = 10 * time.Minute
	option.IdleTimeout = 10 * time.Second

	xClient := client.NewXClient(serviceName, client.Failover, client.RoundRobin, discovery, option)

	// xClient.GetPlugins().Add(p) trace plugin
	//xClient.GetPlugins().Add(new(plugin.ClientSpanPlugin))
	//xClient.GetPlugins().Add(new(plugin.ValidatorArgsPlugin))
	return &RPCXClient{cli: xClient}, nil
}

type ServiceHandle interface {
	HandleName() string
}
type RPCXServer struct {
	addr HostPort
	serv *server.Server
}

func (s *RPCXServer) Shutdown(ctx context.Context) error {
	if s == nil || s.serv == nil {
		return nil
	}
	if err := s.serv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown rpcx server error %+v\n", err)
		return err
	}
	return nil
}
func (s *RPCXServer) Relaunch(ctx context.Context) error {
	if s == nil || s.serv == nil {
		return nil
	}
	if err := s.serv.Restart(ctx); err != nil {
		log.Fatalf("restart rpcx server error %+v\n", err)
		return err
	}
	return nil
}
func (s *RPCXServer) Launch(ctx context.Context) error {
	if s == nil || s.serv == nil {
		return nil
	}
	if err := s.serv.Serve(s.addr.Proto,
		fmt.Sprintf("%s:%d", s.addr.Host, s.addr.Port)); err != nil {
		log.Fatalf("start rpcx server error %+v\n", err)
		return err
	}
	return nil
}

func NewServer(serviceHost HostPort, handles []ServiceHandle,
	onRestart, onShutdown func(s *server.Server)) (*RPCXServer, error) {
	serv := server.NewServer()

	serv.RegisterOnRestart(onRestart)   // on restart
	serv.RegisterOnShutdown(onShutdown) // on shutdown

	for _, hanle := range handles {
		if err := serv.RegisterName(hanle.HandleName(), hanle, ""); err != nil {
			log.Fatalf("server regist service %s error %+v\n", hanle.HandleName(), err)
			return nil, err
		}
	}

	return &RPCXServer{addr: serviceHost, serv: serv}, nil
}
func NewServerAndRegisterConsul(serviceHost HostPort, handles []ServiceHandle,
	onRestart, onShutdown func(s *server.Server), consulHost HostPort) (*RPCXServer, error) {

	serv := server.NewServer()

	rp := serverplugin.ConsulRegisterPlugin{
		ServiceAddress: serviceHost.ProtoHostPort(),
		ConsulServers:  []string{consulHost.HostPort()},
		BasePath:       _base_path,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	if err := rp.Start(); err != nil {
		log.Fatalf("regist service error %+v\n", err)
		return nil, err
	}
	serv.Plugins.Add(rp) //consul
	// serv.Plugins.Add(nil) //trace

	serv.RegisterOnRestart(onRestart)
	serv.RegisterOnShutdown(onShutdown)

	for _, hanle := range handles {
		if err := serv.RegisterName(hanle.HandleName(), hanle, ""); err != nil {
			log.Fatalf("server regist service %s error %+v\n", hanle.HandleName(), err)
			return nil, err
		}
		if err := rp.Register(hanle.HandleName(), hanle, ""); err != nil {
			log.Fatalf("consul regist service %s error %+v\n", hanle.HandleName(), err)
			return nil, err
		}
	}

	return &RPCXServer{addr: serviceHost, serv: serv}, nil
}
