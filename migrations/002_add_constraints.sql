-- +goose Up

-- +goose StatementBegin
ALTER TABLE posts
    ADD CONSTRAINT fk_user_posts FOREIGN KEY (author_id) REFERENCES users_info(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE main_comments
    ADD CONSTRAINT fk_post_maincomment FOREIGN KEY (post_id) REFERENCES posts(id),
    ADD CONSTRAINT fk_user_maincomment FOREIGN KEY (author_id) REFERENCES users_info(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE comments
    ADD CONSTRAINT fk_main_comment FOREIGN KEY (comment_id) REFERENCES main_comments(id),
    ADD CONSTRAINT fk_user_comment FOREIGN KEY (author_id) REFERENCES users_info(user_id);
-- +goose StatementEnd


-- +goose Down

-- +goose StatementBegin
ALTER TABLE posts
DROP CONSTRAINT fk_user_posts;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE main_comments
    DROP CONSTRAINT fk_post_maincomment,
    DROP CONSTRAINT fk_user_maincomment;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE comments
    DROP CONSTRAINT fk_main_comment,
    DROP CONSTRAINT fk_user_comment;
-- +goose StatementEnd


