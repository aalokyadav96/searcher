package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Baito struct {
	BaitoId      string    `bson:"baitoid,omitempty" json:"baitoid"`
	Title        string    `bson:"title" json:"title"`
	Description  string    `bson:"description" json:"description"`
	Category     string    `bson:"category" json:"category"`
	SubCategory  string    `bson:"subcategory" json:"subcategory"`
	Location     string    `bson:"location" json:"location"`
	Wage         string    `bson:"wage" json:"wage"`
	Phone        string    `bson:"phone" json:"phone"`
	Requirements string    `bson:"requirements" json:"requirements"`
	BannerURL    string    `bson:"banner,omitempty" json:"banner,omitempty"`
	Images       []string  `bson:"images" json:"images"`
	WorkHours    string    `bson:"workHours" json:"workHours"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	OwnerID      string    `bson:"ownerId" json:"ownerId"`
}

type BaitoApplication struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BaitoID     string             `bson:"baitoid" json:"baitoid"`
	UserID      string             `bson:"userid" json:"userid"`
	Username    string             `bson:"username" json:"username"`
	Pitch       string             `bson:"pitch" json:"pitch"`
	SubmittedAt time.Time          `bson:"submittedAt" json:"submittedAt"`
}

type BaitoWorker struct {
	UserID      string    `json:"userid" bson:"userid"`
	BaitoUserID string    `json:"baito_user_id" bson:"baito_user_id"`
	Name        string    `json:"name" bson:"name"`
	Age         int       `json:"age" bson:"age"`
	Phone       string    `json:"phone_number" bson:"phone_number"`
	Location    string    `json:"address" bson:"address"`
	Preferred   string    `json:"preferred_roles" bson:"preferred_roles"`
	Bio         string    `json:"bio" bson:"bio"`
	ProfilePic  string    `json:"profile_picture" bson:"profile_picture"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
