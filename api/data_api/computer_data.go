package data_api

import (
	"blogX_server/common/resp"
	"blogX_server/utils/computer"
	"github.com/gin-gonic/gin"
)

type ComputerDataResponse struct {
	CpuPercent  float64 `json:"cpu_percent"`
	MemPercent  float64 `json:"mem_percent"`
	DiskPercent float64 `json:"disk_percent"`
}

// ComputerDataView 电脑数据
func (DataApi) ComputerDataView(c *gin.Context) {
	var data = ComputerDataResponse{
		CpuPercent:  computer.GetCpuPercent(),
		MemPercent:  computer.GetMemPercent(),
		DiskPercent: computer.GetDiskPercent(),
	}
	resp.OkWithData(data, c)
}
