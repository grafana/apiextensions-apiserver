// Package filepath provides filepath storage related utilities.
package filepath

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

type StorageProvider func(s *runtime.Scheme, g genericregistry.RESTOptionsGetter) (rest.Storage, error)

// Object must be implemented by all resources served by the apiserver.
type Object interface {
	// Object allows the apiserver libraries to operate on the Object
	runtime.Object

	// GetObjectMeta returns the object meta reference.
	GetObjectMeta() *metav1.ObjectMeta

	// Scoper is used to qualify the resource as either namespace scoped or non-namespace scoped.
	rest.Scoper

	// New returns a new instance of the resource -- e.g. &v1.Deployment{}
	New() runtime.Object

	// NewList return a new list instance of the resource -- e.g. &v1.DeploymentList{}
	NewList() runtime.Object

	// GetGroupVersionResource returns the GroupVersionResource for this resource.  The resource should
	// be the all lowercase and pluralized kind.s
	GetGroupVersionResource() schema.GroupVersionResource

	// IsStorageVersion returns true if the object is also the internal version -- i.e. is the type defined
	// for the API group an alias to this object.
	// If false, the resource is expected to implement MultiVersionObject interface.
	IsStorageVersion() bool
}

// ABOVE HERE ADDED TO remove dependencies on rest of tilt
//-------------------------
//-------------------------
//-------------------------

// NewJSONFilepathStorageProvider use local host path as persistent layer storage:
//
//   - For namespaced-scoped resources: the resource will be written under the root-path in
//     the following structure:
//
//     -- (root-path) --- /namespace1/ --- resource1
//     |                |
//     |                --- resource2
//     |
//     --- /namespace2/ --- resource3
//
//   - For cluster-scoped resources, there will be no mid-layer folders for namespaces:
//
//     -- (root-path) --- resource1
//     |
//     --- resource2
//     |
//     --- resource3
//
// Args:
//
// fs: An abstraction over the filesystem, so that the JSON can be stored in memory or on-disk.
// watchSet: Storage for watchers to be notified of this resource type. Each type should have its own
//
//	WatchSet, but subresources (like the status subresource) should share a WatchSet with their parent.
func NewJSONFilepathStorageProvider(obj Object, rootPath string, fs FS, watchSet *WatchSet, strategy Strategy) StorageProvider {
	return func(scheme *runtime.Scheme, getter generic.RESTOptionsGetter) (rest.Storage, error) {
		gr := obj.GetGroupVersionResource().GroupResource()
		opt, err := getter.GetRESTOptions(gr)
		if err != nil {
			return nil, err
		}
		codec := opt.StorageConfig.Codec
		return NewFilepathREST(
			fs,
			watchSet,
			strategy,
			gr,
			codec,
			rootPath,
			obj.New,
			obj.NewList,
		), nil
	}
}
