package service

// GrpcApplication represents a gRPC Application
type GrpcApplication interface {
	Application
	RegisterServer()
	ListenAndServe()
}

// GrpcApp is the concrete type of GrpcApplication
type GrpcApp struct {
}

// ListenAndServe implements the GrpcApplication interface
func (app *GrpcApp) ListenAndServe() {

}

// RegisterServer implements the GrpcApplication interface
func (app *GrpcApp) RegisterServer() {

}
