package gormimpl

import (
	"encoding/json"

	namespacev1 "github.com/aide-family/magicbox/domain/namespace/v1"
	"github.com/aide-family/magicbox/domain/namespace/v1/gormimpl/model"
	"github.com/aide-family/magicbox/enum"
)

func ConvertNamespaceItemSelect(namespaceDo *model.Namespace) *namespacev1.SelectNamespaceItem {
	metadata, err := json.Marshal(namespaceDo.Metadata.Map())
	if err != nil {
		return nil
	}
	return &namespacev1.SelectNamespaceItem{
		Value:    namespaceDo.UID.Int64(),
		Label:    namespaceDo.Name,
		Disabled: namespaceDo.DeletedAt.Valid || namespaceDo.Status != uint8(enum.GlobalStatus_ENABLED),
		Tooltip:  string(metadata),
	}
}
