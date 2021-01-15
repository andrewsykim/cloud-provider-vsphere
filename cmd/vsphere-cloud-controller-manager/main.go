/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// The external controller manager is responsible for running controller loops that
// are cloud provider dependent. It uses the API to listen to new events on resources.

package main

import (
	goflag "flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/cloud-provider-vsphere/pkg/cloudprovider/vsphere"
	"k8s.io/cloud-provider-vsphere/pkg/cloudprovider/vsphere/loadbalancer"
	"k8s.io/cloud-provider/app"
	"k8s.io/cloud-provider/options"
	"k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // for client metrics registration
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	klog "k8s.io/klog/v2"
)

// AppName is the full name of this CCM
const AppName string = "vsphere-cloud-controller-manager"

var version string

func main() {
	loadbalancer.Version = version
	loadbalancer.AppName = AppName

	rand.Seed(time.Now().UTC().UnixNano())

	s, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	// There's a bug in k8s.io/cloud-provider where validation happens prior to
	// flags being read. This is a temporary work around to that bug since we only
	// accept --cloud-provider=vsphere anyways. This should be removed once
	// https://github.com/kubernetes/cloud-provider/issues/44 is fixed.
	s.KubeCloudShared.CloudProvider.Name = vsphere.ProviderName

	c, err := s.Config([]string{}, app.ControllersDisabledByDefault.List())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// initialize cloud provider with the cloud provider name and config file provided
	cloud, err := cloudprovider.InitCloudProvider(vsphere.ProviderName, c.ComponentConfig.KubeCloudShared.CloudProvider.CloudConfigFile)
	if err != nil {
		klog.Fatalf("Cloud provider could not be initialized: %v", err)
	}
	if cloud == nil {
		klog.Fatalf("cloud provider is nil")
	}

	if !cloud.HasClusterID() {
		if c.ComponentConfig.KubeCloudShared.AllowUntaggedCloud {
			klog.Warning("detected a cluster without a ClusterID.  A ClusterID will be required in the future.  Please tag your cluster to avoid any future issues")
		} else {
			klog.Fatalf("no ClusterID found.  A ClusterID is required for the cloud provider to function properly.  This check can be bypassed by setting the allow-untagged-cloud option")
		}
	}

	// Initialize the cloud provider with a reference to the clientBuilder
	cloud.Initialize(c.ClientBuilder, make(chan struct{}))
	// Set the informer on the user cloud object
	if informerUserCloud, ok := cloud.(cloudprovider.InformerUser); ok {
		informerUserCloud.SetInformers(c.SharedInformers)
	}

	controllerInitializers := app.DefaultControllerInitializers(c.Complete(), cloud)
	command := app.NewCloudControllerManagerCommand(s, c, controllerInitializers)

	// TODO: once we switch everything over to Cobra commands, we can go back to calling
	// utilflag.InitFlags() (by removing its pflag.Parse() call). For now, we have to set the
	// normalize func and add the go flag set by hand.
	pflag.CommandLine.SetNormalizeFunc(flag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	// utilflag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	klog.Infof("vsphere-cloud-controller-manager version: %s", version)

	// Set cloud-provider flag to vsphere
	command.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "cloud-provider" {
			flag.Value.Set(vsphere.ProviderName)
			flag.DefValue = vsphere.ProviderName
			return
		}
	})

	var versionFlag *pflag.Value
	var clusterNameFlag *pflag.Value
	pflag.CommandLine.VisitAll(func(flag *pflag.Flag) {
		switch flag.Name {
		case "cluster-name":
			clusterNameFlag = &flag.Value
		case "version":
			versionFlag = &flag.Value
		}
	})

	command.Use = AppName
	innerRun := command.Run
	command.Run = func(cmd *cobra.Command, args []string) {
		if versionFlag != nil && (*versionFlag).String() != "false" {
			fmt.Printf("%s %s\n", AppName, version)
			os.Exit(0)
		}
		if clusterNameFlag != nil {
			loadbalancer.ClusterName = (*clusterNameFlag).String()
		}
		innerRun(cmd, args)
	}

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
