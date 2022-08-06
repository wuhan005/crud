# 🤖 crud ![Go](https://github.com/wuhan005/crud/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/crud)](https://goreportcard.com/report/github.com/wuhan005/crud) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/wuhan005/crud)

## Demo

1. There is a table `posts` in the database, whose structure is as follows:

```sql
CREATE SEQUENCE IF NOT EXISTS post_id_seq;

-- Table Definition
CREATE TABLE "public"."posts" (
    "id"         int4 NOT NULL DEFAULT nextval('post_id_seq'::regclass),
    "uid"        bpchar(1),
    "title"      bpchar(1),
    "content"    bpchar(1),
    "created_at" time,
    "updated_at" time,
    "deleted_at" time,
    PRIMARY KEY ( "id" )
);
```

2. Run the following command to auto generate code based on the table structure:

```bash
crud gen --dsn=postgres://postgres:postgres@localhost:5432/crud?sslmode=disable
```

3. Here is what you get, amazing!

```go
package db

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ PostsStore = (*posts)(nil)

// Posts is the default instance of the PostsStore.
var Posts PostsStore

// PostsStore is the persistent interface for posts.
type PostsStore interface {
	// GetByID returns a post with the given id.
	// The zero value in the options will be ignored.
	GetByID(ctx context.Context, id int64) (*Post, error)
}

// NewPostsStore returns a PostsStore instance with the given database connection.
func NewPostsStore(db *gorm.DB) PostsStore {
	return &posts{db}
}

// Post represents the posts.
type Post struct {
	gorm.Model

	UID     string
	Title   string
	Content string
}

type posts struct {
	*gorm.DB
}

var (
	ErrPostNotExists = errors.New("post dose not exist")
)

func (db *posts) GetByID(ctx context.Context, id int64) (*Post, error) {
	var post Post
	if err := db.WithContext(ctx).Model(&Post{}).Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotExists
		}
	}
	return &post, nil
}
```

## License

MIT License
