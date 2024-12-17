package deployment

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	biblobstore "github.com/cloudfoundry/bosh-cli/v7/blobstore"
	bicloud "github.com/cloudfoundry/bosh-cli/v7/cloud"
	bidisk "github.com/cloudfoundry/bosh-cli/v7/deployment/disk"
	biinstance "github.com/cloudfoundry/bosh-cli/v7/deployment/instance"
	bideplmanifest "github.com/cloudfoundry/bosh-cli/v7/deployment/manifest"
	bivm "github.com/cloudfoundry/bosh-cli/v7/deployment/vm"
	bistemcell "github.com/cloudfoundry/bosh-cli/v7/stemcell"
	biui "github.com/cloudfoundry/bosh-cli/v7/ui"
)

type Deployer interface {
	Deploy(
		cloud bicloud.Cloud,
		deploymentManifest bideplmanifest.Manifest,
		cloudStemcell bistemcell.CloudStemcell,
		vmManager bivm.Manager,
		blobstore biblobstore.Blobstore,
		skipDrain bool,
		diskCIDs []string,
		deployStage biui.Stage,
	) (Deployment, error)
}

type deployer struct {
	vmManagerFactory       bivm.ManagerFactory
	instanceManagerFactory biinstance.ManagerFactory
	deploymentFactory      Factory
	logger                 boshlog.Logger
	logTag                 string
}

func NewDeployer(
	vmManagerFactory bivm.ManagerFactory,
	instanceManagerFactory biinstance.ManagerFactory,
	deploymentFactory Factory,
	logger boshlog.Logger,
) Deployer {
	return &deployer{
		vmManagerFactory:       vmManagerFactory,
		instanceManagerFactory: instanceManagerFactory,
		deploymentFactory:      deploymentFactory,
		logger:                 logger,
		logTag:                 "deployer",
	}
}

func (d *deployer) Deploy(
	cloud bicloud.Cloud,
	deploymentManifest bideplmanifest.Manifest,
	cloudStemcell bistemcell.CloudStemcell,
	vmManager bivm.Manager,
	blobstore biblobstore.Blobstore,
	skipDrain bool,
	diskCIDs []string,
	deployStage biui.Stage,
) (Deployment, error) {
	instanceManager := d.instanceManagerFactory.NewManager(cloud, vmManager, blobstore)

	pingTimeout := 10 * time.Second
	pingDelay := 500 * time.Millisecond
	if err := instanceManager.DeleteAll(pingTimeout, pingDelay, skipDrain, deployStage); err != nil {
		return nil, err
	}

	instances, disks, err := d.createAllInstances(deploymentManifest, instanceManager, cloudStemcell, diskCIDs, deployStage)
	if err != nil {
		return nil, err
	}

	stemcells := []bistemcell.CloudStemcell{cloudStemcell}
	return d.deploymentFactory.NewDeployment(instances, disks, stemcells), nil
}

func (d *deployer) createAllInstances(
	deploymentManifest bideplmanifest.Manifest,
	instanceManager biinstance.Manager,
	cloudStemcell bistemcell.CloudStemcell,
	diskCIDs []string,
	deployStage biui.Stage,
) ([]biinstance.Instance, []bidisk.Disk, error) {
	instances := []biinstance.Instance{}
	disks := []bidisk.Disk{}

	if len(deploymentManifest.Jobs) != 1 {
		return instances, disks, bosherr.Errorf("There must only be one job, found %d", len(deploymentManifest.Jobs))
	}

	for _, jobSpec := range deploymentManifest.Jobs {
		if jobSpec.Instances != 1 {
			return instances, disks, bosherr.Errorf("Job '%s' must have only one instance, found %d", jobSpec.Name, jobSpec.Instances)
		}
		for instanceID := 0; instanceID < jobSpec.Instances; instanceID++ {
			instance, instanceDisks, err := instanceManager.Create(jobSpec.Name, instanceID, deploymentManifest, cloudStemcell, diskCIDs, deployStage)
			if err != nil {
				return instances, disks, bosherr.WrapErrorf(err, "Creating instance '%s/%d'", jobSpec.Name, instanceID)
			}
			instances = append(instances, instance)
			disks = append(disks, instanceDisks...)

			err = instance.UpdateJobs(deploymentManifest, deployStage)
			if err != nil {
				return instances, disks, err
			}
		}
	}

	return instances, disks, nil
}
