package main

import (
	"context"
	"github.com/rancher/kontainer-engine/types"
)

type Driver struct {
	driverCapabilities types.Capabilities
}

func (d Driver) GetDriverCreateOptions(ctx context.Context) (*types.DriverFlags, error) {
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["name"] = &types.Flag{
		Type:  types.StringType,
		Usage: "the internal name of the cluster in Rancher",
	}
	driverFlag.Options["api-client-id"] = &types.Flag{
		Type:     types.StringType,
		Password: true,
		Usage:    "Kamatera API Client ID",
	}
	driverFlag.Options["api-secret"] = &types.Flag{
		Type:     types.StringType,
		Password: true,
		Usage:    "Kamatera API Secret",
	}
	driverFlag.Options["node-pools"] = &types.Flag{
		Type:  types.StringSliceType,
		Usage: "The list of node pools created for the cluster",
	}
	return &driverFlag, nil
}

func (d Driver) GetDriverUpdateOptions(ctx context.Context) (*types.DriverFlags, error) {
	driverFlag := types.DriverFlags{
		Options: make(map[string]*types.Flag),
	}
	driverFlag.Options["node-pools"] = &types.Flag{
		Type:  types.StringSliceType,
		Usage: "The list of node pools created for the cluster",
	}
	return &driverFlag, nil
}

func (d Driver) Create(ctx context.Context, opts *types.DriverOptions, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d Driver) Update(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d Driver) PostCheck(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d Driver) Remove(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	//TODO implement me
	panic("implement me")
}

func (d Driver) GetVersion(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.KubernetesVersion, error) {
	//TODO implement me
	panic("implement me")
}

func (d Driver) SetVersion(ctx context.Context, clusterInfo *types.ClusterInfo, version *types.KubernetesVersion) error {
	//TODO implement me
	panic("implement me")
}

func (d Driver) GetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.NodeCount, error) {
	//TODO implement me
	panic("implement me")
}

func (d Driver) SetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo, count *types.NodeCount) error {
	//TODO implement me
	panic("implement me")
}

func (d Driver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	return &d.driverCapabilities, nil
}

func (d Driver) RemoveLegacyServiceAccount(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	//TODO implement me
	panic("implement me")
}

func (d Driver) ETCDSave(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	//TODO implement me
	panic("implement me")
}

func (d Driver) ETCDRestore(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) (*types.ClusterInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d Driver) ETCDRemoveSnapshot(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	//TODO implement me
	panic("implement me")
}

func (d Driver) GetK8SCapabilities(ctx context.Context, opts *types.DriverOptions) (*types.K8SCapabilities, error) {
	capabilities := &types.K8SCapabilities{
		//L4LoadBalancer: &types.LoadBalancerCapabilities{
		//	Enabled:              true,
		//	Provider:             "NodeBalancer", // what are the options?
		//	ProtocolsSupported:   []string{"TCP", "UDP"},
		//	HealthCheckSupported: true,
		//},
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

	return driver
}
