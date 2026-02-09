// Package fileimpl is the implementation of the file repository for the namespace service.
package fileimpl

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/fsnotify/fsnotify"
	klog "github.com/go-kratos/kratos/v2/log"
	"go.yaml.in/yaml/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/aide-family/magicbox/config"
	namespacev1 "github.com/aide-family/magicbox/domain/namespace/v1"
	"github.com/aide-family/magicbox/domain/namespace/v1/fileimpl/model"
	"github.com/aide-family/magicbox/enum"
	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/merr"
	"github.com/aide-family/magicbox/pointer"
)

func init() {
	namespacev1.RegisterNamespaceV1Factory(config.DomainConfig_FILE, NewFileRepository)
}

func NewFileRepository(c *config.DomainConfig) (namespacev1.Repository, func() error, error) {
	fileConfig := &config.FileConfig{}
	if pointer.IsNotNil(c.GetOptions()) {
		if err := anypb.UnmarshalTo(c.GetOptions(), fileConfig, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, nil, merr.ErrorInternalServer("unmarshal file config failed: %v", err)
		}
	}

	// 确保目录存在
	if err := os.MkdirAll(fileConfig.Path, 0o755); err != nil {
		return nil, nil, merr.ErrorInternalServer("create directory failed: %v", err)
	}

	tmpFilepath := filepath.Join(fileConfig.Path, fmt.Sprintf("%s.tmp", fileConfig.Filename))
	filepath := filepath.Join(fileConfig.Path, fileConfig.Filename)
	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return nil, nil, err
	}
	f := &fileRepository{
		repoConfig:      c,
		fileConfig:      fileConfig,
		tmpFilepath:     tmpFilepath,
		filepath:        filepath,
		stopChan:        make(chan struct{}),
		storageInterval: fileConfig.StorageInterval.AsDuration(),
		node:            node,
		namespaces:      make([]*model.NamespaceModel, 0),
	}
	if err := f.load(); err != nil {
		return nil, nil, err
	}
	f.watch()
	return f, func() error {
		close(f.stopChan)
		return f.save()
	}, nil
}

type fileRepository struct {
	repoConfig      *config.DomainConfig
	fileConfig      *config.FileConfig
	tmpFilepath     string
	filepath        string
	mu              sync.RWMutex
	namespaces      []*model.NamespaceModel
	nextID          uint32
	stopChan        chan struct{}
	storageInterval time.Duration
	changed         bool
	node            *snowflake.Node
}

func (f *fileRepository) load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 如果文件不存在，初始化为空列表
	if _, err := os.Stat(f.filepath); os.IsNotExist(err) {
		f.namespaces = make([]*model.NamespaceModel, 0)
		f.nextID = 0
		return nil
	}

	file, err := os.Open(f.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var namespaces []*model.NamespaceModel
	if err := yaml.NewDecoder(file).Decode(&namespaces); err != nil {
		// 如果文件为空，EOF 是正常情况，初始化为空列表
		if err == io.EOF {
			f.namespaces = make([]*model.NamespaceModel, 0)
			f.nextID = 0
			return nil
		}
		return err
	}
	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].ID < namespaces[j].ID
	})

	f.nextID = namespaces[len(namespaces)-1].ID
	for _, namespace := range namespaces {
		if namespace.ID == 0 {
			f.nextID++
			namespace.ID = f.nextID
		}
		// 确保已删除的 namespace 不会被重置 UID
		if namespace.UID == 0 {
			namespace.UID = f.node.Generate().Int64()
		}
	}

	f.namespaces = namespaces
	return nil
}

func (f *fileRepository) save() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.changed = false
	file, err := os.Create(f.tmpFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := yaml.NewEncoder(file).Encode(f.namespaces); err != nil {
		return err
	}
	if err := os.Rename(f.tmpFilepath, f.filepath); err != nil {
		return err
	}
	klog.Debugw("msg", "save namespaces to file", "filepath", f.filepath)
	return nil
}

func (f *fileRepository) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Errorw("msg", "create watcher failed", "error", err)
		return
	}
	defer watcher.Close()
	watcher.Add(f.filepath)
	go func() {
		ticker := time.NewTicker(f.storageInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if f.changed {
					f.save()
				}
			case err := <-watcher.Errors:
				if err != nil {
					klog.Warnw("msg", "watch file failed", "error", err)
				}
			case <-f.stopChan:
				klog.Debugw("msg", "stop watch namespaces")
				return
			}
		}
	}()
}

// SelectNamespace implements [namespacev1.Repository].
func (f *fileRepository) SelectNamespace(ctx context.Context, req *namespacev1.SelectNamespaceRequest) (*namespacev1.SelectNamespaceResponse, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	namespaces := make([]*namespacev1.SelectNamespaceItem, 0, len(f.namespaces))
	sort.Slice(f.namespaces, func(i, j int) bool {
		return f.namespaces[i].UID > f.namespaces[j].UID
	})
	count := 0
	lessFunc := func(i, j int64) bool {
		return i > j
	}
	for _, namespace := range f.namespaces {
		if req.Status > enum.GlobalStatus_GlobalStatus_UNKNOWN && namespace.Status != req.Status {
			continue
		}
		if req.Keyword != "" && !strings.Contains(namespace.Name, req.Keyword) {
			continue
		}
		if lessFunc(namespace.UID, req.LastUID) {
			continue
		}
		count++
		namespaces = append(namespaces, convertNamespaceItemSelect(namespace))
		if count >= int(req.Limit) {
			break
		}
	}

	return &namespacev1.SelectNamespaceResponse{
		Items:   namespaces[:req.Limit],
		Total:   int64(len(namespaces)),
		LastUID: namespaces[len(namespaces)-1].Value,
		HasMore: count == int(req.Limit),
	}, nil
}
