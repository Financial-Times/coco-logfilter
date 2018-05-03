package main

import (
	"errors"
	"github.com/domainr/dnsr"
	"strings"
	"time"
)

// cacheTime time to cache the cluster status
const cacheTime = 30

type cluster struct {
	dnsAddress string
	tag        string
}

type clusterService interface {
	isActive() (bool, error)
}

type monitoredClusterService struct {
	instance       cluster
	cachedStatus   bool
	cacheTimestamp time.Time
}

func newMonitoredClusterService(dnsAddress string, tag string) monitoredClusterService {
	instance := cluster{dnsAddress: dnsAddress, tag: tag}
	return monitoredClusterService{instance: instance}
}

func (mc monitoredClusterService) isActive() (bool, error) {
	// if the DNS address contains the cluster tag (and implicitly the region)
	// than this means that there is no failover mechanism in place
	if strings.Contains(mc.instance.dnsAddress, mc.instance.tag) {
		return true, nil
	}

	if time.Now().Sub(mc.cacheTimestamp).Seconds() < cacheTime {
		return mc.cachedStatus, nil
	}

	resolver := dnsr.New(5)
	cNames := resolver.Resolve(mc.instance.dnsAddress, "CNAME")
	mc.cacheTimestamp = time.Now()
	if len(cNames) > 0 {
		mc.cachedStatus = strings.Contains(cNames[0].Value, mc.instance.tag)
		return mc.cachedStatus, nil
	}
	mc.cachedStatus = false
	return false, errors.New("address could not be resolved, maybe it is invalid")
}
