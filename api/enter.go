package api

import (
	"blogX_server/api/image_api"
	"blogX_server/api/log_api"
	"blogX_server/api/site_api"
)

type Api struct {
	SiteApi  site_api.SiteApi
	LogApi   log_api.LogApi
	ImageApi image_api.ImageApi
}

// App 实例化 以供外部调用Api
var App = Api{}
