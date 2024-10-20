package contact

import (
	"github.com/core-go/search"
)

type ContactFilter struct {
	*search.Filter
	Id          string            `yaml:"id" mapstructure:"id" json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"-" avro:"id" validate:"max=40" operator:"="`
	Name        string            `yaml:"name" mapstructure:"name" json:"name" gorm:"column:name" bson:"name" dynamodbav:"name" firestore:"name" avro:"name" validate:"required,max=100"`
	Country     string            `yaml:"country" mapstructure:"country" json:"country" gorm:"column:country" bson:"country" dynamodbav:"country" firestore:"country" avro:"country" validate:"required,max=100"`
	Company     string            `yaml:"company" mapstructure:"company" json:"company" gorm:"column:company" bson:"company" dynamodbav:"company" firestore:"company" avro:"company" validate:"required,max=100"`
	JobTitle    string            `yaml:"job_title" mapstructure:"job_title" json:"jobTitle" gorm:"column:job_title" bson:"jobTitle" dynamodbav:"jobTitle" firestore:"jobTitle" avro:"jobTitle" validate:"required,max=100"`
	Email       string            `yaml:"email" mapstructure:"email" json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" avro:"email" validate:"email,max=120"`
	Phone       string            `yaml:"phone" mapstructure:"phone" json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" avro:"phone" validate:"required,phone,max=18" operator:"like"`
	SubmittedAt *search.TimeRange `yaml:"submitted_at" json:"submittedAt,omitempty" gorm:"column:submitted_at" bson:"submittedAt,omitempty" dynamodbav:"submittedAt,omitempty" firestore:"submittedAt,omitempty"`
}
