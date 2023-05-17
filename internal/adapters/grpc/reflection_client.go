package grpc

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var ErrorMethodNotFind = errors.New("could not find method from remote")

func GetReflectionInfo(remoteAddr, fullName string, opts ...grpc.DialOption) (*desc.MethodDescriptor,
	*desc.MessageDescriptor,
	*desc.MessageDescriptor,
	error) {
	if fullName[0] == '/' {
		fullName = fullName[1:]
	}
	names := strings.Split(fullName, "/")
	serviceName := names[0]
	methodName := names[1]
	// create a new connection to remote server
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(remoteAddr, opts...)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error to get reflection information from server")
	}
	// create reflection client from connection
	client := grpcreflect.NewClientAuto(context.Background(), conn)
	defer conn.Close()
	fileDescriptor, err := client.FileContainingSymbol(serviceName)
	if err != nil {
		return nil, nil, nil, err
	}
	service := fileDescriptor.FindService(serviceName)
	method := service.FindMethodByName(methodName)
	if method == nil {
		return nil, nil, nil, errors.Wrap(ErrorMethodNotFind, "cloud not find method "+fullName)
	}
	return method, method.GetInputType(), method.GetOutputType(), nil
}

func ReflectionInvoke(req *http.Request, requestData []byte, method *desc.MethodDescriptor, input *desc.MessageDescriptor) (*Message, error) {
	ctx := req.Context()
	ctx = appendMetadata(ctx, req)
	conn, err := grpc.Dial(req.URL.Host, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to remote server")
	}
	defer conn.Close()
	stub := grpcdynamic.NewStub(conn)
	inputMessage := dynamic.NewMessage(input)
	err = inputMessage.UnmarshalJSONPB(&jsonpb.Unmarshaler{AllowUnknownFields: true}, requestData)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal input data into request grpc message")
	}
	response, err := stub.InvokeRpc(ctx, method, inputMessage)
	if err != nil {
		return nil, err
	}
	r, _ := dynamic.AsDynamicMessage(response)
	return &Message{r}, err
}
