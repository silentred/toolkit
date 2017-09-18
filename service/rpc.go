package service

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"google.golang.org/grpc"
)

// GrpcApplication represents a gRPC Application
type GrpcApplication interface {
	Application
	SetServer(*grpc.Server)
	GetServer() *grpc.Server
	ListenAndServe()
}

// GrpcApp is the concrete type of GrpcApplication
type GrpcApp struct {
	App
	server *grpc.Server
}

// NewGrpcApp returns a new GrpcApp
func NewGrpcApp(s *grpc.Server) *GrpcApp {
	app := &GrpcApp{
		App:    NewApp(),
		server: s,
	}
	app.Set("app.rpc", app, new(GrpcApplication))
	return app
}

// Initialize web application
func (app *GrpcApp) Initialize() {
	initialize(app)
}

// ListenAndServe implements the GrpcApplication interface
func (app *GrpcApp) ListenAndServe() {
	var port = fmt.Sprintf("%s:%d", app.GetConfig().Host, app.GetConfig().Port)
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	info := app.server.GetServiceInfo()
	if len(info) == 0 {
		log.Fatalf("grpc server has to register service first. %v", info)
	}

	go func() {
		err = app.server.Serve(l)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	app.handleSignal()
}

// SetServer implements the GrpcApplication interface
func (app *GrpcApp) SetServer(s *grpc.Server) {
	app.server = s
}

// GetServer implements the GrpcApplication interface
func (app *GrpcApp) GetServer() *grpc.Server {
	return app.server
}

func (app *GrpcApp) handleSignal() {
	var (
		c chan os.Signal
		s os.Signal
	)
	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
		syscall.SIGINT, syscall.SIGSTOP)
	// Block until a signal is received.
	for {
		s = <-c
		log.Printf("get a signal %s \n", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			app.server.GracefulStop()
			runHooks(ShutdownHook, app)
			return
		case syscall.SIGHUP:
			// TODO reload
			//return
		default:
			return
		}
	}
}

// ListenAndServe the gRPC server at addr
func ListenAndServe(s *grpc.Server, addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	info := s.GetServiceInfo()
	if len(info) == 0 {
		log.Fatalf("grpc server has to register service first. %v", info)
	}

	err = s.Serve(l)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
