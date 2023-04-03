package filepath

import (
	"path"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	customStorage "k8s.io/apiextensions-apiserver/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

func Storage(gr schema.GroupResource,
	kind, listKind schema.GroupVersionKind,
	strategy customStorage.Strategy,
	optsGetter generic.RESTOptionsGetter,
	tableConvertor rest.TableConvertor,
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
		func() runtime.Object { return &apiextensions.CustomResourceDefinition{} },
		func() runtime.Object { return &apiextensions.CustomResourceDefinitionList{} },
	)
	return store, nil
}
