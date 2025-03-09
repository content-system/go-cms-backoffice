package category

type Category struct {
	Id       string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id;primary_key" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name     string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Path     string `yaml:"path" mapstructure:"path" json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty"`
	Resource string `yaml:"resource" mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource_key" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Icon     string `yaml:"icon" mapstructure:"icon" json:"icon,omitempty" gorm:"column:icon" bson:"icon,omitempty" dynamodbav:"icon,omitempty" firestore:"icon,omitempty"`
	Sequence int    `yaml:"sequence" mapstructure:"sequence" json:"sequence,omitempty" gorm:"column:sequence" bson:"sequence" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
	Type     string `yaml:"type" mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Parent   string `yaml:"parent" mapstructure:"parent" json:"parent,omitempty" gorm:"column:parent" bson:"parent,omitempty" dynamodbav:"parent,omitempty" firestore:"parent,omitempty"`
	Status   string `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
}
