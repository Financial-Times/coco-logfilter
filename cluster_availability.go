package logfilter

import (
	"errors"
	"github.com/domainr/dnsr"
	"strings"
	"time"
)

// CACHE_TIME time to cache the cluster status
const CACHE_TIME = 30

type Cluster struct {
	readDNS string
	tag     string
}

type ClusterService interface {
	IsActive() (bool, error)
}

type MonitoredClusterService struct {
	instance       Cluster
	cachedStatus   bool
	cacheTimestamp time.Time
}

func NewMonitoredClusterService(readDNS string, tag string) MonitoredClusterService {
	instance := Cluster{readDNS: readDNS, tag: tag}
	return MonitoredClusterService{instance: instance}
}

func (mc MonitoredClusterService) IsActive() (bool, error) {
	// if the read address contains the cluster tag (and implicitly the region)
	// than this means that there is no failover mechanism in place
	if strings.Contains(mc.instance.readDNS, mc.instance.tag) {
		return true, nil
	}

	if time.Now().Sub(mc.cacheTimestamp).Seconds() < CACHE_TIME {
		return mc.cachedStatus, nil
	}

	resolver := dnsr.New(5)
	cNames := resolver.Resolve(mc.instance.readDNS, "CNAME")
	mc.cacheTimestamp = time.Now()
	if len(cNames) > 0 {
		isActive := strings.Contains(cNames[0].Value, mc.instance.tag)
		mc.cachedStatus = isActive
		return isActive, nil
	}
	mc.cachedStatus = false
	return false, errors.New("address could not be resolved, maybe it is invalid")
}
