package resourceclient

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateOrUpdate[S client.Object](ctx context.Context, c client.Client, objectType, obj S) error {
	namespacedName := types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
	err := c.Get(ctx, namespacedName, objectType)
	if err != nil {
		if errors.IsNotFound(err) {
			if err := c.Create(ctx, obj); err != nil {
				return err
			}
			return nil
		}
		return err
	} else {
		if err := c.Update(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}
