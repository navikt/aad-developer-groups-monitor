package metrics

import (
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "navikt"
	subsystem = "aad_developer_groups_monitor"
)

var developers = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Subsystem: subsystem,
	Name:      "developers",
	Help:      "Number of developers, labeled with group name",
}, []string{"group_name", "group_id"})

func SetDeveloperCount(numDevelopers int, groupName string, groupID uuid.UUID) {
	developers.With(prometheus.Labels{
		"group_name": groupName,
		"group_id":   groupID.String(),
	}).Set(float64(numDevelopers))
}
