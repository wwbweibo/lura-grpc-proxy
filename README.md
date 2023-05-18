# lura-grpc-proxy

lura-grpc-proxy is one of the client plugin that support you to convert incoming http request to grpc request, then send to backend service. this plugin also convert the grpc response to http response.

## usage

the recommended way to start with lura project is using [krankend](https://www.krakend.io/), which base on lura.

this plugin was tested on krankend-ce, but it should also sopport other gateway based on lura. If find any promblem on other gateway, please make an issue.

> To use this plugin, your server is required to set the grpc-reflection on.

### start with single krakend.json file.

If you are using single krakend.json file, use this plugin as a simple http client plugin. 

To configure this plugin, using extra_config in your endpoint. The endpoint section should like this.

```json
{
    "endpoint": "/test/api/hello/v1",
    "method": "POST",
    "output_encoding": "no-op",
    "input_headers": ["*"],
    "input_query_strings": ["*"],
    "backend": [
    {
        "url_pattern": "/api/hello/v1",
        "encoding": "no-op",
        "sd": "static",
        "method": "GET",
        "host": [ "demo:9999" ],
        "disable_host_sanitize": false,
        "extra_config": {
                "plugin/http-client": {
                    "name": "lura-grpc-proxy",
                    "lura-grpc-proxy": {
                        "method_name": "api.hello.v1.UserService/Login"
                    }
                }
            }
        }
    ]
}
```

There are some limitation you should care about.

1. currently, this plugin will only support one bakend with one host.
2. the `url_pattern` in backend shoud be the same as `endpoint`
3. the `method_name` in plugin config is required. 

### start with flexible configuration

about how to using flexible configuration in lura, please visit [Flexible Configuration: template-based config](https://www.krakend.io/docs/configuration/flexible-config/)

There are some requirement on your setting json file. Your setting json file must meet the following format requirements. If not, the plugin will failed to start up.

```json
{
    "routes": [
        {
            "endpoint": "",
            "method": "",
            "backends": [
                {
                    "url": "",
                    "host": "",
                    "grpc_method": ""
                }
            ]
        }
    ]
}
```

## request and reponse convert

This project are using [jhump/protoreflect](https://github.com/jhump/protoreflect) to process grpc reflection invoke. 

Plugin will read parameter from these location, then merge them into a map, finally using protoreflect to convert then into grpc message send to targe service:

1. path parameter
2. query string
3. all request body

for example, a post request like this:

```
curl http://example.com/api/v1/user/mike?ax=bcs -d '{"name": "mike", "age": 15}' -X POST -H "Content-Type: application/json"
```

which the endpoint is defined as `/api/v1/user/{user_name}`, the final request send to backend is like this: 

```json
{
    "user_name": "mike",
    "ax": "bcs",
    "name": "mike",
    "age": 15
}
```
