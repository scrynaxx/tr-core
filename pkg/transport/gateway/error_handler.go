package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errorHandler(_ context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	if isInternalError(err) {
		err = status.Error(codes.Internal, "internal server error")
	}

	st, httpStatus := statusFromError(err)
	body := errorResponse{
		Code:    statusCode(st),
		Message: st.Message(),
	}

	w.Header().Set("Content-Type", marshaler.ContentType(body))
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(body)
}

func isInternalError(err error) bool {
	var httpErr *runtime.HTTPStatusError
	if errors.As(err, &httpErr) {
		return httpErr.HTTPStatus == http.StatusInternalServerError
	}

	st, ok := status.FromError(err)
	if !ok {
		return true
	}

	return runtime.HTTPStatusFromCode(st.Code()) == http.StatusInternalServerError
}

func statusFromError(err error) (*status.Status, int) {
	var httpErr *runtime.HTTPStatusError
	if errors.As(err, &httpErr) {
		st, ok := status.FromError(httpErr.Err)
		if !ok {
			st = status.New(codes.Unknown, httpErr.Err.Error())
		}

		return st, httpErr.HTTPStatus
	}

	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Unknown, err.Error())
	}

	return st, runtime.HTTPStatusFromCode(st.Code())
}

func statusCode(st *status.Status) string {
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok && info.GetReason() != "" {
			return info.GetReason()
		}
	}

	return strings.ToLower(st.Code().String())
}
