package operatorsource

import (
	"context"
	"errors"
	"fmt"
	"strings"

	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	marketplace "github.com/operator-framework/operator-marketplace/pkg/apis/operators/v1"
	"github.com/operator-framework/operator-marketplace/pkg/appregistry"
	interface_client "github.com/operator-framework/operator-marketplace/pkg/client"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/phase"
	"github.com/operator-framework/operator-marketplace/pkg/registry"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewConfiguringReconciler returns a Reconciler that reconciles
// an OperatorSource object in "Configuring" phase.
func NewConfiguringReconciler(logger *log.Entry, factory appregistry.ClientFactory, datastore datastore.Writer, reader datastore.Reader, client client.Client, refresher PackageRefreshNotificationSender) Reconciler {
	return NewConfiguringReconcilerWithInterfaceClient(logger, factory, datastore, reader, interface_client.NewClient(client), refresher)
}

// NewConfiguringReconcilerWithInterfaceClient returns a configuring
// Reconciler that reconciles an OperatorSource object in "Configuring"
// phase. It uses the Client interface which is a wrapper to the raw
// client provided by the operator-sdk, instead of the raw client itself.
// Using this interface facilitates mocking of kube client interaction
// with the cluster, while using fakeclient during unit testing.
func NewConfiguringReconcilerWithInterfaceClient(logger *log.Entry, factory appregistry.ClientFactory, datastore datastore.Writer, reader datastore.Reader, client interface_client.Client, refresher PackageRefreshNotificationSender) Reconciler {
	return &configuringReconciler{
		logger:    logger,
		factory:   factory,
		datastore: datastore,
		client:    client,
		refresher: refresher,
		builder:   &CatalogSourceConfigBuilder{},
		reader:    reader,
	}
}

// configuringReconciler is an implementation of Reconciler interface that
// reconciles an OperatorSource object in "Configuring" phase.
type configuringReconciler struct {
	logger    *log.Entry
	factory   appregistry.ClientFactory
	datastore datastore.Writer
	client    interface_client.Client
	refresher PackageRefreshNotificationSender
	builder   *CatalogSourceConfigBuilder
	reader    datastore.Reader
}

// Reconcile reconciles an OperatorSource object that is in "Configuring" phase.
// It ensures that a corresponding CatalogSourceConfig object exists.
//
// in represents the original OperatorSource object received from the sdk
// and before reconciliation has started.
//
// out represents the OperatorSource object after reconciliation has completed
// and could be different from the original. The OperatorSource object received
// (in) should be deep copied into (out) before changes are made.
//
// nextPhase represents the next desired phase for the given OperatorSource
// object. If nil is returned, it implies that no phase transition is expected.
//
// Upon success, it returns "Succeeded" as the next and final desired phase.
// On error, the function returns "Failed" as the next desied phase
// and Message is set to appropriate error message.
//
// If the corresponding CatalogSourceConfig object already exists
// then no further action is taken.
func (r *configuringReconciler) Reconcile(ctx context.Context, in *marketplace.OperatorSource) (out *marketplace.OperatorSource, nextPhase *marketplace.Phase, err error) {
	if in.GetCurrentPhaseName() != phase.Configuring {
		err = phase.ErrWrongReconcilerInvoked
		return
	}

	out = in

	r.logger.Infof("Downloading metadata from Namespace [%s] of [%s]", in.Spec.RegistryNamespace, in.Spec.Endpoint)

	metadata, err := r.getManifestMetadata(&in.Spec, in.Namespace)
	if err != nil {
		nextPhase = phase.GetNextWithMessage(phase.Configuring, err.Error())
		return
	}

	if len(metadata) == 0 {
		err = errors.New("The OperatorSource endpoint returned an empty manifest list")

		// Moving it to 'Failed' phase since human intervention is required to
		// resolve this situation. As soon as the user pushes new operator
		// manifest(s) registry sync will detect a new release and will trigger
		// a new reconciliation.
		nextPhase = phase.GetNextWithMessage(phase.Failed, err.Error())
		return
	}

	r.logger.Infof("%d manifest(s) scheduled for download in the operator-registry pod", len(metadata))

	isResyncNeeded, err := r.writeMetadataToDatastore(in, out, metadata)
	if err != nil {
		// No operator metadata was written, move to Failed phase.
		nextPhase = phase.GetNextWithMessage(phase.Failed, err.Error())
		return
	}

	// Now that we have updated the datastore, let's check if the opsrc is new.
	// If it is, let's force a resync for CatalogSourceConfig.
	if isResyncNeeded {
		r.logger.Info("New opsrc detected. Refreshing catalogsourceconfigs.")
		r.refresher.SendRefresh()
	}

	packages := r.datastore.GetPackageIDsByOperatorSource(out.GetUID())
	out.Status.Packages = packages

	cscCreate := new(CatalogSourceConfigBuilder).WithTypeMeta().
		WithNamespacedName(in.Namespace, in.Name).
		WithLabels(in.GetLabels()).
		WithSpec(in.Namespace, packages, in.Spec.DisplayName, in.Spec.Publisher).
		WithOwnerLabel(in).
		CatalogSourceConfig()

	err = r.reconcileCatalogSource(cscCreate)
	if err != nil {
		nextPhase = phase.GetNextWithMessage(phase.Configuring, err.Error())
		return
	}

	r.ensurePackagesInStatus(cscCreate)

	nextPhase = phase.GetNext(phase.Succeeded)
	return
}

// reconcileCatalogSource ensures a CatalogSource exists with all the
// resources it requires.
func (r *configuringReconciler) reconcileCatalogSource(csc *marketplace.CatalogSourceConfig) error {
	// Ensure that the packages in the spec are available in the datastore
	err := r.checkPackages(csc)
	if err != nil {
		return err
	}

	// Ensure that a registry deployment is available
	existing_registry := registry.NewRegistry(r.logger, r.client, r.reader, csc, registry.RegistryServerImage)
	err = existing_registry.Ensure()
	if err != nil {
		return err
	}

	// Check if the CatalogSource already exists
	catalogSourceGet := new(registry.CatalogSourceBuilder).WithTypeMeta().CatalogSource()
	key := client.ObjectKey{
		Name:      csc.Name,
		Namespace: csc.Spec.TargetNamespace,
	}
	err = r.client.Get(context.TODO(), key, catalogSourceGet)

	return err
}

// getManifestMetadata gets the package metadata from the OperatorSource endpoint.
// It returns the list of packages to be written to the OperatorSource status. error is set
// when there is an issue downloading the metadata. In that case the list of packages
// will be empty.
func (r *configuringReconciler) getManifestMetadata(spec *marketplace.OperatorSourceSpec, namespace string) ([]*datastore.RegistryMetadata, error) {

	metadata := make([]*datastore.RegistryMetadata, 0)

	options, err := SetupAppRegistryOptions(r.client, spec, namespace)
	if err != nil {
		return metadata, err
	}

	registry, err := r.factory.New(options)
	if err != nil {
		return metadata, err
	}

	metadata, err = registry.ListPackages(spec.RegistryNamespace)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}

// writeMetadataToDatastore checks to see if there are any existing metadata
// before we write to the datastore. If there are not, we are assuming
// this is a new OperatorSource and in this case we should force all
// CatalogSourceConfigs to compare their versions to what's in the datastore
// after we update it. The function returns the whether a resync is needed and an error
func (r *configuringReconciler) writeMetadataToDatastore(in *marketplace.OperatorSource, out *marketplace.OperatorSource, metadata []*datastore.RegistryMetadata) (bool, error) {

	preUpdateDatastorePackageList := r.datastore.GetPackageIDsByOperatorSource(out.GetUID())

	count, err := r.datastore.Write(in, metadata)
	if err != nil {
		if count == 0 {
			return preUpdateDatastorePackageList == "", err
		}
		// Ignore faulty operator metadata
		r.logger.Infof("There were some faulty operator metadata, errors - %v", err)
		err = nil
	}

	r.logger.Infof("Successfully downloaded %d operator metadata", count)
	return preUpdateDatastorePackageList == "", err
}

// updateExistingCatalogSourceConfig updates an existing CatalogSourceConfig
// when the OperatorSource that owns it is updated in any way.
func (r *configuringReconciler) updateExistingCatalogSourceConfig(ctx context.Context, in *marketplace.OperatorSource, packages string) error {

	cscNamespacedName := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	cscExisting := marketplace.CatalogSourceConfig{}
	err := r.client.Get(ctx, cscNamespacedName, &cscExisting)
	if err != nil {
		return err
	}

	cscExisting.EnsureGVK()

	builder := CatalogSourceConfigBuilder{object: cscExisting}
	cscUpdate := builder.WithSpec(in.Namespace, packages, in.Spec.DisplayName, in.Spec.Publisher).
		WithLabels(in.GetLabels()).
		WithOwnerLabel(in).
		CatalogSourceConfig()
	// Update the CatalogSource if it exists else create one.
	if err == nil {
		catalogSourceGet.Spec.Address = existing_registry.GetAddress()
		r.logger.Infof("Updating CatalogSource %s", catalogSourceGet.Name)
		err = r.client.Update(context.TODO(), catalogSourceGet)
		if err != nil {
			r.logger.Errorf("Failed to update CatalogSource : %v", err)
			return err
		}
		r.logger.Infof("Updated CatalogSource %s", catalogSourceGet.Name)
	} else {
		// Create the CatalogSource structure
		catalogSource := newCatalogSource(csc, existing_registry.GetAddress())
		r.logger.Infof("Creating CatalogSource %s", catalogSource.Name)
		err = r.client.Create(context.TODO(), catalogSource)
		if err != nil && !errors.IsAlreadyExists(err) {
			r.logger.Errorf("Failed to create CatalogSource : %v", err)
			return err
		}
		r.logger.Infof("Created CatalogSource %s", catalogSource.Name)
	}

	return nil
}

// ensurePackagesInStatus makes sure that the csc's status.PackageRepositioryVersions
// field is updated at the end of the configuring phase if successful. It iterates
// over the list of packages and creates a new map of PackageName:Version for each
// package in the spec.
func (r *configuringReconciler) ensurePackagesInStatus(csc *marketplace.CatalogSourceConfig) {
	newPackageRepositioryVersions := make(map[string]string)
	packageIDs := csc.GetPackageIDs()
	for _, packageID := range packageIDs {
		version, err := r.reader.ReadRepositoryVersion(packageID)
		if err != nil {
			r.logger.Errorf("Failed to find package: %v", err)
			version = "-1"
		}

		newPackageRepositioryVersions[packageID] = version
	}

	csc.Status.PackageRepositioryVersions = newPackageRepositioryVersions
}

// checkPackages returns an error if there are packages missing from the
// datastore but listed in the spec.
func (r *configuringReconciler) checkPackages(csc *marketplace.CatalogSourceConfig) error {
	missingPackages := []string{}
	packageIDs := csc.GetPackageIDs()
	for _, packageID := range packageIDs {
		if _, err := r.reader.Read(packageID); err != nil {
			missingPackages = append(missingPackages, packageID)
			continue
		}
	}

	if len(missingPackages) > 0 {
		return fmt.Errorf(
			"Still resolving package(s) - %s. Please make sure these are valid packages.",
			strings.Join(missingPackages, ","),
		)
	}
	return nil
}

// newCatalogSource returns a CatalogSource object.
func newCatalogSource(csc *marketplace.CatalogSourceConfig, address string) *olm.CatalogSource {
	builder := new(registry.CatalogSourceBuilder).
		WithOwnerLabel(csc).
		WithMeta(csc.Name, csc.Spec.TargetNamespace).
		WithSpec(olm.SourceTypeGrpc, address, csc.Spec.DisplayName, csc.Spec.Publisher)

	// Check if the operatorsource.DatastoreLabel is "true" which indicates that
	// the CatalogSource is the datastore for an OperatorSource. This is a hint
	// for us to set the "olm-visibility" label in the CatalogSource so that it
	// is not visible in the OLM Packages UI. In addition we will set the
	// "openshift-marketplace" label which will be used by the Marketplace UI
	// to filter out global CatalogSources.
	cscLabels := csc.ObjectMeta.GetLabels()
	datastoreLabel, found := cscLabels[datastore.DatastoreLabel]
	if found && strings.ToLower(datastoreLabel) == "true" {
		builder.WithOLMLabels(cscLabels)
	}

	return builder.CatalogSource()
}
