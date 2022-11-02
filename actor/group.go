package actor

import (
	"github.com/google/uuid"
	"strings"
)

type Group struct {
	ID   uuid.UUID
	Name string
}

func NewGroup(name string) *Group {
	grp := Group{
		ID:   uuid.New(),
		Name: name,
	}
	return &grp
}

func (grp *Group) GetTopic(name string) string {
	str := "group/by-id/" + grp.ID.String() + "/bhv/by-name/" + name
	return strings.ToValidUTF8(str, "--")
}
