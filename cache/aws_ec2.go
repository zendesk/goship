package cache

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/providers"
	"github.com/zendesk/goship/resources"
)

// AwsEc2Cache implements GoshipCache interface for EC2
type AwsEc2Cache struct {
	Path             string
	CacheTimeSeconds uint64
	Provider         providers.AwsEc2Provider
	Result           []*resources.Ec2Instance
}

// NewAwsEc2Cache returns new AwsEc2Cache
func NewAwsEc2Cache(provider providers.AwsEc2Provider) *AwsEc2Cache {
	return &AwsEc2Cache{
		Provider:         provider,
		CacheTimeSeconds: config.GlobalConfig.CacheValidity,
		Result:           []*resources.Ec2Instance{},
	}
}

func (g *AwsEc2Cache) cacheFilePath() string {
	return fmt.Sprintf("%s%s_%s_%s", path.Join(config.GlobalConfig.CacheDirectory, config.GlobalConfig.CacheFilePrefix), g.Provider.Name(), g.Provider.AwsProfileName, g.Provider.AwsRegion)
}

func (g *AwsEc2Cache) cacheOutdated() bool {
	stat, err := os.Stat(g.cacheFilePath())
	if os.IsNotExist(err) || time.Since(stat.ModTime()).Seconds() > float64(g.CacheTimeSeconds) {
		return true
	}
	return false
}

// CacheName returns cache name
func (g *AwsEc2Cache) CacheName() string {
	return fmt.Sprintf("%s.%s", g.Provider.AwsProfileName, g.Provider.AwsRegion)
}

// Resources returns all resources stored in cache
func (g *AwsEc2Cache) Resources() []resources.Resource {
	var asInt []resources.Resource
	for _, v := range g.Result {
		asInt = append(asInt, v)
	}
	return asInt
}

// Len returns count of resources stored in cache
func (g *AwsEc2Cache) Len() int {
	return len(g.Result)
}

// Refresh refreshes resources
func (g *AwsEc2Cache) Refresh(force bool) (refreshed bool, err error) {
	if g.cacheOutdated() == false && force == false {
		g.Read()
		return false, nil
	}

	g.Provider.Init()

	result, err := g.Provider.GetResources()
	if err != nil {
		return false, err
	}
	for _, r := range result {
		g.Result = append(g.Result, r.(*resources.Ec2Instance))
	}
	err = g.Save()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Save saves the cache to cache file
func (g *AwsEc2Cache) Save() error {
	cacheFile, _ := os.Create(g.cacheFilePath())
	defer cacheFile.Close()
	encoder := gob.NewEncoder(cacheFile)
	err := encoder.Encode(&g.Result)
	if err != nil {
		fmt.Printf("Error while attempting to save cache file %s: %s\n", g.cacheFilePath(), err.Error())
	}
	return nil
}

// Read reads the cache from cache file
func (g *AwsEc2Cache) Read() error {
	file, err := os.Open(g.cacheFilePath())
	if err == nil {
		decoder := gob.NewDecoder(file)
		//var r ec2.DescribeInstancesOutput
		err = decoder.Decode(&g.Result)
		if err != nil {
			fmt.Printf("Error while attempting to read cache file %s: %s\n", g.cacheFilePath(), err.Error())
		}
		//fmt.Printf("%v", g.Result)
	}
	file.Close()
	return err
}
