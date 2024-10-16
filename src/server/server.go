package main

import (
	pb "ex03/transmitter"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	lis, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatalf("Listening error: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTransmitterServer(grpcServer, &server{})

	log.Println("Starting server on port :3333")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server launch error: %v", err)
	}
}

type server struct {
	pb.UnimplementedTransmitterServer
}

func (s *server) DataStream(empty *emptypb.Empty, reqStream pb.Transmitter_DataStreamServer) error {
	uu := getUuid()
	mean := getMean()
	stdDev := getStandardDeviation()
	logging(uu, mean, stdDev)
	for {
		response := getResponse(uu, mean, stdDev)

		if err := reqStream.Send(response); err != nil {
			log.Printf("Error with response sending %v", err)
			return err
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func getTime() int64 {
	return time.Now().UTC().Unix()
}
func getUuid() string {
	return uuid.New().String()
}

func getfrequency(mean, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + mean
}

func getMean() float64 {
	return rand.Float64()*20 - 10
}

func getStandardDeviation() float64 {
	return rand.Float64()*1.2 + 0.3
}

func getResponse(uu string, mean float64, stdDev float64) *pb.Response {

	resp := &pb.Response{
		SessionId: uu,
		Frequency: getfrequency(mean, stdDev),
		Time:      getTime(),
	}
	return resp
}

func logging(uu string, mean float64, stdDev float64) {
	log.Printf("\nuuid: %s\nmean: %.2f\nstandart deviation: %.2f\n", uu, mean, stdDev)
}
