package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(id string) *Error {
	return &Error{&errorpb.Error{Id: id}}
}

type Error struct {
	err *errorpb.Error
}

func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.GetCode(), e.GetErrorData())
}

func (e *Error) String() string {
	return fmt.Sprintf("error: id=%s code=%d message=%s metadata=%v",
		e.Proto().Id, e.Proto().Code, e.Proto().Message, e.Proto().Metadata)
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) GetErrorData() string {
	var dt, err = json.Marshal(e.err)
	if err != nil {
		return err.Error()
	}
	return string(dt)
}

func (e *Error) Ok() bool              { return codes.Code(e.err.Code) == codes.OK }
func (e *Error) IsBiz() bool           { return e.err.Id != "" }
func (e *Error) Proto() *errorpb.Error { return e.err }
func (e *Error) WithCode(code codes.Code) *Error {
	var err = proto.Clone(e.err).(*errorpb.Error)
	err.Code = int32(code)
	return &Error{err: err}
}

func (e *Error) Clone() *Error {
	return &Error{err: proto.Clone(e.err).(*errorpb.Error)}
}

func (e *Error) GetCode() codes.Code { return codes.Code(e.err.Code) }

func (e *Error) GetMsg() string { return e.err.Message }
func (e *Error) GetId() string  { return e.err.Id }
func (e *Error) GetMetadata() map[string]string {
	var md = make(map[string]string, len(e.err.Metadata))
	for i := range e.err.Metadata {
		md[e.err.Metadata[i].K] = e.err.Metadata[i].V
	}
	return md
}

func (e *Error) Err(err error) error {
	if err == nil {
		return nil
	}

	if _err, ok := err.(*Error); ok {
		e.Proto().Code = int32(_err.GetCode())
		e.Proto().Message = _err.GetMsg()
		e.Proto().Metadata = append(e.Proto().Metadata, _err.Proto().Metadata...)
		return e
	}

	e.err.Message = err.Error()
	e.Metadata("err_msg", fmt.Sprintf("%v", err))
	e.Metadata("err", e.err.Message)
	return e
}

func (e *Error) Msg(format string, a ...interface{}) error {
	e.err.Message = fmt.Sprintf(format, a...)
	return e
}

func (e *Error) Metadata(k, v string) {
	e.err.Metadata = append(e.err.Metadata, &errorpb.Kv{K: k, V: v})
}

// HTTPStatus returns the Status represented by se.
func (e *Error) HTTPStatus() int {
	switch e.err.Code {
	case 0:
		return http.StatusOK
	case 1:
		return http.StatusInternalServerError
	case 2:
		return http.StatusInternalServerError
	case 3:
		return http.StatusBadRequest
	case 4:
		return http.StatusRequestTimeout
	case 5:
		return http.StatusNotFound
	case 6:
		return http.StatusConflict
	case 7:
		return http.StatusForbidden
	case 8:
		return http.StatusTooManyRequests
	case 9:
		return http.StatusPreconditionFailed
	case 10:
		return http.StatusConflict
	case 11:
		return http.StatusBadRequest
	case 12:
		return http.StatusNotImplemented
	case 13:
		return http.StatusInternalServerError
	case 14:
		return http.StatusServiceUnavailable
	case 15:
		return http.StatusInternalServerError
	case 16:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func FromErr(err error) *Error {
	switch err.(type) {
	case nil:
		return nil
	case *Error:
		return err.(*Error)
	}

	se, ok := err.(interface{ GRPCStatus() *status.Status })
	if ok && se != nil {
		if se.GRPCStatus().Code() == codes.OK {
			return nil
		}

		var e = new(errorpb.Error)
		if json.Unmarshal([]byte(se.GRPCStatus().Message()), e) == nil {
			return &Error{err: e}
		}

		return &Error{err: &errorpb.Error{
			Code:    int32(se.GRPCStatus().Code()),
			Id:      "",
			Message: se.GRPCStatus().Message(),
		}}
	}

	return New("").WithCode(codes.Unknown).Err(err).(*Error)
}
