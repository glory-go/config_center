package nacos

import (
	"github.com/glory-go/glory/log"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfigCenter struct {
	config *NacosConfig
	client config_client.IConfigClient
}

func (c *NacosConfigCenter) LoadConfig(key, group string) (string, error) {
	val, err := c.client.GetConfig(vo.ConfigParam{
		DataId: key,
		Group:  group,
	})
	if err != nil {
		return "", err
	}
	return val, err
}

func (c *NacosConfigCenter) SyncConfig(key, group string, value *string, cancel <-chan struct{}) error {
	val, err := c.LoadConfig(key, group)
	if err != nil {
		return err
	}
	*value = val
	if err := c.client.ListenConfig(vo.ConfigParam{
		DataId: key,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			*value = data
		},
	}); err != nil {
		return err
	}
	// 收到cancel指令后cancel
	go func() {
		defer func() {
			if e := recover(); err != nil {
				log.Error(e)
			}
		}()
		<-cancel
		c.client.CancelListenConfig(vo.ConfigParam{
			DataId: key,
			Group:  group,
		})
	}()

	return nil
}
