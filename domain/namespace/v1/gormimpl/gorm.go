// Package gormimpl is the implementation of the gorm repository for the namespace service.
package gormimpl

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	"github.com/aide-family/magicbox/config"
	"github.com/aide-family/magicbox/connect"
	namespacev1 "github.com/aide-family/magicbox/domain/namespace/v1"
	"github.com/aide-family/magicbox/domain/namespace/v1/gormimpl/query"
	"github.com/aide-family/magicbox/enum"
	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/merr"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
)

func init() {
	namespacev1.RegisterNamespaceV1Factory(config.DomainConfig_GORM, NewGormRepository)
}

func NewGormRepository(c *config.DomainConfig) (namespacev1.Repository, func() error, error) {
	ormConfig := &config.ORMConfig{}
	if pointer.IsNotNil(c.GetOptions()) {
		if err := anypb.UnmarshalTo(c.GetOptions(), ormConfig, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, nil, merr.ErrorInternalServer("unmarshal orm config failed: %v", err)
		}
	}
	db, close, err := connect.NewDB(ormConfig)
	if err != nil {
		return nil, nil, err
	}
	query.SetDefault(db)
	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return nil, nil, err
	}
	return &gormRepository{repoConfig: c, db: db, node: node}, close, nil
}

type gormRepository struct {
	repoConfig *config.DomainConfig
	db         *gorm.DB
	node       *snowflake.Node
}

func (g *gormRepository) GetNamespace(ctx context.Context, uid int64) (*namespacev1.NamespaceModel, error) {
	namespace, err := query.Namespace.WithContext(ctx).Where(query.Namespace.UID.Eq(uid)).First()
	if err != nil {
		return nil, merr.ErrorInternalServer("get namespace failed: %v", err)
	}
	return &namespacev1.NamespaceModel{
		UID:       namespace.UID.Int64(),
		Name:      namespace.Name,
		Metadata:  namespace.Metadata.Map(),
		Status:    enum.GlobalStatus(namespace.Status),
		CreatedAt: namespace.CreatedAt.Unix(),
		UpdatedAt: namespace.UpdatedAt.Unix(),
	}, nil
}

// SelectNamespace implements [namespacev1.Repository].
func (g *gormRepository) SelectNamespace(ctx context.Context, req *namespacev1.SelectNamespaceRequest) (*namespacev1.SelectNamespaceResponse, error) {
	mutation := query.Namespace
	wrappers := mutation.WithContext(ctx)
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(mutation.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status > enum.GlobalStatus_GlobalStatus_UNKNOWN {
		wrappers = wrappers.Where(mutation.Status.Eq(uint8(req.Status)))
	}
	wrappers = wrappers.Order(mutation.UID.Desc())
	total, err := wrappers.Count()
	if err != nil {
		return nil, merr.ErrorInternalServer("count namespace failed: %v", err)
	}
	if req.LastUID > 0 {
		wrappers = wrappers.Where(mutation.UID.Lt(req.LastUID))
	}
	wrappers = wrappers.Limit(int(req.Limit))
	wrappers = wrappers.Select(mutation.UID, mutation.Name, mutation.Status, mutation.DeletedAt)
	queryNamespaces, err := wrappers.Find()
	if err != nil {
		return nil, merr.ErrorInternalServer("select namespace failed: %v", err)
	}
	namespaces := make([]*namespacev1.SelectNamespaceItem, 0, len(queryNamespaces))
	lastUID := int64(0)
	for _, queryNamespace := range queryNamespaces {
		namespace := ConvertNamespaceItemSelect(queryNamespace)
		namespaces = append(namespaces, namespace)
		lastUID = namespace.Value
	}
	return &namespacev1.SelectNamespaceResponse{
		Items:   namespaces,
		Total:   total,
		LastUID: lastUID,
		HasMore: len(queryNamespaces) == int(req.Limit),
	}, nil
}
