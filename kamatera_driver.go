package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rancher/kontainer-engine/drivers/options"
	"github.com/rancher/kontainer-engine/types"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Driver struct {
	driverCapabilities types.Capabilities
}

func (d *Driver) GetDriverCreateOptions(ctx context.Context) (*types.DriverFlags, error) {
	logrus.Info("[Kamatera] GetDriverCreateOptions")
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "the internal name of the cluster in Rancher",
	}
	driverFlag.Options["display-name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "the display name of the cluster in Rancher",
	}
	driverFlag.Options["api-client-id"] = &types.Flag{
		Type: types.StringType,
		//Password: true,
		Usage: "Kamatera API Client ID",
	}
	driverFlag.Options["api-secret"] = &types.Flag{
		Type: types.StringType,
		//Password: true,
		Usage: "Kamatera API Secret",
	}
	driverFlag.Options["datacenter"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The datacenter where the cluster will be created",
	}
	driverFlag.Options["private-network-name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The private network name to assign to all nodes",
	}
	driverFlag.Options["sshkey-private"] = &types.Flag{
		Type: types.StringType,
		//Password: true,
		Usage: "The private SSH key to use for the nodes",
	}
	driverFlag.Options["sshkey-public"] = &types.Flag{
		Type:  types.StringType,
		Usage: "The public SSH key to use for the nodes",
	}
	//driverFlag.Options["node-pools"] = &types.Flag{
	//	Type:  types.StringSliceType,
	//	Usage: "The list of node pools created for the cluster",
	//}
	return &driverFlag, nil
}

func (d *Driver) GetDriverUpdateOptions(ctx context.Context) (*types.DriverFlags, error) {
	logrus.Info("[Kamatera] GetDriverUpdateOptions")
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["api-client-id"] = &types.Flag{
		Type: types.StringType,
		//Password: true,
		Usage: "Kamatera API Client ID",
	}
	driverFlag.Options["api-secret"] = &types.Flag{
		Type: types.StringType,
		//Password: true,
		Usage: "Kamatera API Secret",
	}
	return &driverFlag, nil
}

func (d *Driver) Create(ctx context.Context, opts *types.DriverOptions, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Info("[Kamatera] Create")
	apiClientId := options.GetValueFromDriverOptions(opts, types.StringType, "api-client-id", "apiClientId").(string)
	apiSecret := options.GetValueFromDriverOptions(opts, types.StringType, "api-secret", "apiSecret").(string)
	config := KConfig{
		Cluster: KConfigCluster{
			Name:       options.GetValueFromDriverOptions(opts, types.StringType, "display-name", "displayName").(string),
			Datacenter: options.GetValueFromDriverOptions(opts, types.StringType, "datacenter").(string),
			SshKey: KConfigClusterSshKey{
				Private: options.GetValueFromDriverOptions(opts, types.StringType, "sshkey-private", "sshkeyPrivate").(string),
				Public:  options.GetValueFromDriverOptions(opts, types.StringType, "sshkey-public", "sshkeyPublic").(string),
			},
			PrivateNetwork: KConfigClusterPrivateNetwork{
				Name: options.GetValueFromDriverOptions(opts, types.StringType, "private-network-name", "privateNetworkName").(string),
			},
		},
	}
	res, err := KCreateCluster(apiClientId, apiSecret, config)
	if err != nil {
		return nil, err
	}
	if res.State != "SUCCESS" {
		errorstr := "unknown error"
		if res.Error != nil {
			errorstr = *res.Error
		}
		return nil, fmt.Errorf("failed to create cluster: %s", errorstr)
	}
	if clusterInfo == nil {
		clusterInfo = &types.ClusterInfo{}
	}
	if clusterInfo.Metadata == nil {
		clusterInfo.Metadata = map[string]string{}
	}
	configJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	clusterInfo.Metadata["kconfig"] = string(configJson)
	clusterInfo.Metadata["apiClientId"] = apiClientId
	clusterInfo.Metadata["apiSecret"] = apiSecret
	return clusterInfo, nil
}

func (d *Driver) Update(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {
	logrus.Info("[Kamatera] Update")
	return clusterInfo, nil
}

func (d *Driver) PostCheck(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Info("[Kamatera] PostCheck")
	kconfigJson := clusterInfo.Metadata["kconfig"]
	var config KConfig
	err := json.Unmarshal([]byte(kconfigJson), &config)
	if err != nil {
		return nil, err
	}
	kubeConfig, err := KGetKubeconfig(clusterInfo.Metadata["apiClientId"], clusterInfo.Metadata["apiSecret"], config)
	if err != nil {
		return nil, err
	}
	clusterInfo.RootCaCertificate = kubeConfig.Clusters[0].Cluster.CertificateAuthorityData
	clusterInfo.Endpoint = kubeConfig.Clusters[0].Cluster.Server
	clusterInfo.ClientCertificate = kubeConfig.Users[0].User.ClientCertificateData
	clusterInfo.ClientKey = kubeConfig.Users[0].User.ClientKeyData
	clusterInfo.Username = kubeConfig.Users[0].Name

	// vvvvvvv TODO vvvvvvv
	clusterInfo.Version = "v1.31.0"
	clusterInfo.NodeCount = 1
	// ^^^^^^^ TODO ^^^^^^^

	caData, err := base64.StdEncoding.DecodeString(clusterInfo.RootCaCertificate)
	if err != nil {
		return nil, err
	}
	clientCertificateData, err := base64.StdEncoding.DecodeString(clusterInfo.ClientCertificate)
	if err != nil {
		return nil, err
	}
	clientKeyData, err := base64.StdEncoding.DecodeString(clusterInfo.ClientKey)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(&rest.Config{
		Host: clusterInfo.Endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   caData,
			CertData: clientCertificateData,
			KeyData:  clientKeyData,
		},
	})
	if err != nil {
		return nil, err
	}
	clusterInfo.ServiceAccountToken, err = GenerateServiceAccountToken(clientset, config.Cluster.Name)
	if err != nil {
		return nil, err
	}
	if cloudcliDebug == "true" {
		logrus.Infof("PostCheck info: %v", *clusterInfo)
	}
	return clusterInfo, nil
}

func (d *Driver) Remove(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	logrus.Info("[Kamatera] Remove")
	return fmt.Errorf("implement me")
}

func (d *Driver) GetVersion(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.KubernetesVersion, error) {
	logrus.Info("[Kamatera] GetVersion")
	// TODO
	return &types.KubernetesVersion{
		Version: "v1.31.0",
	}, nil
}

func (d *Driver) SetVersion(ctx context.Context, clusterInfo *types.ClusterInfo, version *types.KubernetesVersion) error {
	logrus.Info("[Kamatera] SetVersion")
	return fmt.Errorf("implement SetVersion")
}

func (d *Driver) GetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.NodeCount, error) {
	logrus.Info("[Kamatera] GetClusterSize")
	return &types.NodeCount{
		Count: 1,
	}, nil
}

func (d *Driver) SetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo, count *types.NodeCount) error {
	logrus.Info("[Kamatera] SetClusterSize")
	return fmt.Errorf("implement SetClusterSize")
}

func (d *Driver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	return &d.driverCapabilities, nil
}

func (d *Driver) RemoveLegacyServiceAccount(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return nil
}

func (d *Driver) ETCDSave(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	logrus.Info("[Kamatera] ETCDSave")
	return fmt.Errorf("implement ETCDSave")
}

func (d *Driver) ETCDRestore(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) (*types.ClusterInfo, error) {
	logrus.Info("[Kamatera] ETCDRestore")
	return nil, fmt.Errorf("implement ETCDRestore")
}

func (d *Driver) ETCDRemoveSnapshot(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	logrus.Info("[Kamatera] ETCDRemoveSnapshot")
	return fmt.Errorf("implement ETCDRemoveSnapshot")
}

func (d *Driver) GetK8SCapabilities(ctx context.Context, opts *types.DriverOptions) (*types.K8SCapabilities, error) {
	capabilities := &types.K8SCapabilities{
		L4LoadBalancer: &types.LoadBalancerCapabilities{
			Enabled: false,
		},
		NodePoolScalingSupported: false,
	}
	return capabilities, nil
}

func NewDriver() types.Driver {
	driver := &Driver{
		driverCapabilities: types.Capabilities{
			Capabilities: make(map[int64]bool),
		},
	}

	driver.driverCapabilities.AddCapability(types.GetVersionCapability)
	driver.driverCapabilities.AddCapability(types.SetVersionCapability)
	driver.driverCapabilities.AddCapability(types.GetClusterSizeCapability)
	driver.driverCapabilities.AddCapability(types.SetClusterSizeCapability)
	driver.driverCapabilities.AddCapability(types.EtcdBackupCapability)

	return driver
}
