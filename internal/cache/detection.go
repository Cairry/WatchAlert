package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
)

type (
	detectionCache struct {
		rc *redis.Client
	}

	InterDetectionCache interface {
	}
)

func newDetectionCacheInterface(r *redis.Client) InterDetectionCache {
	return detectionCache{
		rc: r,
	}
}

func (d detectionCache) Set(reqId string, data models.DetectionData) error {
	dataList, err := d.Get(reqId)
	if err != nil {
		return err
	}

	dataList = append(dataList, data)

	d.rc.Set(reqId, cmd.JsonMarshal(dataList), 0)

	return nil
}

func (d detectionCache) Get(reqId string) ([]models.DetectionData, error) {
	var cacheDataList []models.DetectionData
	result, err := d.rc.Get(reqId).Result()
	if err != nil {
		return cacheDataList, err
	}

	err = json.Unmarshal([]byte(result), &cacheDataList)
	if err != nil {
		return cacheDataList, err
	}

	return cacheDataList, nil
}

func (d detectionCache) Delete() {

}
