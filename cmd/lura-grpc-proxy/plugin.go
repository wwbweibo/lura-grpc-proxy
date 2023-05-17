package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wwbweibo/lura-grpc-proxy/internal"
	"github.com/wwbweibo/lura-grpc-proxy/internal/adapters/grpc"
	"github.com/wwbweibo/lura-grpc-proxy/internal/domain"
	"github.com/wwbweibo/lura-grpc-proxy/internal/helper"

	"google.golang.org/grpc/status"
)

type PluginRegister struct {
	Name string
}

func NewPluginRegister() PluginRegister {
	plugin := PluginRegister{
		Name: internal.PluginName,
	}
	logger.Debug("init plugin: " + internal.PluginName)
	return plugin
}

func (r PluginRegister) RegisterLogger(v interface{}) {
	l, ok := v.(internal.Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", ClientRegisterer))
}

func (r PluginRegister) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(internal.PluginName, r.registerClients)
}

func (r PluginRegister) registerClients(_ context.Context, extra map[string]interface{}) (http.Handler, error) {
	logger.Info("register client plugin")
	config, _ := extra[internal.PluginName].(map[string]interface{})
	gRPCMethodName := config["method_name"].(string)
	// read all config from local, and parse into memory
	// here will meet two scene
	// 1. user use an all-in-one configuration.
	// 2. user use template file.
	endpoint, host, err := internal.LoadRouteFromConfiguration(gRPCMethodName)
	if err != nil {
		panic(err)
	}

	methodDesc, input, _, err := grpc.GetReflectionInfo(host, gRPCMethodName)
	if err != nil {
		panic(err)
	}

	// return the actual handler wrapping or your custom logic, so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		bts, err := grpc.DecodeRequestToMapBytes(req, endpoint)
		if err != nil {
			logger.Error("error to resolve request data", err)
			helper.WriteResponse(http.StatusServiceUnavailable, domain.HttpResponse{Code: -1, Message: "调用后端服务失败，请检查配置是否正确或者联系管理员"}, w)
			return
		}
		response, err := grpc.ReflectionInvoke(req, bts, methodDesc, input)
		if err != nil {
			if stu, ok := status.FromError(err); ok {
				logger.Error("error to resolve request data", err)
				helper.WriteResponse(http.StatusInternalServerError, domain.HttpResponse{Code: int(stu.Code()), Message: stu.Message()}, w)
				return
			} else {
				logger.Error("error to resolve request data", err)
				helper.WriteResponse(http.StatusBadGateway, domain.HttpResponse{Code: -1, Message: "调用后端服务失败，请检查配置是否正确或者联系管理员"}, w)
				return
			}
		} else {
			var rpcResp domain.HttpResponse
			rpcResp.Message = "success"
			rpcResp.Code = 0
			rpcResp.Data = response
			helper.WriteResponse(http.StatusOK, rpcResp, w)
		}
		return
	}), nil
}
