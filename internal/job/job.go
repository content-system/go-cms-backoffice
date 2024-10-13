package job

import "time"

type Job struct {
	Id          string  `json:"id,omitempty" gorm:"primary_key;column:id" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"max=40"`
	Title       string  `json:"title,omitempty" gorm:"column:title" dynamodbav:"title,omitempty" firestore:"title,omitempty" validate:"omitempty,max=120"`
	Description *string `json:"description,omitempty" gorm:"column:description" dynamodbav:"description,omitempty" firestore:"description,omitempty" validate:"omitempty,max=1000"`
	// Skills         []string   `json:"skills,omitempty" gorm:"column:skills" dynamodbav:"skills,omitempty" firestore:"skills,omitempty"`
	PublishedAt    *time.Time `json:"publishedAt,omitempty" gorm:"column:published_at" dynamodbav:"publishedAt,omitempty" firestore:"publishedAt,omitempty"`
	ExpiredAt      *time.Time `json:"expiredAt,omitempty" gorm:"column:expired_at" dynamodbav:"expiredAt,omitempty" firestore:"expiredAt,omitempty"`
	Quantity       int32      `json:"quantity,omitempty" gorm:"column:quantity" dynamodbav:"quantity,omitempty" firestore:"quantity,omitempty"`
	ApplicantCount int32      `json:"applicantCount,omitempty" gorm:"column:applicant_count" dynamodbav:"applicantCount,omitempty" firestore:"applicantCount,omitempty"`
	Requirements   *string    `json:"requirements,omitempty" gorm:"column:requirements" dynamodbav:"requirements,omitempty" firestore:"requirements,omitempty"`
	Benefit        *string    `json:"benefit,omitempty" gorm:"column:benefit" dynamodbav:"benefit,omitempty" firestore:"benefit,omitempty"`
	CompanyId      string     `json:"companyId,omitempty" gorm:"column:company_id" dynamodbav:"companyid,omitempty" firestore:"companyid,omitempty"`
}
