package admin

// NewsPostInput is the body for creating/updating a Home feed post.
type NewsPostInput struct {
	Author   string `json:"author" binding:"required"`
	Avatar   string `json:"avatar"`
	Title    string `json:"title" binding:"required"`
	Body     string `json:"body" binding:"required"`
	Image    string `json:"image"`
	Category string `json:"category"`
	Likes    int    `json:"likes"`
}

// UpdateProgressInput is the body for editing a user's progress on a level.
type UpdateProgressInput struct {
	Stars     int  `json:"stars"`
	Completed bool `json:"completed"`
}
