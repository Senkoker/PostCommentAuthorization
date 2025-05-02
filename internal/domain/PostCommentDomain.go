package domain

import (
	"VK_posts/internal/logger"
	"VK_posts/internal/models"
	"errors"
	"mime/multipart"
)

func subtractSlices(main []string, toRemove []string) []string {
	removeMap := make(map[string]struct{})
	for _, item := range toRemove {
		removeMap[item] = struct{}{}
	}
	result := []string{}
	for _, item := range main {
		if _, exists := removeMap[item]; !exists {
			result = append(result, item)
		}
	}
	return result
}
func uniqueSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0, len(slice))
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
		}
	}
	for entry, _ := range keys {
		list = append(list, entry)
	}
	return list
}

type Domain struct {
	PostgresPostComment PostgresPostAndCommentInterface
	PostgresUserInfo    PostgresUserInformationInterface
	SelectelS3          SelectelS3Interface
	Redis               RedisPostHashInterface
}
type PostgresPostAndCommentInterface interface {
	CreatePost(models.NewPost) (string, error)
	GetPosts(postsIds []string) ([]models.Post, []string, error)
	GetPostWithHashtags(hashtags []string, limit, offset string) ([]models.Post, []string, error)
	CreateMainComment(comment models.NewComment) (string, error)
	CreateComment(comment models.NewComment) (string, error)
	GetMainComment(postID, limit, offset string) ([]models.MainComment, []string, error)
	GetComment(parentID, limit, offset string) ([]models.Comment, []string, error)
}

type PostgresUserInformationInterface interface {
	GetUserInfo(users []string) (map[string]models.UserInfo, error)
}

type RedisPostHashInterface interface {
	GetPostHashtagHash(hashtags []string, limit, offset string) ([]models.Post, error)
	GetPostHash(postIds []string) ([]models.Post, []string, error)
	CreatePopularPostHash(post []models.Post) error
}

type SelectelS3Interface interface {
	SendImage(img *multipart.FileHeader) (string, error)
}

func NewDomain(postgresPostComment PostgresPostAndCommentInterface, postgresUserInfo PostgresUserInformationInterface, selectel SelectelS3Interface, redis RedisPostHashInterface) *Domain {
	return &Domain{PostgresPostComment: postgresPostComment, PostgresUserInfo: postgresUserInfo, SelectelS3: selectel, Redis: redis}
}
func (d *Domain) FeedCreatePost(post models.NewPost) (string, error) {
	op := "Domain CreatePost"
	loggers := logger.GetLogger().With("op", op)
	imgUrl, err := d.SelectelS3.SendImage(post.Img)
	if err != nil {
		loggers.Error("Error sending image to Selectel", "err", err.Error())
		imgUrl = "some_string"
	}
	post.Content = post.Content + "&" + imgUrl
	uuid, err := d.PostgresPostComment.CreatePost(post)
	if err != nil {
		loggers.Error("Problem add Post in Postgres", "err", err.Error())
		return "", err
	}
	return uuid, nil
}
func (d *Domain) FeedCreateComment(comment models.NewComment) (string, error) {
	op := "Domain CreateComment"
	logger := logger.GetLogger().With("op", op)
	if comment.Main {
		uuid, err := d.PostgresPostComment.CreateMainComment(comment)
		if err != nil {
			logger.Error("Problem add Post in Postgres", "err", err.Error())
			return "", err
		}
		return uuid, nil
	}
	uuid, err := d.PostgresPostComment.CreateComment(comment)
	if err != nil {
		logger.Error("Problem too create comment", "err", err.Error())
		return "", err
	}
	return uuid, nil

}

func (d *Domain) FeedGetPosts(interestPostIds []string) ([]models.Post, error) {
	op := "FeedGetPost"
	localLogger := logger.GetLogger().With("op", op)
	posts, postIds, redisErr := d.Redis.GetPostHash(interestPostIds)
	if redisErr != nil {
		localLogger.Error("Problem get post hash", "err", redisErr)

	}
	interestPostIds = subtractSlices(interestPostIds, postIds)
	var popularPost []models.Post
	var userInfo map[string]models.UserInfo
	if interestPostIds != nil {
		postgresPosts, users, err := d.PostgresPostComment.GetPosts(interestPostIds)
		users = uniqueSlice(users)
		if err != nil && redisErr != nil {
			localLogger.Error("Redis error", "err", redisErr)
			localLogger.Error("Postgres error", "err", err)
			return nil, err
		} else if err != nil {
			localLogger.Error("Postgres error", "err", err)
			return posts, nil
		}
		userInfo, err = d.PostgresUserInfo.GetUserInfo(users)
		if err != nil {
			localLogger.Error("Problem get user info postgres", "err", err)
			return posts, err
		}
		for i := 0; i < len(postgresPosts); i++ {
			author := postgresPosts[i].AuthorId
			postgresPosts[i].Author = userInfo[author].FirstName + " " + userInfo[author].SecondName
			postgresPosts[i].ImgPersonURL = userInfo[author].ImgUrl
			if postgresPosts[i].Watched > -1 {
				popularPost = append(popularPost, postgresPosts[i])
			}
			if i == (len(postgresPosts) - 1) {
				err = d.Redis.CreatePopularPostHash(popularPost)
				if err != nil {
					localLogger.Error("Problem to sent post in Redis", "err", err)
				}
			}
		}
		posts = append(posts, postgresPosts...)
	}
	return posts, nil
}
func (d *Domain) FeedGetPostsWithHashtag(hashtags []string, limit, offset, redisStatus string) ([]models.Post, error) {
	op := "Domain GetPosts"
	logger := logger.GetLogger().With("op", op)
	if redisStatus == "true" {
		postWithHashtagsRedis, err := d.Redis.GetPostHashtagHash(hashtags, limit, offset)
		return postWithHashtagsRedis, err
	}
	//Todo: Дописать заполнение автора в пост
	postHashtags, users, err := d.PostgresPostComment.GetPostWithHashtags(hashtags, limit, offset)
	if err != nil {
		if errors.Is(err, errors.New("DoesNotExist")) {
			logger.Error("Problem get Post with hashtags", "err", err.Error())
			return nil, err
		}
		logger.Error("Problem get Post with hashtags", "err", err.Error())
		return nil, err
	}
	usersInfo, err := d.PostgresUserInfo.GetUserInfo(users)
	if err != nil {
		logger.Error("Problem get user info postgres", "err", err.Error())
		return nil, err
	}
	for i := 0; i < len(postHashtags); i++ {
		authorID := postHashtags[i].AuthorId
		info := usersInfo[authorID]
		postHashtags[i].Author = info.FirstName + " " + info.SecondName
		postHashtags[i].ImgPersonURL = info.ImgUrl
	}
	return postHashtags, nil
}

func (d *Domain) FeedGetMainComments(postId string, limit, offset string) ([]models.MainComment, error) {
	op := "Domain GetMainComments"
	localoger := logger.GetLogger().With("op", op)
	mainComments, users, err := d.PostgresPostComment.GetMainComment(postId, limit, offset)
	if err != nil {
		localoger.Error("Errors", "err", err.Error())
		return nil, err
	}
	usersInfo, err := d.PostgresUserInfo.GetUserInfo(users)
	if err != nil {
		localoger.Error("Problem get user info postgres", "err", err.Error())
		return nil, err
	}
	for i := 0; i < len(mainComments); i++ {
		authorID := mainComments[i].AuthorID
		info := usersInfo[authorID]
		mainComments[i].AuthorName = info.FirstName + " " + info.SecondName
		mainComments[i].AuthorImgUrl = info.ImgUrl
	}
	return mainComments, nil
}
func (d *Domain) FeedGetComments(parentID, limit, offset string) ([]models.Comment, error) {
	op := "DomainGetMainComments"
	localoger := logger.GetLogger().With("op", op)
	comments, users, err := d.PostgresPostComment.GetComment(parentID, limit, offset)
	if err != nil {
		localoger.Error("Errors", "err", err.Error())
		return nil, err
	}
	usersInfo, err := d.PostgresUserInfo.GetUserInfo(users)
	for i := 0; i < len(comments); i++ {
		authorID := comments[i].AuthorID
		info := usersInfo[authorID]
		comments[i].AuthorName = info.FirstName + " " + info.SecondName
		comments[i].AuthorImgUrl = info.ImgUrl
	}
	return comments, nil
}
