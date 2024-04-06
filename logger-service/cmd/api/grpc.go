package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, r *logs.LogRequest) (*logs.LogResponse, error) {
	input := r.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	if err := l.Models.LogEntry.Insert(logEntry); err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged via gRPC!",
	}

	return res, nil
}

func (app *Config) gRPCListen(){
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Printf("gRPC Server started on port %s", gRPCPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}