package job

import "github.com/core-go/search"

type JobFilter struct {
	*search.Filter
	Id             string            `json:"id,omitempty" gorm:"primary_key;column:id" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"max=40" match:"equal"`
	Title          string            `json:"title,omitempty" gorm:"column:title" dynamodbav:"title,omitempty" firestore:"title,omitempty" validate:"omitempty,max=120"`
	Description    string            `json:"description,omitempty" gorm:"column:description" dynamodbav:"description,omitempty" firestore:"description,omitempty" validate:"omitempty,max=1000"`
	Requirements   string            `json:"requirements,omitempty" gorm:"column:requirements" dynamodbav:"requirements,omitempty" firestore:"requirements,omitempty"`
	Benefit        string            `json:"benefit,omitempty" gorm:"column:benefit" dynamodbav:"benefit,omitempty" firestore:"benefit,omitempty"`
	PublishedAt    *search.TimeRange `json:"publishedAt,omitempty" gorm:"column:published_at" dynamodbav:"publishedAt,omitempty" firestore:"publishedAt,omitempty"`
	ExpiredAt      *search.TimeRange `json:"expiredAt,omitempty" gorm:"column:expired_at" dynamodbav:"expiredAt,omitempty" firestore:"expiredAt,omitempty"`
	Skills         []string          `json:"skills,omitempty" gorm:"column:skills" dynamodbav:"skills,omitempty" firestore:"skills,omitempty"`
	Location       *string           `json:"location,omitempty" gorm:"column:location" dynamodbav:"location,omitempty" firestore:"location,omitempty"`
	Quantity       *int32            `json:"quantity,omitempty" gorm:"column:quantity" dynamodbav:"quantity,omitempty" firestore:"quantity,omitempty"`
	ApplicantCount *int32            `json:"applicantCount,omitempty" gorm:"column:applicant_count" dynamodbav:"applicantCount,omitempty" firestore:"applicantCount,omitempty"`
	CompanyId      string            `json:"companyId,omitempty" gorm:"column:company_id" dynamodbav:"companyid,omitempty" firestore:"companyid,omitempty"`
}
