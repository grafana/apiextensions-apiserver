package filepath

import (
	"path"

	customStorage "k8s.io/apiextensions-apiserver/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

var _ customStorage.NewStorageFunc = Storage

func Storage(gr schema.GroupResource,
	strategy customStorage.Strategy,
	optsGetter generic.RESTOptionsGetter,
	tableConvertor rest.TableConvertor,
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
) (customStorage.Storage, error) {
	fs := RealFS{}
	ws := NewWatchSet()

	opt, err := optsGetter.GetRESTOptions(gr)
	if err != nil {
		return nil, err
	}
	codec := opt.StorageConfig.Codec
	store := NewFilepathREST(
		fs,
		ws,
		strategy,
		gr,
		codec,
		tableConvertor,
		path.Join("data/k8s/resources", gr.String()),
		newFunc,
		newListFunc,
	)
	return store, nil
}
