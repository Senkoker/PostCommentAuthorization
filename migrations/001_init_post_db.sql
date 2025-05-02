-- +goose Up

-- +goose StatementBegin
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id UUID NOT NULL ,
    tag_ids TEXT[],
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    watched INTEGER NOT NULL DEFAULT 0,
    likes INTEGER NOT NULL DEFAULT 0,
    private BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX post_id_tags_search ON posts(id,tag_ids)
-- +goose StatementEnd


-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS main_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    author_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX post_id_created_at_search ON main_comments(post_id,created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    comment_id UUID NOT NULL,
    author_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX comment_reply_created_search ON comments(comment_id,created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS selection_for_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() ,
    userid UUID NOT NULL,
    new_posts UUID[] ,
    viewed_posts UUID[],
    friends UUID[]
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender UUID NOT NULL,
    recipient UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS server_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    userid UUID NOT NULL,
    server int NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_info (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    first_name VARCHAR(64) NOT NULL,
    second_name VARCHAR(64) NOT NULL,
    images TEXT[],
    img_url TEXT NOT NULL,
    birth_date DATE,
    education VARCHAR(64),
    country VARCHAR(128),
    city VARCHAR,
    postIDs TEXT[]
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX user_search ON user_info(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE user_friends (
    user_id UUID,
    friend_id UUID,
    status BOOLEAN,
    PRIMARY KEY (user_id, friend_id)
);
-- +goose StatementEnd
-- +goose Down

-- +goose StatementBegin
DROP TABLE users_info,server_log,messages,selection_for_user,comments,main_comments,posts,user_friends
-- +goose StatementEnd
