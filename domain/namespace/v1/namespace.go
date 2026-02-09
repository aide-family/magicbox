// Package namespacev1 is the namespace service implementation.
package namespacev1

import (
	context "context"

	"github.com/aide-family/magicbox/enum"
)

type (
	Repository interface {
		SelectNamespace(ctx context.Context, req *SelectNamespaceRequest) (*SelectNamespaceResponse, error)
	}

	SelectNamespaceRequest struct {
		Keyword string            `json:"keyword"`
		Limit   int32             `json:"limit"`
		LastUID int64             `json:"last_uid"`
		Status  enum.GlobalStatus `json:"status"`
	}
	SelectNamespaceResponse struct {
		Items   []*SelectNamespaceItem `json:"items"`
		Total   int64                  `json:"total"`
		LastUID int64                  `json:"last_uid"`
		HasMore bool                   `json:"has_more"`
	}

	SelectNamespaceItem struct {
		Value    int64  `json:"value"`
		Label    string `json:"label"`
		Disabled bool   `json:"disabled"`
		Tooltip  string `json:"tooltip"`
	}
)
