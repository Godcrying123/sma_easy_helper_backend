package models

import (
	"errors"

	"github.com/bitly/go-simplejson"
	"github.com/mitchellh/mapstructure"

	"github.com/astaxie/beego"
)

// Cluster model is structing the cluster model with below attribute
type Cluster struct {
	ClusterID     int
	Name          string
	NumOfMachines int
	ClusterLabel  string
	Machines      []Machine
}

// ClusterList is the collection of Cluster
var (
	ClusterMap       map[string]Cluster
	ClusterList      []Cluster
	ClusterID        *int
	ClusterFilesChan chan string
	ClusterChan      chan Cluster
)

func init() {
	ClusterMap = make(map[string]Cluster)
	ClusterList = make([]Cluster, 0)
	ClusterFilesChan = make(chan string)
	ClusterChan = make(chan Cluster)
	var WholeClusterFileString string
	go func() {
		FileListDir("./saved_infos/cluster/", ClusterFilesChan)
		close(ClusterFilesChan)
	}()
	for FileName := range ClusterFilesChan {
		fileString, err := FileRead(FileName)
		if err != nil {
			beego.Error(err)
		} else {
			WholeClusterFileString += fileString
		}
	}
	go func() {
		ClusterJSONRead([]byte(WholeClusterFileString), ClusterChan)
		close(ClusterChan)
	}()
	for clusterEntity := range ClusterChan {
		ClusterMap[clusterEntity.Name] = clusterEntity
		ClusterList = append(ClusterList, clusterEntity)
	}
}

// ClusterJSONRead function is for transfering JSON file to structure
func ClusterJSONRead(fileByte []byte, ClusterChan chan<- Cluster) (chan<- Cluster, error) {
	jsons, _ := simplejson.NewJson(fileByte)
	for _, jsonEntity := range jsons.MustArray() {
		cluster := Cluster{}
		err := mapstructure.WeakDecode(jsonEntity.(map[string]interface{}), &cluster)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		ClusterChan <- cluster
	}
	return ClusterChan, nil
}

// GetCluster function is for getting the cluster by its name
func GetCluster(cname string) (c Cluster, err error) {
	if c, ok := ClusterMap[cname]; ok {
		return c, nil
	}
	return c, errors.New("Cluster not exist")
}

// GetAllClusters function is for getting all clusters
func GetAllClusters() []Cluster {
	return ClusterList
}

// DeleteCluster function is for deleting the cluster by its name
func DeleteCluster(cname string) {
	delete(ClusterMap, cname)
}

// AddCluster function is for adding the cluster entity and return its ID
func AddCluster(c Cluster) int {
	c.ClusterID = *ClusterID
	*ClusterID++
	ClusterMap[c.Name] = c
	return c.ClusterID
}
