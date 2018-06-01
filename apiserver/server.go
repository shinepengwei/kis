package apiserver

import (
	"github.com/spf13/cobra"
	"kis/apiserver/storage"
	"kis/master"
	"kis/apiserver/options"
	"kis/apiserver/genericapiserver"
	"net"
)

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.

type completedServerRunOptions struct {
	*options.ServerRunOptions
}

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "kis-apiserver",
		Long: `The kis API server test.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// set default options
			completedOptions, err := Complete(s)
			if err != nil {
				return err
			}


			stopCh := SetupSignalHandler()
			return Run(completedOptions, stopCh)
		},
	}
	s.AddFlags(cmd.Flags())

	return cmd
}

func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	//apiserverConfig := BuildStorageFactory(&completeOptions)
	server, err := CreateServerChain(completeOptions, stopCh)
	if err != nil {
		return err
	}
	server.Run(stopCh)
	return nil
}

func CreateServerChain(completedOptions completedServerRunOptions, stopCh <-chan struct{}) (*master.Master, error){
	kubeAPIServerConfig, err := CreateKubeAPIServerConfig(completedOptions)
	if err != nil {
		return nil, err
	}
	m := CreateKubeAPIServer(kubeAPIServerConfig)
	return m,nil
}


func CreateKubeAPIServerConfig(completedOptions completedServerRunOptions) (*master.Config,error){
	apiserverConfig := BuildGenericAPIServerConfig(completedOptions.ServerRunOptions)

	storageFactory := BuildStorageFactory(completedOptions.ServerRunOptions)
	config := &master.Config{
		GenericConfig: apiserverConfig,
		ExtraConfig: master.ExtraConfig{
			StorageFactory:          storageFactory,
		},
	}
	return config, nil
}

func CreateKubeAPIServer(kubeAPIServerConfig *master.Config) (* master.Master){
	m := kubeAPIServerConfig.New()
	return m

}
func BuildStorageFactory(s *options.ServerRunOptions) *storage.DefaultStorageFactory{
	return storage.NewStorageFactory(s.Etcd.StorageConfig,s.Etcd.DefaultStorageMediaType)
}

func BuildGenericAPIServerConfig(s *options.ServerRunOptions)(*genericapiserver.Config){
	apiserverConfig := genericapiserver.NewGenericAPIServerConfig()
	storageFactory := BuildStorageFactory(s)
	s.Etcd.ApplyWithStorageFactoryTo(storageFactory, apiserverConfig)
	return apiserverConfig
}

func Complete(s *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	// set defaults
	s.GenericServerRunOptions.AdvertiseAddress = net.IP{0}
	s.Etcd.StorageConfig.ServerList = []string{"http://localhost:2379"}

	options.ServerRunOptions = s
	return options, nil
}
