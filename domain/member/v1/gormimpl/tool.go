package gormimpl

import (
	"strings"
	"time"

	apiv1 "github.com/aide-family/magicbox/api/v1"
	"github.com/aide-family/magicbox/domain/member/v1/gormimpl/model"
)

func convertMemberItem(m *model.Member) *apiv1.MemberItem {
	if m == nil {
		return nil
	}
	return &apiv1.MemberItem{
		Uid:       m.ID.Int64(),
		Email:     m.Email,
		Phone:     m.Phone,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Format(time.DateTime),
		UpdatedAt: m.UpdatedAt.Format(time.DateTime),
		UserUID:   m.UserUID.Int64(),
		Name:      m.Name,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		Remark:    m.Remark,
	}
}

func convertMemberSelectItem(m *model.Member) *apiv1.SelectMemberItem {
	if m == nil {
		return nil
	}
	return &apiv1.SelectMemberItem{
		Value:    m.ID.Int64(),
		Label:    strings.Join([]string{m.Name, "(", m.Nickname, ")"}, " "),
		Disabled: false,
		Tooltip:  m.Remark,
	}
}
