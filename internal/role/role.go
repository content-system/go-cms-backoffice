package role

import "time"

type Role struct {
	RoleId     string     `json:"roleId,omitempty" gorm:"column:role_id;primary_key" bson:"_id,omitempty" dynamodbav:"roleId,omitempty" firestore:"roleId,omitempty" validate:"max=40"`
	RoleName   string     `json:"roleName,omitempty" gorm:"column:role_name" bson:"roleName,omitempty" dynamodbav:"roleName,omitempty" firestore:"roleName,omitempty" validate:"required,max=255"`
	Status     string     `json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty" match:"equal" validate:"required,max=1,code"`
	Remark     *string    `json:"remark,omitempty" gorm:"column:remark" bson:"remark,omitempty" dynamodbav:"remark,omitempty" firestore:"remark,omitempty"`
	CreatedBy  *string    `json:"createdBy,omitempty" gorm:"column:created_by" bson:"createdBy,omitempty" dynamodbav:"createdBy,omitempty" firestore:"createdBy,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty" gorm:"column:created_at" bson:"createdAt,omitempty" dynamodbav:"createdAt,omitempty" firestore:"createdAt,omitempty"`
	UpdatedBy  *string    `json:"updatedBy,omitempty" gorm:"column:updated_by" bson:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty" firestore:"updatedBy,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at" bson:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty" firestore:"updatedAt,omitempty"`
	Privileges []string   `json:"privileges,omitempty" bson:"privileges,omitempty" dynamodbav:"privileges,omitempty" firestore:"privileges,omitempty"`
}
