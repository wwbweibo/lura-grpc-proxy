package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func CallRemoteService(req *http.Request, method string, input, output interface{}) error {
	ctx := req.Context()
	ctx = appendMetadata(ctx, req)
	conn, err := grpc.Dial(req.URL.Host, grpc.WithInsecure())
	if err != nil {
		return err
	}
	return conn.Invoke(ctx, method, input, output)
}

func appendMetadata(ctx context.Context, req *http.Request) context.Context {
	md := metadata.MD{}
	for k, v := range req.Header {
		md.Append(strings.ToLower(k), v...)
	}
	span := trace.FromContext(ctx)
	sampled := "0"
	if span.SpanContext().TraceOptions.IsSampled() {
		sampled = "1"
	}
	md.Append("uber-trace-id", fmt.Sprintf("%s:%s:%s", span.SpanContext().TraceID.String(), span.SpanContext().SpanID, sampled))
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}

func DecodeRequest(grpcRequest interface{}, req *http.Request, endpoint string) error {
	params := make(map[string]interface{})
	decodePath(req, &params, endpoint)
	decodeQueryString(req, &params)
	err := decodeJsonBody(req, &params)
	if err != nil {
		return errors.Wrap(err, "error when decode json body")
	}
	bts, err := json.Marshal(params)
	if err != nil {
		return errors.Wrap(err, "error when marshal decoded input")
	}
	err = json.Unmarshal(bts, &grpcRequest)
	if err != nil {
		return errors.Wrap(err, "error when unmashal grpc request body")
	}
	return nil
}

func DecodeRequestToMapBytes(req *http.Request, endpoint string) ([]byte, error) {
	params := make(map[string]interface{})
	decodePath(req, &params, endpoint)
	decodeQueryString(req, &params)
	err := decodeJsonBody(req, &params)
	if err != nil {
		return nil, errors.Wrap(err, "error when decode json body")
	}
	bts, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "error when marshal decoded input")
	}
	return bts, nil
}

func decodePath(req *http.Request, set *map[string]interface{}, endpoint string) {
	// global.Logger.Debug("defined endpoint", "endpoint", endpoint)
	// global.Logger.Debug("requested url", "url", req.URL.Path)
	// todo: 获取路径，并从路径中拿到路径参数
	defined := strings.Split(endpoint, "/")
	request := strings.Split(req.URL.Path, "/")
	if len(defined) != len(request) {
		panic("defined route not match requested route")
	}
	for i := 1; i < len(defined); i++ {
		if defined[i][0] == '{' {
			(*set)[defined[i][1:len(defined[i])-1]] = request[i]
		}
	}
}

func decodeQueryString(req *http.Request, set *map[string]interface{}) {
	query := map[string][]string(req.URL.Query())
	for k, v := range query {
		(*set)[k] = v[0]
	}
}

func decodeJsonBody(req *http.Request, set *map[string]interface{}) error {
	bts, err := io.ReadAll(req.Body)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		} else {
			return errors.Wrap(err, "error when read request body")
		}
	}
	if len(bts) == 0 {
		return nil
	}
	temp := make(map[string]interface{})
	err = json.Unmarshal(bts, &temp)
	if err != nil {
		return errors.Wrap(err, "error when unmarshal request body")
	}
	for k, v := range temp {
		(*set)[k] = v
	}
	return nil
}

func CreateRequestModel(input, output string, types map[string]reflect.Type) (interface{}, interface{}, error) {
	var inputType, outputType reflect.Type
	if t, ok := types[input]; ok {
		inputType = t
	} else {
		return nil, nil, errors.New("type not find")
	}
	if t, ok := types[output]; ok {
		outputType = t
	} else {
		return nil, nil, errors.New("type not find")
	}
	inputValue := reflect.New(inputType).Interface()
	outputValue := reflect.New(outputType).Interface()
	return inputValue, outputValue, nil
}

func DecodeName(method string, descs map[string]protoreflect.FileDescriptor) (string, string, error) {
	if method[0] == '/' {
		method = method[1:]
	}
	names := strings.Split(method, "/")
	idx := strings.LastIndex(names[0], ".")
	pkgName := names[0][:idx]
	serviceName := names[0][idx+1:]
	if desciptor, ok := descs[pkgName]; !ok {
		return "", "", errors.New("could not package")
	} else {
		for i := 0; i < desciptor.Services().Len(); i++ {
			service := desciptor.Services().Get(i)
			if string(service.Name()) == serviceName {
				for j := 0; j < service.Methods().Len(); j++ {
					method := service.Methods().Get(j)
					if string(method.Name()) == names[1] {
						input := string(method.Input().FullName())
						output := string(method.Output().FullName())
						return input, output, nil
					}
				}
			}
		}
	}
	return "", "", errors.New("no matching method find")
}
