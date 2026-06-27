package article

import (
	"time"

	"github.com/core-go/core/convert"
)

type Article struct {
	Id          string     `json:"id" gorm:"primary_key;column:id" bson:"_id" dynamodbav:"id,omitempty" firestore:"-"`
	Slug        string     `json:"slug,omitempty" gorm:"column:slug" bson:"slug,omitempty" dynamodbav:"slug,omitempty" firestore:"slug,omitempty"`
	Title       string     `json:"title,omitempty" gorm:"column:title" bson:"title,omitempty" dynamodbav:"title,omitempty" firestore:"title,omitempty"`
	Description string     `json:"description,omitempty" gorm:"column:description" bson:"description" dynamodbav:"description,omitempty" firestore:"description,omitempty"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" gorm:"column:published_at" bson:"publishedAt,omitempty" dynamodbav:"publishedAt,omitempty" firestore:"publishedAt,omitempty"`
	Content     string     `json:"content,omitempty" gorm:"column:content" bson:"content,omitempty" dynamodbav:"content,omitempty" firestore:"content,omitempty"`
	Thumbnail   string     `json:"thumbnail,omitempty" gorm:"column:thumbnail" bson:"thumbnail,omitempty" dynamodbav:"thumbnail,omitempty" firestore:"thumbnail,omitempty"`
	Tags        []string   `json:"tags,omitempty" gorm:"column:tags" bson:"tags,omitempty" dynamodbav:"tags,omitempty" firestore:"tags,omitempty"`
	Status      string     `json:"status,omitempty" gorm:"column:status" bson:"status" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	AuthorId    string     `json:"authorId,omitempty" gorm:"column:authorid" bson:"authorId,omitempty" dynamodbav:"authorId,omitempty" firestore:"authorId,omitempty"`

	SubmittedBy string     `json:"submittedBy,omitempty" gorm:"column:submitted_by" bson:"submittedBy,omitempty" dynamodbav:"submittedBy,omitempty" firestore:"submittedBy,omitempty"`
	SubmittedAt *time.Time `json:"submittedAt,omitempty" gorm:"column:submitted_at" bson:"submittedAt,omitempty" dynamodbav:"submittedAt,omitempty" firestore:"submittedAt,omitempty"`
	ApprovedBy  string     `json:"approvedBy,omitempty" gorm:"column:approved_by" bson:"approvedBy,omitempty" dynamodbav:"approvedBy,omitempty" firestore:"approvedBy,omitempty"`
	ApprovedAt  *time.Time `json:"approvedAt,omitempty" gorm:"column:approved_at" bson:"approvedAt,omitempty" dynamodbav:"approvedAt,omitempty" firestore:"approvedAt,omitempty"`
	CreatedBy   string     `json:"createdBy,omitempty" gorm:"column:created_by" bson:"createdBy,omitempty" dynamodbav:"createdBy,omitempty" firestore:"createdBy,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty" gorm:"column:created_at" bson:"createdAt,omitempty" dynamodbav:"createdAt,omitempty" firestore:"createdAt,omitempty"`
	UpdatedBy   string     `json:"updatedBy,omitempty" gorm:"column:updated_by" bson:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty" firestore:"updatedBy,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at" bson:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty" firestore:"updatedAt,omitempty"`
	// Name        string     `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
}

var RemovedColumns = []string{"submittedBy", "submittedAt", "approvedBy", "approvedAt", "createdBy", "createdAt", "updatedBy", "updatedAt"}

func (a Article) GetData() map[string]interface{} {
	return convert.ToMap(a)
}
