package playground

import (
	"fmt"
	"time"
)

func ExampleData() Blog {
	profileJohn := UserProfile{
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Biography: "A passionate writer about Go.",
		Website:   "https://johndoe.com",
	}

	profileJane := UserProfile{
		Birthday:  time.Date(1985, 2, 15, 0, 0, 0, 0, time.UTC),
		Biography: "A web developer and writer.",
		Website:   "https://jane.com",
	}

	authorJohn := Author{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Profile:   profileJohn,
	}

	authorJane := Author{
		ID:        2,
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		Profile:   profileJane,
	}

	posts := []Post{
		{
			ID:          1,
			Title:       "Understanding Golang",
			Content:     "Go is an open-source programming language...",
			PublishDate: time.Now().AddDate(0, -1, 0),
			Author:      authorJohn,
			Comments:    generateComments(3),
			Tags:        []string{"Go", "Programming"},
			Likes:       100,
		},
		{
			ID:          2,
			Title:       "Exploring Python",
			Content:     "Python is versatile...",
			PublishDate: time.Now().AddDate(0, -2, 0),
			Author:      authorJane,
			Comments:    generateComments(4),
			Tags:        []string{"Python", "Development"},
			Likes:       150,
		},
		{
			ID:          3,
			Title:       "Web Development Basics",
			Content:     "The world of web development...",
			PublishDate: time.Now().AddDate(0, -3, 0),
			Author:      authorJane,
			Comments:    generateComments(5),
			Tags:        []string{"Web", "HTML", "CSS"},
			Likes:       125,
		},
		{
			ID:          4,
			Title:       "Machine Learning in a Nutshell",
			Content:     "ML is revolutionizing industries...",
			PublishDate: time.Now().AddDate(0, -5, 0),
			Author:      authorJohn,
			Comments:    generateComments(6),
			Tags:        []string{"ML", "AI"},
			Likes:       200,
		},
		{
			ID:          5,
			Title:       "JavaScript: The Good Parts",
			Content:     "JavaScript powers the web...",
			PublishDate: time.Now().AddDate(0, -4, 0),
			Author:      authorJane,
			Comments:    generateComments(3),
			Tags:        []string{"JavaScript", "Web"},
			Likes:       170,
		},
	}

	blog := Blog{
		Posts:      make([]Post, len(posts)),
		Authors:    map[int]Author{authorJohn.ID: authorJohn, authorJane.ID: authorJane},
		TotalViews: 10000,
		TotalPosts: len(posts),
		TotalLikes: 0,
	}

	for i, post := range posts {
		blog.Posts[i] = post
		blog.TotalLikes += post.Likes
	}

	return blog
}

func generateComments(count int) []Comment {
	var comments []Comment
	for i := 1; i <= count; i++ {
		comment := Comment{
			ID:          i,
			AuthorName:  fmt.Sprintf("Commenter %d", i),
			Content:     fmt.Sprintf("This is comment %d!", i),
			CommentDate: time.Now().AddDate(0, 0, -i),
			Upvotes:     i * 5,
		}
		comments = append(comments, comment)
	}
	return comments
}
