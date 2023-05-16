package grpcserver

import (
	"context"
	"errors"
	"github.com/artems723/monik/internal/server/config"
	"github.com/artems723/monik/internal/server/domain"
	pb "github.com/artems723/monik/internal/server/proto"
	"github.com/artems723/monik/internal/server/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

func New(serv *service.Service, cfg config.Config) *MetricsServer {
	return &MetricsServer{
		service: serv,
		Cfg:     cfg,
	}
}

// start GRPC server
func (s *MetricsServer) Start() {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	serv := grpc.NewServer()
	pb.RegisterMetricsServer(serv, s)
	if err := serv.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	service *service.Service
	Cfg     config.Config
}

func (s *MetricsServer) Save(ctx context.Context, in *pb.SaveMetricRequest) (*pb.SaveMetricResponse, error) {
	var resp pb.SaveMetricResponse
	metric, err := ConvertGRPCtoMetric(in.Metric)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, err = s.service.WriteMetric(ctx, &metric)
	if err != nil {
		return nil, status.Error(errMapping(err), err.Error())
	}
	return &resp, nil

}
func (s *MetricsServer) SaveList(ctx context.Context, in *pb.SaveListMetricsRequest) (*pb.SaveListMetricsResponse, error) {
	var metrics []*domain.Metric
	var resp pb.SaveListMetricsResponse
	for _, v := range in.Metric {
		m, err := ConvertGRPCtoMetric(v)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		metrics = append(metrics, &m)
	}
	err := s.service.WriteMetrics(ctx, &domain.Metrics{Metrics: metrics})
	if err != nil {
		return nil, status.Error(errMapping(err), err.Error())
	}
	return &resp, nil
}

func (s *MetricsServer) Get(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var resp pb.GetMetricResponse
	metricValue, err := s.service.GetMetric(ctx, domain.NewMetric(in.MetricName, ""))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp.Metric = &pb.Metric{
		Id:    metricValue.ID,
		Mtype: string(metricValue.MType),
		Value: *metricValue.Value,
		Delta: *metricValue.Delta,
	}
	return &resp, nil
}

func (s *MetricsServer) GetList(ctx context.Context, in *pb.GetListMetricRequest) (*pb.GetListMetricResponse, error) {
	var resp pb.GetListMetricResponse
	var result []*pb.Metric
	metricList, err := s.service.GetAllMetrics(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	for _, v := range metricList.Metrics {
		result = append(result, ConvertMetrictoGRPC(*v))
	}
	resp.Metric = result
	return &resp, nil
}

func (s *MetricsServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	var resp pb.PingResponse
	err := s.service.Ping()
	if err != nil {
		resp.Ping = false
		return &resp, nil
	}
	resp.Ping = true
	return &resp, nil
}

var (
	ErrStatusNotFound            = errors.New("status not found (404)")
	ErrStatusBadRequest          = errors.New("wrong request (400)")
	ErrStatusNotImplemented      = errors.New("wrong type (501)")
	ErrStatusInternalServerError = errors.New("internal server error(500)")
	ErrWrongType                 = errors.New("wrong type")
)

func errMapping(err error) codes.Code {
	switch {
	case errors.Is(err, ErrStatusBadRequest):
		return codes.InvalidArgument
	case errors.Is(err, ErrStatusNotFound):
		return codes.NotFound
	case errors.Is(err, ErrStatusNotImplemented):
		return codes.Unimplemented
	default:
		return codes.Internal
	}
}

func ConvertGRPCtoMetric(in *pb.Metric) (domain.Metric, error) {
	metric := domain.Metric{
		ID:    in.Id,
		MType: domain.MetricType(in.Mtype),
		Hash:  in.Hash,
	}
	switch in.Mtype {
	case "gauge":
		value := in.Value
		metric.Value = &value
	case "counter":
		value := in.Delta
		metric.Delta = &value
	default:
		return domain.Metric{}, ErrWrongType
	}
	return metric, nil
}

func ConvertMetrictoGRPC(in domain.Metric) *pb.Metric {
	var value float64
	if in.Value != nil {
		value = *in.Value
	}
	var delta int64
	if in.Delta != nil {
		delta = *in.Delta
	}
	result := pb.Metric{
		Id:    in.ID,
		Mtype: string(in.MType),
		Value: value,
		Delta: delta,
		Hash:  in.Hash,
	}
	return &result
}
