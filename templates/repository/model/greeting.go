package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kumparan/go-lib/utils"
	log "github.com/sirupsen/logrus"

	pb "gitlab.kumparan.com/yowez/skeleton-service/pb/skeleton"
)

// Greeting :nodoc:
type Greeting struct {
	ID        int64     `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewFromProto :nodoc:
func NewFromProto(p *pb.Greeting) *Greeting {
	g := &Greeting{
		ID:   p.GetId(),
		Name: p.GetName(),
	}

	if p.GetCreatedAt() != "" {
		createdAt, err := time.Parse(time.RFC3339, p.GetCreatedAt())
		if err != nil {
			log.WithField("p.GetCreatedAt()", p.GetCreatedAt()).Error(err)
		}
		g.CreatedAt = createdAt
	}

	if p.GetUpdatedAt() != "" {
		UpdatedAt, err := time.Parse(time.RFC3339, p.GetUpdatedAt())
		if err != nil {
			log.WithField("p.GetUpdatedAt()", p.GetUpdatedAt()).Error(err)
		}
		g.UpdatedAt = UpdatedAt
	}

	return g
}

// NewGreetingFromInterface converts interface to greeting model.
func NewGreetingFromInterface(i interface{}) (m *Greeting, err error) {
	bt, _ := i.([]byte)

	err = json.Unmarshal(bt, &m)
	if err != nil {
		log.WithField("i", utils.Dump(i)).Error(err)
	}

	return
}

// NewGreetingCacheKeyByID :nodoc:
func NewGreetingCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:greeting:%d", id)
}

// GetName :nodoc:
func (g *Greeting) GetName() string {
	return g.Name
}

// ToProto :nodoc:
func (g *Greeting) ToProto() *pb.Greeting {
	return &pb.Greeting{
		Id:        g.ID,
		Name:      g.Name,
		CreatedAt: g.CreatedAt.Format(time.RFC3339),
		UpdatedAt: g.UpdatedAt.Format(time.RFC3339),
	}
}

// ImmutableColumns :nodoc:
func (g *Greeting) ImmutableColumns() []string {
	return []string{"created_at"}
}
