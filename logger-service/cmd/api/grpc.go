package main

import (
	"context"
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

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	res := &logs.LogResponse{Result: "failed"}
	err := l.Models.Entry.Insert(logEntry)
	if err != nil {
		res = &logs.LogResponse{Result: "Logged!"}
	}

	return res, err
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", gRpcPort)
	if err != nil {
		log.Panicln("Failed to listen for gRPC", err)
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	if err = s.Serve(lis); err != nil {
		log.Panicln("Failed to start gRPC server", err)
	}
}
