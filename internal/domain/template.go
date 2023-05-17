package domain

// Template is the define of FC_SETTINGS dir file, user can add more field in files, but must follow
// this basic define, otherwise, error will occured when plugin load these files.
type Template struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Backends []struct {
		Url        string `json:"url"`
		Host       string `json:"host"`
		GrpcMethod string `json:"grpc_method"`
	} `json:"backends"`
}

/*
	so the settings/xxx.json shold be like this
	{
		"routes": [
			{
				"endpoint": "",
				"method": "",
				"backends": [
					"url": "",
					"host": "",
					"grpc_method": ""
				]
			}
		]
	}
*/
