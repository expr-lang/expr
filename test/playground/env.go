package playground

import (
	"time"
)

type UserProfile struct {
	Birthday  time.Time
	Biography string
	Website   string
}

type Author struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Profile   UserProfile
}

type Post struct {
	ID          int
	Title       string
	Content     string
	PublishDate time.Time
	Author      Author
	Comments    []Comment
	Tags        []string
	Likes       int
}

type Comment struct {
	ID          int
	AuthorName  string
	Content     string
	CommentDate time.Time
	Upvotes     int
}

type Blog struct {
	Posts      []Post
	Authors    map[int]Author
	TotalViews int
	TotalPosts int
	TotalLikes int
}
