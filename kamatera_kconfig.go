package main

type KConfig struct {
	Cluster   KConfigCluster             `json:"cluster"`
	NodePools map[string]KConfigNodePool `json:"node-pools"`
}

type KConfigCluster struct {
	Name           string                       `json:"name"`
	Datacenter     string                       `json:"datacenter"`
	SshKey         KConfigClusterSshKey         `json:"ssh-key"`
	PrivateNetwork KConfigClusterPrivateNetwork `json:"private-network"`
}

type KConfigClusterSshKey struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

type KConfigClusterPrivateNetwork struct {
	Name string `json:"name"`
}

type KConfigNodePool struct {
	Nodes []int `json:"nodes"`
}
