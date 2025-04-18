package handlers

const (
	getUserInfo         = "SELECT user_id,first_name,second_name,img_url FROM users_info WHERE id = ANY ($1)"
	createPost          = "INSERT INTO posts (author_id ,tag_ids,content,private) VALUES ($1,$2,$3,$4) returning id;"
	getPosts            = "SELECT id,author_id,tag_ids,content,crated_at,watched,likes FROM posts WHERE id = ANY ($1)"
	getPostsWithHashtag = "SELECT id,author_id,tag_ids,content,crated_at,watched,likes FROM posts ORDER BY watched WHERE tag_ids @> $1  DESC LIMIT $2 OFFSET $3"
	createMainComment   = "INSERT INTO main_comments (post_id,author_id,content,created_at) VALUES ($1,$2,$3,$4) returning id"
	createComment       = "INSERT INTO comments (comment_id,author_id,content,created_at) VALUES ($1,$2,$3,$4) returning id"
	getMainComment      = "SELECT post_id,author_id,content,created_at FROM main_comments WHERE id = $1 limit $2 offset $3"
	getComment          = "SELECT comment_id,author_id,content,created_at FROM comments WHERE id = $1 limit $2 offset $3"
)
