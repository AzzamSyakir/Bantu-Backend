package cache

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

type JobCache struct {
	Client *configs.RedisConfig
}

func NewJobCache(client *configs.RedisConfig) *JobCache {
	return &JobCache{
		Client: client,
	}
}

func (jobCache *JobCache) GetJobsFromRedis(searchQuery string, provinceID int64, regencyID int64) ([]entity.JobEntity, error) {
	ctx := context.Background()
	var cursor uint64
	var allKeys []string
	var jobs []entity.JobEntity

	for {
		keys, newCursor, err := jobCache.Client.Redis.Connection.Scan(ctx, cursor, "job:*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("error scanning job keys: %v", err)
		}

		allKeys = append(allKeys, keys...)
		cursor = newCursor

		if cursor == 0 {
			break
		}
	}

	if len(allKeys) == 0 {
		return jobs, nil
	}

	jobData, err := jobCache.Client.Redis.Connection.MGet(ctx, allKeys...).Result()
	if err != nil {
		return nil, err
	}

	for _, data := range jobData {
		if data == nil {
			continue
		}

		var job entity.JobEntity
		err := json.Unmarshal([]byte(data.(string)), &job)
		if err != nil {
			return nil, errors.New("error unmarshalling job data")
		}

		if searchQuery != "" && !strings.Contains(strings.ToLower(job.Title), strings.ToLower(searchQuery)) {
			continue
		}
		if provinceID > 0 && job.ProvinceID != provinceID {
			continue
		}
		if regencyID > 0 && job.RegencyID != regencyID {
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (jobCache *JobCache) SaveJobToRedis(job *entity.JobEntity) error {
	ctx := context.Background()
	json, err := json.Marshal(job)
	if err != nil {
		return err
	}
	jobKey := fmt.Sprintf("job:%s", job.ID)
	err = jobCache.Client.Redis.Connection.Set(ctx, jobKey, json, 0).Err()
	if err != nil {
		return err
	}
	log.Printf("Saving job to Redis: %+v\n", jobKey)
	return nil
}

func (jobCache *JobCache) UpdateJobInRedis(job *entity.JobEntity) error {
	ctx := context.Background()
	key := fmt.Sprintf("job:%s", job.ID)

	exists, err := jobCache.Client.Redis.Connection.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return err
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return errors.New("job not found in redis")
	}

	err = jobCache.Client.Redis.Connection.Set(ctx, key, jobJSON, 0).Err()
	if err != nil {
		return err
	}
	fmt.Println("success update job in redis")
	return nil
}

func (jobCache *JobCache) DeleteJobFromRedis(id string) error {
	ctx := context.Background()
	jobKey := fmt.Sprintf("job:%s", id)

	deleted, err := jobCache.Client.Redis.Connection.Del(ctx, jobKey).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return errors.New("job not found in redis")
	}

	fmt.Println("success delete job in redis")
	return nil
}
