package gormimpl

import (
	apiv1 "github.com/aide-family/magicbox/api/v1"
	"github.com/aide-family/magicbox/domain/namespace/v1/gormimpl/model"
	"github.com/aide-family/magicbox/enum"
)

func ConvertNamespaceItemSelect(namespaceDo *model.Namespace) *apiv1.NamespaceItemSelect {
	if namespaceDo == nil {
		return nil
	}
	return &apiv1.NamespaceItemSelect{
		Value:    namespaceDo.UID.Int64(),
		Label:    namespaceDo.Name,
		Disabled: namespaceDo.DeletedAt.Valid || namespaceDo.Status != enum.GlobalStatus_ENABLED,
		Tooltip:  "",
	}
}

func ConvertNamespaceItem(namespaceDo *model.Namespace) *apiv1.NamespaceItem {
	if namespaceDo == nil {
		return nil
	}
	return &apiv1.NamespaceItem{
		Uid:       namespaceDo.UID.Int64(),
		Name:      namespaceDo.Name,
		Metadata:  namespaceDo.Metadata.Map(),
		Status:    namespaceDo.Status,
		CreatedAt: namespaceDo.CreatedAt.String(),
		UpdatedAt: namespaceDo.UpdatedAt.String(),
	}
}
