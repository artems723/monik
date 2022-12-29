package storage

type Repository interface {
	GetMetric(agentID, metricName string) (string, bool)
	WriteMetric(agentID, metricName, metricValue string)
}
