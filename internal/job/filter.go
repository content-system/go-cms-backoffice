package job

import "github.com/core-go/search"

type JobFilter struct {
	*search.Filter
	Id             string            `json:"id,omitempty" gorm:"primary_key;column:id" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"max=40" match:"equal"`
	Title          string            `json:"title,omitempty" gorm:"column:title" dynamodbav:"title,omitempty" firestore:"title,omitempty" validate:"omitempty,max=120"`
	Description    string            `json:"description,omitempty" gorm:"column:description" dynamodbav:"description,omitempty" firestore:"description,omitempty" validate:"omitempty,max=1000"`
	PublishedAt    *search.TimeRange `json:"publishedAt,omitempty" gorm:"column:published_at" dynamodbav:"publishedAt,omitempty" firestore:"publishedAt,omitempty"`
	ExpiredAt      *search.TimeRange `json:"expiredAt,omitempty" gorm:"column:expired_at" dynamodbav:"expiredAt,omitempty" firestore:"expiredAt,omitempty"`
	Position       *string           `json:"position,omitempty" gorm:"column:position" dynamodbav:"position,omitempty" firestore:"position,omitempty" validate:"omitempty,max=100"`
	Quantity       *int32            `json:"quantity,omitempty" gorm:"column:quantity" dynamodbav:"quantity,omitempty" firestore:"quantity,omitempty"`
	Location       *string           `json:"location,omitempty" gorm:"column:location" dynamodbav:"location,omitempty" firestore:"location,omitempty"`
	Skills         []string          `json:"skills,omitempty" gorm:"column:skills" dynamodbav:"skills,omitempty" firestore:"skills,omitempty"`
	ApplicantCount *int32            `json:"applicantCount,omitempty" gorm:"column:applicant_count" dynamodbav:"applicantCount,omitempty" firestore:"applicantCount,omitempty"`
	CompanyId      string            `json:"companyId,omitempty" gorm:"column:company_id" dynamodbav:"companyid,omitempty" firestore:"companyid,omitempty"`
	// Status         []string          `json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty" match:"equal" validate:"required,max=1,code"`
}
