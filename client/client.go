package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "github.com/dev-code0101/golang_cron/pb" // Update with the actual path
)

func main() {
    // Set up a connection to the server.
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    // Create a new client for the ReportService.
    client := pb.NewReportServiceClient(conn)

    // Create a context with a timeout for the health check request.
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // Call the HealthCheck method.
    resp, err := client.HealthCheck(ctx, &pb.HealthCheckRequest{})
    if err != nil {
        log.Fatalf("could not check health: %v", err)
    }

    // Log the response.
    log.Printf("Health Check Status: %s", resp.Status)
}
