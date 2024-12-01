package article

import "time"

type Article struct {
	Id          string     `json:"id" gorm:"primary_key;column:id" bson:"_id" dynamodbav:"id,omitempty" firestore:"-"`
	Title       string     `json:"title,omitempty" gorm:"column:title" bson:"title,omitempty" dynamodbav:"title,omitempty" firestore:"title,omitempty"`
	Description string     `json:"description,omitempty" gorm:"column:description" bson:"description" dynamodbav:"description,omitempty" firestore:"description,omitempty"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" gorm:"column:published_at" bson:"publishedAt,omitempty" dynamodbav:"publishedAt,omitempty" firestore:"publishedAt,omitempty"`
	Content     string     `json:"content,omitempty" gorm:"column:content" bson:"content,omitempty" dynamodbav:"content,omitempty" firestore:"content,omitempty"`
	Thumbnail   string     `json:"thumbnail,omitempty" gorm:"column:thumbnail" bson:"thumbnail,omitempty" dynamodbav:"thumbnail,omitempty" firestore:"thumbnail,omitempty"`
	Tags        []string   `json:"tags,omitempty" gorm:"column:tags" bson:"tags,omitempty" dynamodbav:"tags,omitempty" firestore:"tags,omitempty"`
	// Type        string     `json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty" validate:"required"`
	Status *string `json:"status,omitempty" gorm:"column:status" bson:"status" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	// AuthorId    string     `json:"authorId,omitempty" gorm:"column:authorid" bson:"authorId,omitempty" dynamodbav:"authorId,omitempty" firestore:"authorId,omitempty"`
	// Name        string     `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
}
