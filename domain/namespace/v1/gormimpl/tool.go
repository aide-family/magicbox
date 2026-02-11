package gormimpl

import (
	namespacev1 "github.com/aide-family/magicbox/domain/namespace/v1"
	"github.com/aide-family/magicbox/domain/namespace/v1/gormimpl/model"
	"github.com/aide-family/magicbox/enum"
)

func ConvertNamespaceItemSelect(namespaceDo *model.Namespace) *namespacev1.SelectNamespaceItem {
	return &namespacev1.SelectNamespaceItem{
		Value:    namespaceDo.UID.Int64(),
		Label:    namespaceDo.Name,
		Disabled: namespaceDo.DeletedAt.Valid || namespaceDo.Status != uint8(enum.GlobalStatus_ENABLED),
		Tooltip:  "",
	}
}
