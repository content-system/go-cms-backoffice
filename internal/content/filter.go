package content

import "github.com/core-go/search"

type ContentFilter struct {
	*search.Filter
	Id          string            `json:"id" gorm:"primary_key;column:id" bson:"_id" dynamodbav:"id,omitempty" firestore:"-" match:"equal"`
	Lang        string            `json:"lang,omitempty" gorm:"primary_key;column:lang" bson:"lang,omitempty" dynamodbav:"lang,omitempty" firestore:"lang,omitempty" match:"equal"`
	Title       string            `json:"title,omitempty" gorm:"column:title" bson:"title,omitempty" dynamodbav:"title,omitempty" firestore:"title,omitempty"`
	Body        string            `json:"body,omitempty" gorm:"column:body" bson:"body,omitempty" dynamodbav:"body,omitempty" firestore:"body,omitempty"`
	PublishedAt *search.TimeRange `json:"publishedAt,omitempty" gorm:"column:published_at" bson:"publishedAt,omitempty" dynamodbav:"publishedAt,omitempty" firestore:"publishedAt,omitempty"`
	Tags        []string          `json:"tags,omitempty" gorm:"column:tags" bson:"tags,omitempty" dynamodbav:"tags,omitempty" firestore:"tags,omitempty"`
	Status      []string          `json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty" match:"equal" validate:"required,max=1,code"`
}
