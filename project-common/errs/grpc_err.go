package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	common "gwh.com/project-common"
)

func GrpcError(err *BError) error {
	return status.Errorf(codes.Code(err.Code), err.Msg)
}

func PraseGrpcError(err error) (common.BusinessCode, string) {
	fromError, _ := status.FromError(err)
	code := fromError.Code()
	msg := fromError.Message()
	return common.BusinessCode(code), msg
}
