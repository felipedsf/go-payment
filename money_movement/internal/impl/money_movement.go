package mm

import (
	"context"
	"database/sql"
	pb "github.com/felipedsf/go-payment/money_movement/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcMoneyMovement struct {
	db *sql.DB
	pb.UnimplementedMoneyMovementServiceServer
}

func NewGrpcMoneyMovement(db *sql.DB) *GrpcMoneyMovement {
	return &GrpcMoneyMovement{
		db: db,
	}
}

func (this *GrpcMoneyMovement) Authorize(ctx context.Context, payload *pb.AuthorizePayload) (*pb.AuthorizeResponse, error) {
	return nil, nil
}

func (this *GrpcMoneyMovement) Capture(ctx context.Context, in *pb.CapturePayload) (*emptypb.Empty, error) {
	return nil, nil
}
