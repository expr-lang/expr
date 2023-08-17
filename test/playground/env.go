package playground

import (
	"time"
)

type UserProfile struct {
	Birthday  time.Time
	Biography string
	Website   string
}

func (u UserProfile) Age() int {
	return time.Now().Year() - u.Birthday.Year()
}

type Author struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Profile   UserProfile
}

func (a Author) FullName() string {
	return a.FirstName + " " + a.LastName
}

func (a Author) IsAdmin() bool {
	return a.ID == 1
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

func (p Post) Published() bool {
	return !p.PublishDate.IsZero()
}

func (p Post) After(date time.Time) bool {
	return p.PublishDate.After(date)
}

func (p Post) Before(date time.Time) bool {
	return p.PublishDate.Before(date)
}

func (p Post) Compare(other Post) int {
	if p.PublishDate.Before(other.PublishDate) {
		return -1
	} else if p.PublishDate.After(other.PublishDate) {
		return 1
	}
	return 0
}

func (p Post) Equal(other Post) bool {
	return p.Compare(other) == 0
}

func (p Post) IsZero() bool {
	return p.ID == 0 && p.Title == "" && p.Content == "" && p.PublishDate.IsZero()
}

type Comment struct {
	ID          int
	AuthorName  string
	Content     string
	CommentDate time.Time
	Upvotes     int
}

func (c Comment) Upvoted() bool {
	return c.Upvotes > 0
}

func (c Comment) AuthorEmail() string {
	return c.AuthorName + "@example.com"
}

type Blog struct {
	Posts      []Post
	Authors    map[int]Author
	TotalViews int
	TotalPosts int
	TotalLikes int
}

func (b Blog) RecentPosts() []Post {
	var posts []Post
	for _, post := range b.Posts {
		if post.Published() {
			posts = append(posts, post)
		}
	}
	return posts
}

func (b Blog) PopularPosts() []Post {
	var posts []Post
	for _, post := range b.Posts {
		if post.Likes > 150 {
			posts = append(posts, post)
		}
	}
	return posts
}

func (b Blog) TotalUpvotes() int {
	var upvotes int
	for _, post := range b.Posts {
		for _, comment := range post.Comments {
			upvotes += comment.Upvotes
		}
	}
	return upvotes
}

func (b Blog) TotalComments() int {
	var comments int
	for _, post := range b.Posts {
		comments += len(post.Comments)
	}
	return comments
}

func (Blog) Add(a, b float64) float64 {
	return a + b
}

func (Blog) Sub(a, b float64) float64 {
	return a - b
}

func (Blog) Title(post Post) string {
	return post.Title
}

func (Blog) HasTag(post Post, tag string) bool {
	for _, t := range post.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (Blog) IsAdmin(author Author) bool {
	return author.IsAdmin()
}

func (Blog) IsZero(post Post) bool {
	return post.IsZero()
}

func (Blog) WithID(posts []Post, id int) Post {
	for _, post := range posts {
		if post.ID == id {
			return post
		}
	}
	return Post{}
}
