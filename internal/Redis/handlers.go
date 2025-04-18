package Redis

import (
	"VK_posts/internal/logger"
	"VK_posts/internal/models"
	"VK_posts/pkg/Redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type HandlerRedis struct {
	storage *Redis.RedisStorage
	ctxTime time.Duration
}

func NewRedisHandler(storage *Redis.RedisStorage, ctxTime time.Duration) *HandlerRedis {
	return &HandlerRedis{storage: storage, ctxTime: ctxTime}
}

func (r *HandlerRedis) GetPostHashtagHash(hashtags []string, limit, offset string) ([]models.Post, error) {
	var allQuery string
	for i, hashtag := range hashtags {
		text := fmt.Sprintf("@tags_ids:{%s}", hashtag)
		if i == 0 {
			allQuery = text
		} else {
			allQuery = allQuery + " " + text
		}
	}
	ctx := context.Background()
	res, err := r.storage.RStorage.Do(ctx,
		"FT.SEARCH",
		"Post_index",
		allQuery,
		"LIMIT", offset, limit,
	).Result()
	if err != nil {
		return nil, fmt.Errorf("Problem to return data from redis:%w", err)
	}
	var posts []models.Post
	var post models.Post
	result := res.([]interface{})
	for i := 1; i < len(result); i += 2 {
		fields := result[i+1].([]interface{})
		postRedis := fields[1].(string)
		err = json.Unmarshal([]byte(postRedis), &post)
		if err != nil {
			logger.GetLogger().Info("Problem to convert data to posthashtag struct", "err", err)
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}
func (r *HandlerRedis) GetPostHash(postIds []string) ([]models.Post, []string, error) {
	pipe := r.storage.RStorage.Pipeline()
	cmds := make([]*redis.StringCmd, len(postIds))
	ctx := context.Background()
	for i, postId := range postIds {
		cmds[i] = pipe.Get(ctx, fmt.Sprintf("Post_id:%s", postId))
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("problem to connect any statement for send query to redis:%w", err)
	}
	posts := make([]models.Post, len(postIds))
	var existsPostsIds []string
	var data []byte
	for i, cmd := range cmds {
		data, err = cmd.Bytes()
		if err == redis.Nil {
			logger.GetLogger().Error("Data is empty", "err", err)
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		var post models.Post
		err = json.Unmarshal(data, &post)
		if err != nil {
			return nil, nil, err
		}
		existsPostsIds = append(existsPostsIds, post.PostID)
		posts[i] = post
	}
	return posts, existsPostsIds, nil
}
func (r *HandlerRedis) CreatePopularPostHash(posts []models.Post) error {
	pipe := r.storage.RStorage.Pipeline()
	for _, post := range posts {
		jsonPost, err := json.Marshal(post)
		if err != nil {
			return fmt.Errorf("problem to marshal post to redis:%w", err)
		}
		ctx := context.Background()
		pipe.Set(ctx, fmt.Sprintf("Post_id:%s", post.PostID), jsonPost, 0)
	}
	ctx := context.Background()
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("problem to connect any statement for send query to redis:%w", err)
	}
	return nil
}
