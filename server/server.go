package main

import (
    "context"
    "log"
    "net"
    "sync"
    "time"

    "github.com/robfig/cron/v3"
    "google.golang.org/grpc"
    pb "github.com/dev-code0101/golang_cron/pb" // Update with the actual path
)

type server struct {
    pb.UnimplementedReportServiceServer
    reports map[string]string
    mu      sync.Mutex
}

func (s *server) GenerateReport(ctx context.Context, req *pb.GenerateReportRequest) (*pb.GenerateReportResponse, error) {
    reportID := time.Now().Format("20060102150405") + "_" + req.UserId
    s.mu.Lock()
    s.reports[reportID] = req.UserId
    s.mu.Unlock()

    log.Printf("Generated report for UserID: %s, ReportID: %s", req.UserId, reportID)
    return &pb.GenerateReportResponse{ReportId: reportID, Error: ""}, nil
}

func (s *server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    return &pb.HealthCheckResponse{Status: "Healthy"}, nil
}

func main() {
    reports := make(map[string]string)
    srv := &server{reports: reports}

    // Create a TCP listener on port 50051
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterReportServiceServer(grpcServer, srv)

    // Set up cron job
    c := cron.New()
    c.AddFunc("@every 10s", func() {
        userIDs := []string{"user1", "user2", "user3"}
        for _, userID := range userIDs {
            _, err := srv.GenerateReport(context.Background(), &pb.GenerateReportRequest{UserId: userID})
            if err != nil {
                log.Printf("Error generating report for UserID %s: %v", userID, err)
            }
        }
    })
    c.Start()

    log.Println("Starting gRPC server on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
