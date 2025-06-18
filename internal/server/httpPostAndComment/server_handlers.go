package httpPostAndComment

import (
	"VK_posts/internal/models"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

type PortInterface interface {
	FeedCreatePost(post models.NewPost) (string, error)
	FeedCreateComment(comment models.NewComment) (string, error)
	FeedGetPosts(interestingPostIDS []string) ([]models.Post, error)
	FeedGetPostsWithHashtag(hashtags []string, limit, offset, redisStatus string) ([]models.Post, error)
	FeedGetMainComments(postId, limit, offset string) ([]models.MainComment, error)
	FeedGetComments(parentID, limit, offset string) ([]models.Comment, error)
}
type Port struct {
	Domain PortInterface
}

func NewHandlers(domain PortInterface) *Port {
	return &Port{domain}
}

func (p *Port) CreatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		authorId := c.Get("userID").(string)
		hashtagsStr := c.Request().FormValue("hashtags")
		hashtags := strings.Split(hashtagsStr, "#")
		hashtags = hashtags[1:]
		fmt.Println("HASTAGS", hashtags)
		content := c.Request().FormValue("content")
		_, img, err := c.Request().FormFile("img")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		var private bool
		strPrivate := c.Request().FormValue("private")
		if strPrivate == "true" {
			private = true
		} else {
			private = false
		}
		newPost := models.NewPost{
			AuthorId: authorId,
			Hashtags: hashtags,
			Content:  content,
			Img:      img,
			Private:  private,
		}
		postUUID, err := p.Domain.FeedCreatePost(newPost)
		if err != nil {
			//Todo: обработать различные типы ошибок
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, models.CreatePostResponse{Id: postUUID})
	}
}

func (p *Port) CreateComment() echo.HandlerFunc {
	return func(c echo.Context) error {
		authorID := c.Get("userID").(string)
		queryMain := c.QueryParam("main")
		var main bool
		if queryMain == "true" {
			main = true
		} else {
			main = false
		}
		var comment models.NewComment
		err := c.Bind(&comment)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		comment.Main = main
		comment.AuthorID = authorID
		UUID, err := p.Domain.FeedCreateComment(comment)
		if err != nil {
			//Todo: обработать ошибки
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, models.CreateCommentResponse{Id: UUID})

	}
}

func (p *Port) GetPosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		//Todo: сравнить совпадают ли Keys
		var usersPosts models.InterestingPost
		err := c.Bind(&usersPosts)
		fmt.Println(usersPosts)
		if err != nil {
			fmt.Println("Ошибка здесь")
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		hashtagsStr := usersPosts.Hashtags
		var hashtags []string
		if hashtagsStr != "" {
			hashtags = strings.Split(hashtagsStr, "#")
			fmt.Println("Hashtags test")
			hashtags = hashtags[1:]
		}
		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")
		redisStatus := c.QueryParam("redis")
		fmt.Println(hashtags, limit, offset, redisStatus)
		fmt.Println(len(hashtags))
		if hashtags == nil {
			posts, err := p.Domain.FeedGetPosts(usersPosts.PostsIDs)
			if err != nil {
				//Todo: обработать ошибки
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, models.GetPostResponse{Posts: posts})
		}
		fmt.Println("C хэштэгами")
		posts, err := p.Domain.FeedGetPostsWithHashtag(hashtags, limit, offset, redisStatus)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, models.GetPostResponse{Posts: posts})
	}
}

func (p *Port) GetMainComment() echo.HandlerFunc {
	return func(c echo.Context) error {
		postID := c.QueryParam("post_id")
		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")
		mainComments, err := p.Domain.FeedGetMainComments(postID, limit, offset)
		if err != nil {
			// Todo: обработать ошибки
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, models.GetMainCommentResponse{MainComments: mainComments})
	}
}

func (p *Port) GetComment() echo.HandlerFunc {
	return func(c echo.Context) error {
		postID := c.QueryParam("reply_id")
		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")
		comments, err := p.Domain.FeedGetComments(postID, limit, offset)
		if err != nil {
			// Todo: обработать ошибки
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, models.GetCommentResponse{Comments: comments})
	}
}
