package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimoktorr/monitoring/internal/app"
	"github.com/dimoktorr/monitoring/internal/pkg/api"
	"github.com/dimoktorr/monitoring/internal/pkg/metrics"
	"github.com/dimoktorr/monitoring/internal/pkg/persistent/repository"
	"github.com/dimoktorr/monitoring/internal/pkg/persistent/storage"
	v1 "github.com/dimoktorr/monitoring/pkg/api/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type App struct {
	metricsServer *http.Server
	cfg           *Config
	grpcServer    *grpc.Server
	service       *app.Service
}

func New(ctx context.Context, cfg *Config) (*App, error) {
	a := &App{
		cfg: cfg,
	}

	repoConn, err := repository.NewPostgresConn(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(repoConn.Pgx, repoConn.ScanAPI)

	redis, err := storage.NewRedisUniversalClient(ctx, a.cfg.Redis)
	if err != nil {
		return nil, err
	}

	storageProduct := storage.NewStorage(redis, a.cfg.Redis.TTL, "example-storage")

	if err := a.newTracer(); err != nil {
		return nil, err
	}

	a.newService(
		metrics.New(),
		repo,
		storageProduct,
	)

	a.newMetricsServer()
	a.newGRPCServer()

	return a, nil
}

func (a *App) newTracer() error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(a.cfg.Tracing.URL)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(a.cfg.Tracing.ServiceName),
			attribute.String("environment", a.cfg.Tracing.Environment),
			attribute.Int64("solObjectID", a.cfg.Tracing.SolObjectID),
			attribute.Int64("imsSystemID", a.cfg.Tracing.ImsSystemID),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	return nil
}

func (a *App) newService(
	metrics *metrics.Metrics,
	repo *repository.Repository,
	storageClient *storage.Storage,
) {
	a.service = app.New(
		metrics,
		repo,
		storageClient,
	)
}

func (a *App) newMetricsServer() {
	a.metricsServer = &http.Server{
		Handler: promhttp.Handler(),
		Addr:    fmt.Sprintf("%s:%s", a.cfg.Prometheus.Host, a.cfg.Prometheus.Port),
	}
}

func (a *App) Start() {
	go func() {
		log.Println("metrics http server started, port", a.cfg.Prometheus.Port)
		if err := a.metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("start metrics http server failed")
		}
	}()

	l, err := newTCPListener(a.cfg.Service.Host, a.cfg.Service.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("grpc server started", a.cfg.Service.GRPCPort)
		err := a.grpcServer.Serve(l)
		if err != nil {
			log.Fatal("grpc server failed", err)
		}
	}()
}

func (a *App) Stop(ctx context.Context) error {
	<-ctx.Done()

	graceCtx, graceCancel := context.WithTimeout(ctx, a.cfg.Service.ShutdownContextTimeout)
	defer graceCancel()

	if err := a.metricsServer.Shutdown(graceCtx); err != nil {
		return fmt.Errorf("could not gracefully shutdown metrics http server, err: %w", err)
	}

	a.grpcServer.GracefulStop()

	log.Println("http servers stopped...")

	return nil
}

func (a *App) newGRPCServer() {
	unaryInterceptor := grpc.ChainUnaryInterceptor(
		UnaryRequestIDServerInterceptor(),
	)

	var options []grpc.ServerOption
	options = append(options,
		unaryInterceptor,
	)

	a.grpcServer = grpc.NewServer(options...)

	serverApi := api.NewServer(a.service)
	a.grpcServer.RegisterService(&v1.ExampleService_ServiceDesc, serverApi)
}

func newTCPListener(host, port string) (net.Listener, error) {
	l, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %q: %w", port, err)
	}

	return l, nil
}
