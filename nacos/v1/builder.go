package nacos

import (
	"strconv"

	configmanager "github.com/glory-go/glory/config_manager"
	"github.com/glory-go/glory/utils"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfig struct {
	NamespaceID string `json:"namespace_id,omitempty"`

	Address string `json:"address,omitempty"`
	Path    string `json:"path,omitempty"`
	Port    string `json:"port,omitempty"`
	Scheme  string `json:"scheme,omitempty"`
}

func NacosConfigCenterBuilder(conf map[string]string) (configmanager.ConfigCenter, error) {
	targetConf := &NacosConfig{}
	if err := utils.ConvertInto(conf, targetConf); err != nil {
		return nil, err
	}
	if targetConf.Port == "" {
		targetConf.Port = "80"
	}
	if targetConf.Scheme == "" {
		targetConf.Scheme = "http"
	}
	port, err := strconv.Atoi(targetConf.Port)
	if err != nil {
		return nil, err
	}
	// 连接到nacos
	clientConfig := constant.ClientConfig{
		NamespaceId:         targetConf.NamespaceID, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      targetConf.Address,
			ContextPath: targetConf.Address,
			Port:        uint64(port),
			Scheme:      targetConf.Scheme,
		},
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, err
	}

	return &NacosConfigCenter{
		config: targetConf,
		client: configClient,
	}, nil
}
