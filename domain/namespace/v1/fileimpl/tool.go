package fileimpl

import (
	"encoding/json"

	namespacev1 "github.com/aide-family/magicbox/domain/namespace/v1"
	"github.com/aide-family/magicbox/domain/namespace/v1/fileimpl/model"
	"github.com/aide-family/magicbox/enum"
)

func convertNamespaceItemSelect(namespaceModel *model.NamespaceModel) *namespacev1.SelectNamespaceItem {
	metadata, _ := json.Marshal(namespaceModel.Metadata)
	return &namespacev1.SelectNamespaceItem{
		Value:    namespaceModel.UID,
		Label:    namespaceModel.Name,
		Disabled: namespaceModel.DeletedAt != 0 || namespaceModel.Status != enum.GlobalStatus_ENABLED,
		Tooltip:  string(metadata),
	}
}
