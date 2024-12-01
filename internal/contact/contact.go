package contact

import "time"

type Contact struct {
	Id          string     `yaml:"id" mapstructure:"id" json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"-" avro:"id" validate:"max=40" operator:"="`
	Name        string     `yaml:"name" mapstructure:"name" json:"name" gorm:"column:name" bson:"name" dynamodbav:"name" firestore:"name" avro:"name" validate:"required,max=100"`
	Country     string     `yaml:"country" mapstructure:"country" json:"country" gorm:"column:country" bson:"country" dynamodbav:"country" firestore:"country" avro:"country" validate:"required,max=100"`
	Company     string     `yaml:"company" mapstructure:"company" json:"company" gorm:"column:company" bson:"company" dynamodbav:"company" firestore:"company" avro:"company" validate:"required,max=100"`
	JobTitle    string     `yaml:"job_title" mapstructure:"job_title" json:"jobTitle" gorm:"column:job_title" bson:"jobTitle" dynamodbav:"jobTitle" firestore:"jobTitle" avro:"jobTitle" validate:"required,max=100"`
	Email       string     `yaml:"email" mapstructure:"email" json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" avro:"email" validate:"email,max=120"`
	Phone       string     `yaml:"phone" mapstructure:"phone" json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" avro:"phone" validate:"required,phone,max=18" operator:"like"`
	Message     string     `yaml:"message" mapstructure:"message" json:"message" gorm:"column:message" bson:"message" dynamodbav:"message" firestore:"message" avro:"message" validate:"required,max=400"`
	SubmittedAt *time.Time `yaml:"submitted_at" json:"submittedAt,omitempty" gorm:"column:submitted_at" bson:"submittedAt,omitempty" dynamodbav:"submittedAt,omitempty" firestore:"submittedAt,omitempty"`
	ContactedAt *time.Time `yaml:"contacted_at" json:"contactedAt,omitempty" gorm:"column:contacted_at" bson:"contactedAt,omitempty" dynamodbav:"contactedAt,omitempty" firestore:"contactedAt,omitempty"`
	ContactedBy *string    `yaml:"contacted_by" json:"contactedBy,omitempty" gorm:"column:contacted_by" bson:"contactedBy,omitempty" dynamodbav:"contactedBy,omitempty" firestore:"contactedBy,omitempty"`
}
