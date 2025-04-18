CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id UUID NOT NULL ,
    tag_ids TEXT[],
    content TEXT NOT NULL,
    crated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    watched INTEGER NOT NULL DEFAULT 0,
    likes INTEGER NOT NULL DEFAULT 0,
    private BOOLEAN NOT NULL,
    status BOOLEAN NOT NULL DEFAULT True
    CONSTRAINT fk_user_posts FOREIGN KEY (author_id) REFERENCES users(id)
);
CREATE TABLE main_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    author_id UUID NOT NULL,
    content TEXT NOT NULL,
    crated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_post_maincomment FOREIGN KEY (post_id) REFERENCES posts(id),
    CONSTRAINT fk_user_maincomment FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    comment_id UUID NOT NULL,
    author_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
    CONSTRAINT fk_main_comment FOREIGN KEY (comment_id) REFERENCES main_comments(post_id)
    CONSTRAINT fk_user_comment FOREIGN KEY (author_id) REFERENCES users(id)
);
CREATE TABLE selection_for_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() ,
    userid UUID,
    new_posts UUID[],
    viewed_posts UUID[],
    friends UUID[]
);