package options

import (
	"net"
	"kis/apiserver/storage/storagebackend"
	"github.com/spf13/pflag"
)

type GenericServerRunOptions struct{
	AdvertiseAddress		net.IP
}

func (s *GenericServerRunOptions) AddUniversalFlags(fs *pflag.FlagSet){
	fs.IPVar(&s.AdvertiseAddress, "advertise-address", s.AdvertiseAddress, ""+
		"The IP address on which to advertise the apiserver to members of the cluster. This "+
		"address must be reachable by the rest of the cluster. If blank, the --bind-address "+
		"will be used. If --bind-address is unspecified, the host's default interface will "+
		"be used.")
}
// ServerRunOptions runs a kubernetes api server.
type ServerRunOptions struct {
	GenericServerRunOptions		*GenericServerRunOptions
	Etcd                    *EtcdOptions
}

func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	s.GenericServerRunOptions.AddUniversalFlags(fs)
	s.Etcd.AddFlags(fs)
}



func NewEtcdOptions(backendConfig *storagebackend.Config) *EtcdOptions {
	options := &EtcdOptions{
		StorageConfig:           *backendConfig,
		DefaultStorageMediaType: "application/json",
		DeleteCollectionWorkers: 1,
		EnableGarbageCollection: true,
		EnableWatchCache:        true,
		DefaultWatchCacheSize:   100,
	}
	return options
}
const (
	DefaultEtcdPathPrefix = "/registry"
)

// NewServerRunOptions creates a new ServerRunOptions object with default parameters
func NewServerRunOptions() *ServerRunOptions {
	backendConfig := storagebackend.Config{
		Prefix:DefaultEtcdPathPrefix,
	}
	s := ServerRunOptions{
		GenericServerRunOptions:	&GenericServerRunOptions{},
		Etcd:                 NewEtcdOptions(&backendConfig),
	}

	// Overwrite the default for storage data format.
	s.Etcd.DefaultStorageMediaType = "application/vnd.kubernetes.protobuf"

	return &s
}


type completedServerRunOptions struct {
	*ServerRunOptions
}

func (c *completedServerRunOptions)Validate()(bool){
	return true
}