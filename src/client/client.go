package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"math"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	pb "ex03/transmitter"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type statistics struct {
	counter int
	mean    float64
	stdDev  float64
}

type Report struct {
	ID        int64     `bun:",pk,autoincrement"`
	SessionID string    `bun:"session_id"`
	Frequency float64   `bun:"frequency"`
	Timestamp time.Time `bun:"timestamp"`
}

var (
	k float64
)

func main() {
	db := connectDB()
	defer db.Close()
	flag.Float64Var(&k, "k", 2.0, "Anomaly detection coefficient")
	flag.Parse()
	conn, err := grpc.NewClient("localhost:3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection problem:(: %v", err)
	}
	defer conn.Close()

	client := pb.NewTransmitterClient(conn)

	req := &emptypb.Empty{}

	stream, err := client.DataStream(context.Background(), req)
	if err != nil {
		log.Fatalf("error with DataStream calling: %v", err)
	}
	stats := statistics{}
	for {
		response, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error with message receiving: %v", err)
		}

		updateStats(&stats, response.Frequency)
		if stats.counter > 0 && stats.counter%100 == 0 {
			log.Printf("\nNumber of received parametres=%d\nFrequency=%.2f\nTime=%d", stats.counter, response.Frequency, response.Time)
		}
		if stats.counter >= 100 {
			if isAnomaly(&stats, response.Frequency) {
				anomaly := Report{
					SessionID: response.SessionId,
					Frequency: response.Frequency,
					Timestamp: time.Now().UTC(),
				}
				instertAno(db, &anomaly)
				log.Printf("\nANOMALY DETECTED!!!\nReceived frequency=%.2f\nMean=%.2f\nStandard deviation=%.2f\n", response.Frequency, stats.mean, stats.stdDev)
			}
		}
	}
}

func updateStats(stats *statistics, frequency float64) {
	stats.counter++
	// mean
	oldMean := stats.mean
	stats.mean += (frequency - oldMean) / float64(stats.counter)
	// standart deviation
	if stats.counter > 1 {
		variance := ((float64(stats.counter-1) * math.Pow(stats.stdDev, 2)) + (frequency-oldMean)*(frequency-stats.mean)) / float64(stats.counter)
		stats.stdDev = math.Sqrt(variance)
	}
}

func isAnomaly(stats *statistics, frequency float64) bool {
	if stats.stdDev == 0 {
		return false
	}

	lowerBound := stats.mean - k*stats.stdDev
	upperBound := stats.mean + k*stats.stdDev

	return stats.counter > 0 && (frequency < lowerBound || frequency > upperBound)
}

func connectDB() *bun.DB {
	sqldb, err := sql.Open("postgres", "user=postgres password=1 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return bun.NewDB(sqldb, pgdialect.New())
}

func instertAno(db *bun.DB, anomaly *Report) {
	_, err := db.NewInsert().Model(anomaly).Exec(context.Background())
	if err != nil {
		log.Printf("Error inserting anomaly: %v", err)
	}
}
