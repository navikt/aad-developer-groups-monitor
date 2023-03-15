package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace  = "navikt"
	subsystem  = "aad_developer_groups_monitor"
	labelGroup = "groupName"
)

var developers = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Subsystem: subsystem,
	Name:      "developers",
	Help:      "Number of developers, labeled with group name",
}, []string{labelGroup})

func SetDeveloperCount(numDevelopers int, groupName string) {
	developers.With(prometheus.Labels{
		labelGroup: groupName,
	}).Set(float64(numDevelopers))
}
