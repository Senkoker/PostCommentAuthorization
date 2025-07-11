package handlers

const (
	getUserInfo         = "SELECT user_id,first_name,second_name,img_url FROM users_info WHERE user_id = ANY ($1)"
	createPost          = "INSERT INTO posts (author_id ,tag_ids,content,private) VALUES ($1,$2,$3,$4) returning id;"
	getPosts            = "SELECT id,author_id,tag_ids,content,created_at,watched,likes FROM posts WHERE id = ANY ($1)"
	getPostsWithHashtag = "SELECT id,author_id,tag_ids,content,created_at,watched,likes FROM posts WHERE tag_ids @> $1 ORDER BY watched DESC LIMIT $2 OFFSET $3"
	createMainComment   = "INSERT INTO main_comments (post_id,author_id,content,created_at) VALUES ($1,$2,$3,$4) returning id"
	createComment       = "INSERT INTO comments (comment_id,author_id,content,created_at) VALUES ($1,$2,$3,$4) returning id"
	getMainComment      = "SELECT id,author_id,content,created_at FROM main_comments WHERE post_id = $1 limit $2 offset $3"
	getComment          = "SELECT comment_id,author_id,content,created_at FROM comments WHERE comment_id = $1 limit $2 offset $3"
)
