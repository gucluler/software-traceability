package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	db "traceability/database"

	guuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ViewKind string

const (
	UserStory   ViewKind = "userStory"
	Functional  ViewKind = "functional"
	Development ViewKind = "development"
	None        ViewKind = "none"
)

// ArchViewComponent is the component of a view
// swagger:model
type ArchViewComponent struct {
	// the id for the member
	//
	// required: false
	ID string `json:"id"`

	// UserKind is used for architecture user stories to make things simpler
	//
	// required: false
	UserKind string `json:"userKind"`

	// Kind, "userStory", "functional", "development"
	//
	// required: false
	Kind ViewKind `json:"kind" validate:"oneof=userStory functional development"`

	LinksList []string `json:"links,omitempty"`
	// required: true
	Desctription string `json:"description" validate:"required"`
	// component belongs to view with id
	// required: true
	ViewID string `json:"viewID" validate:"required"`

	// component belongs to view with id
	// required: true
	ProjectID string `json:"projectID" validate:"required"`

	//FunctionList is used for development view to show functions of a component
	FunctionList []string `json:"functions"`

	//VarList is used for development view to show variables of a component
	VarList []string `json:"variables"`

	// Comments will be in here
}

// UnmarshalJSON parses from json
func (ac *ArchViewComponent) UnmarshalJSON(data []byte) error {
	// Define a secondary tycursive cape so that we don't end up with a rell to json.Unmarshal
	type Aux ArchViewComponent
	var a *Aux = (*Aux)(ac)
	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	// Validate the valid enum values
	switch ac.Kind {
	case UserStory, Functional, Development, None:
		return nil
	default:
		ac.Kind = ""
		return errors.New("invalid value for Key")
	}
}

// ArchView general purpose architecture view
// swagger:model
type ArchView struct {
	// the id for the project
	//
	// required: false
	ID string `json:"id,omitempty" bson:"omitempty"`

	// the name of the project
	//
	// required: true
	// max length: 30
	Name string `json:"name" validate:"required"`

	// belonging project's id
	//
	// required: true
	ProjectID string `json:"projectID"`

	// Kind, "userStory", "functional", "development"
	//
	// required: false
	Kind string `json:"kind" validate:"oneof=userStory functional development"`

	// description
	//
	// required: false
	Desctription string `json:"description"`

	// Component IDs of the view
	//
	// required: false
	Components []string `json:"components,omitempty" bson:"omitempty"`
}

// AddArchView adds a new project to the database
func AddArchView(v ArchView) error {
	v.ID = guuid.New().String()

	collection := db.DB.Collection(db.ArchViewCollectionName)
	insertResult, err := collection.InsertOne(context.TODO(), v)

	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

// FindArchViewByID returns an ArchView or error
func FindArchViewByID(id string) (ArchView, error) {
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()

	collection := db.DB.Collection(db.ArchViewCollectionName)

	var resultArchView ArchView

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	err := collection.FindOne(ctx, filter).Decode(&resultArchView)
	return resultArchView, err
}

// AddArchViewComponent adds component to the ArchView
func AddArchViewComponent(c ArchViewComponent) {
	c.ID = guuid.New().String()
	archViewID := c.ViewID
	archViewCollection := db.DB.Collection(db.ArchViewCollectionName)
	componentCollection := db.DB.Collection(db.ArchViewComponentCollectionName)

	query := bson.M{"id": archViewID}
	update := bson.M{"$push": bson.M{"components": c.ID}}

	updateResult, err := archViewCollection.UpdateOne(context.TODO(), query, update)
	insertResult, err := componentCollection.InsertOne(context.TODO(), c)

	if err != nil {
		panic(err)
	}

	fmt.Println("Upserted a single document:", updateResult, "\n Inserted a single document: ", insertResult)
}

// FindArchViewComponentByID returns an ArchView or error
func FindArchViewComponentByID(id string, archViewID string) (ArchViewComponent, error) {
	exp := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()
	collection := db.DB.Collection(db.ArchViewComponentCollectionName)

	var resultComponent ArchViewComponent

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	err := collection.FindOne(ctx, filter).Decode(&resultComponent)
	return resultComponent, err
}
