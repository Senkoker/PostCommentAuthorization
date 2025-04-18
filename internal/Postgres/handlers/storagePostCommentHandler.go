package handlers

import (
	"VK_posts/internal/logger"
	"VK_posts/internal/models"
	"VK_posts/pkg/Postgres"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StorageHandler struct {
	storage *Postgres.Storage
	ctxTime time.Duration
}

func NewStorageHandler(storage *Postgres.Storage, ctxTime time.Duration) *StorageHandler {
	return &StorageHandler{storage: storage, ctxTime: ctxTime}
}
func (s *StorageHandler) GetUserInfo(users []string) (map[string]models.UserInfo, error) {
	stmt, err := s.storage.Db.Prepare(getUserInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare GetUserInfo: %w", err)
	}
	ctx := context.Background()
	rows, err := stmt.QueryContext(ctx, users)
	if err != nil {
		return nil, fmt.Errorf("Poblem to return data GetUserInfo: %w", err)
	}
	usersInfo := make(map[string]models.UserInfo)
	var user models.UserInfo
	var userId string
	for rows.Next() {
		err = rows.Scan(&userId, &user.FirstName, &user.SecondName, &user.ImgUrl)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.GetLogger().Error("User does not exist")
				continue
			}
			return nil, fmt.Errorf("failed to scan rows GetUserInfo: %w", err)
		}
		usersInfo[userId] = user

	}
	return usersInfo, nil
}

func (s *StorageHandler) CreatePost(post models.NewPost) (string, error) {
	stmt, err := s.storage.Db.Prepare(createPost)
	if err != nil {
		return "", fmt.Errorf("failed to prepare CreatePost: %w", err)
	}
	var id string
	ctx := context.Background()
	err = stmt.QueryRowContext(ctx, post.AuthorId, post.Hashtags, post.Content, post.Private).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to scan rows CreatePost: %w", err)
	}
	return id, nil
}
func (s *StorageHandler) GetPosts(postIDS []string) ([]models.Post, []string, error) {
	stmt, err := s.storage.Db.Prepare(getPosts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to prepare GetPosts: %w", err)
	}
	ctx := context.Background()
	row, err := stmt.QueryContext(ctx, postIDS)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query GetPosts: %w", err)
	}
	var post models.Post
	var posts []models.Post
	var users []string

	for row.Next() {
		err = row.Scan(&post.PostID, &post.AuthorId, &post.TagsIds, &post.Content, &post.CreatedAt, &post.Watched, &post.Likes)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.GetLogger().Error("Post does not exist")
			}
			continue
		}
		posts = append(posts, post)
		users = append(users, post.AuthorId)
	}
	return posts, users, nil

}

func (s *StorageHandler) GetPostWithHashtags(hashtags []string, limit, offset string) ([]models.Post, []string, error) {
	stmt, err := s.storage.Db.Prepare(getPostsWithHashtag)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to prepare GetPosts: %w", err)
	}
	ctx := context.Background()
	row, err := stmt.QueryContext(ctx, hashtags, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query GetPostsWithHashtags: %w", err)
	}
	var post models.Post
	var posts []models.Post
	var users []string

	for row.Next() {
		err = row.Scan(&post.PostID, &post.AuthorId, &post.TagsIds, &post.Content, &post.CreatedAt, &post.Watched, &post.Likes)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.GetLogger().Error("Post with hashtag does not exist")
			}
			continue
		}
		posts = append(posts, post)
		users = append(users, post.AuthorId)
	}
	return posts, users, nil
}

func (s *StorageHandler) CreateMainComment(comment models.NewComment) (string, error) {
	stmt, err := s.storage.Db.Prepare(createMainComment)
	if err != nil {
		return "", fmt.Errorf("failed to prepare CreateMainComment: %w", err)
	}
	var id string
	ctx := context.Background()
	err = stmt.QueryRowContext(ctx, comment.ReplyTo, comment.AuthorID, comment.Content, comment.CreatedAt).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to scan rows CreateMainComment: %w", err)
	}
	return id, nil
}

func (s *StorageHandler) CreateComment(comment models.NewComment) (string, error) {
	stmt, err := s.storage.Db.Prepare(createComment)
	if err != nil {
		return "", fmt.Errorf("failed to prepare CreateComment: %w", err)
	}
	var id string
	ctx := context.Background()
	err = stmt.QueryRowContext(ctx, comment.ReplyTo, comment.AuthorID, comment.Content, comment.CreatedAt).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to scan rows CreateComment: %w", err)
	}
	return id, nil
}

func (s *StorageHandler) GetMainComment(postID, limit, offset string) ([]models.MainComment, []string, error) {
	stmt, err := s.storage.Db.Prepare(getMainComment)
	if err != nil {
		return nil, nil, err
	}
	var mainComment models.MainComment
	var mainComments []models.MainComment
	var users []string
	ctx := context.Background()
	rows, err := stmt.QueryContext(ctx, postID, limit, offset)
	for rows.Next() {
		err = rows.Scan(&mainComment.PostID, &mainComment.AuthorID, &mainComment.Content, &mainComment.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.GetLogger().Error("Main comment does not exist")
			}
			continue
		}
		users = append(users, mainComment.AuthorID)
		mainComments = append(mainComments, mainComment)
	}
	return mainComments, users, nil

}

func (s *StorageHandler) GetComment(parentId, limit, offset string) ([]models.Comment, []string, error) {
	stmt, err := s.storage.Db.Prepare(getComment)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to prepare GetComment: %w", err)
	}
	var mainComment models.Comment
	var mainComments []models.Comment
	var users []string
	ctx := context.Background()
	rows, err := stmt.QueryContext(ctx, parentId, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query GetComment: %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&mainComment.MainCommentID, &mainComment.AuthorID, &mainComment.Content, &mainComment.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.GetLogger().Error("Comment does not exist")
			}
			continue
		}
		users = append(users, mainComment.AuthorID)
		mainComments = append(mainComments, mainComment)
	}
	return mainComments, users, nil

}
