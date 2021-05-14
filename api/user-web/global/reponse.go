package global

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

const RequestSuccess = 0

const (
	ClientApiParamsError = 1000 * iota
	ClientGRPCParamsError
	ClientErrorRequest
)

const (
	ServerGRPCInternalError = 2000 * iota
	ServerGRPCUnknownError
	ServerAPIInternalError
)

type Response struct {
	ErrorCode int         `json:"error_code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
}

func JsonError(errCode int, msg string) Response {
	return Response{
		ErrorCode: errCode,
		Msg:       msg,
		Data:      []int{},
	}
}

func JsonSuccess(msg string, data interface{}) Response {
	return Response{
		ErrorCode: RequestSuccess,
		Msg:       msg,
		Data:      data,
	}
}

func HandleGrpcErrorToHttpError(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, JsonError(ClientErrorRequest, e.Message()))
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, JsonError(ServerGRPCInternalError, "internal error"))
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, JsonError(ClientGRPCParamsError, "invalid argument"))
			case codes.Unavailable:
				ctx.JSON(http.StatusInternalServerError, JsonError(ClientGRPCParamsError, "user service unavailable"))
			default:
				ctx.JSON(http.StatusInternalServerError, JsonError(ServerGRPCUnknownError, "unknown error"))
			}
		}
	}
}
