package service

import (
	"context"

	pb "grpc-course/pb"
)

type LaptopServer struct {

}

func NewLaptopServer() *LaptopServer {
	return &LaptopServer{}
}

//CreateLaptop is a unary RPC to create a new laptop
func CreateLaptop(
	ctx context.Context, a
	req *pb.CreateLaptopRequest
) (*pb.CreateLaptopResponse, error) {
	laptop:= req.GetLaptop();

}
