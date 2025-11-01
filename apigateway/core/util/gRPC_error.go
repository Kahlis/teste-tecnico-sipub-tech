package util

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCError struct {
	Code    codes.Code `json:"code"`
	Message string     `json:"message"`
	Details []any      `json:"details,omitempty"`
}

func ParseGRPCError(err error) *GRPCError {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return &GRPCError{
			Code:    st.Code(),
			Message: err.Error(),
		}
	}

	return &GRPCError{
		Code:    st.Code(),
		Message: st.Message(),
		Details: st.Details(),
	}
}

func GRPCToZap(grpcErr *GRPCError) (zap.Field, zap.Field, zap.Field) {
	code := zap.Uint32("code", uint32(grpcErr.Code))
	message := zap.String("message", grpcErr.Message)
	details := zap.Any("details", grpcErr.Details)
	return code, message, details
}
