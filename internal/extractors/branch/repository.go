package branch

import (
	"context"
	"time"

	"github.com/CesarDelgadoM/extractor-reports/internal/models"
	"github.com/CesarDelgadoM/extractor-reports/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

type IBranchRepository interface {
	Find(userid uint, name string) (*models.Restaurant, error)
	Size(userid uint, name string) int
	GetPage(userid uint, name string, skip, limit int64) (*[]models.Branch, error)
}

type branchRepository struct {
	mongodb *database.MongoDB
}

func NewBranchRepository(mongodb *database.MongoDB) IBranchRepository {
	return &branchRepository{
		mongodb,
	}
}

func (r *branchRepository) Find(userid uint, name string) (*models.Restaurant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data models.Restaurant

	filter := bson.M{
		"userid": userid,
		"name":   name,
	}

	err := r.mongodb.CollectionRestaurant().FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r branchRepository) Size(userid uint, name string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipe := bson.A{

		bson.M{
			"$match": bson.M{
				"userid": userid,
				"name":   name,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id": 0,
				"number": bson.M{
					"$size": "$branches",
				},
			},
		},
	}

	cursor, err := r.mongodb.CollectionRestaurant().Aggregate(ctx, pipe)
	if err != nil {
		return -1
	}

	size := struct {
		Number int `json:"number"`
	}{}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&size); err != nil {
			return -1
		}
	}

	return size.Number
}

func (r *branchRepository) GetPage(userid uint, name string, skip, limit int64) (*[]models.Branch, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipe := bson.A{

		bson.M{
			"$match": bson.M{
				"userid": userid,
				"name":   name,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id": 0,
				"branches": bson.M{
					"$slice": bson.A{
						"$branches",
						skip,
						limit,
					},
				},
			},
		},
	}

	cursor, err := r.mongodb.CollectionRestaurant().Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	result := struct {
		Branches []models.Branch `json:"branches"`
	}{}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	return &result.Branches, nil
}
