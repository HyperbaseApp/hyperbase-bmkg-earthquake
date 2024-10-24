package main

import (
	"bmkgearthquakecollector/collector"
	hyperbaseclient "bmkgearthquakecollector/hyperbase"
	"bmkgearthquakecollector/model"
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mnaufalhilmym/goasync"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	hyperbaseBaseURL := os.Getenv("HYPERBASE_BASE_URL")

	hyperbaseAuthTokenID := uuid.MustParse(os.Getenv("HYPERBASE_AUTH_TOKEN_ID"))
	hyperbaseAuthToken := os.Getenv("HYPERBASE_AUTH_TOKEN")
	hyperbaseAuthCollectionID := uuid.MustParse(os.Getenv("HYPERBASE_AUTH_COLLECTION_ID"))
	hyperbaseAuthCredential := map[string]any{
		"username": os.Getenv("HYPERBASE_AUTH_CREDENTIAL_USERNAME"),
		"password": os.Getenv("HYPERBASE_AUTH_CREDENTIAL_PASSWORD"),
	}

	hyperbaseProjectID := uuid.MustParse(os.Getenv("HYPERBASE_PROJECT_ID"))
	hyperbaseCollectionID := uuid.MustParse(os.Getenv("HYPERBASE_COLLECTION_ID"))

	schedulerInterval, err := time.ParseDuration(os.Getenv("SCHEDULER_INTERVAL"))
	if err != nil {
		panic(err)
	}

	hyperbase := hyperbaseclient.New(hyperbaseBaseURL)
	hyperbase.Authenticate(
		hyperbaseAuthTokenID,
		hyperbaseAuthToken,
		hyperbaseAuthCollectionID,
		hyperbaseAuthCredential,
	)
	hyperbaseProject := hyperbase.SetProject(hyperbaseProjectID)
	hyperbaseCollection := hyperbaseProject.SetCollection(hyperbaseCollectionID)

	scheduler(schedulerInterval, hyperbaseCollection)
}

func scheduler(d time.Duration, hyperbaseCollection *hyperbaseclient.HyperbaseCollection) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		log.Println("Scheduled to run at", time.Now().Add(d))
		<-ticker.C // Waiting for the next schedule to run

		autoGempaTask := goasync.Spawn(func(ctx context.Context) (*model.AutoGempaModel, error) {
			return collector.AutoGempa()
		})
		gempaTerkiniTask := goasync.Spawn(func(ctx context.Context) (*model.DataModel, error) {
			return collector.GempaTerkini()
		})
		gempaDirasakanTask := goasync.Spawn(func(ctx context.Context) (*model.DataModel, error) {
			return collector.GempaDirasakan()
		})

		autoGempaData, err := autoGempaTask.Await(context.Background())
		if err != nil {
			panic(err)
		}

		data, err := goasync.TryJoin(context.Background(), gempaTerkiniTask, gempaDirasakanTask)
		if err != nil {
			panic(err)
		}

		if err := hyperbaseCollection.InsertOne(autoGempaData.ToMap()); err != nil {
			if !errors.Is(err, hyperbaseclient.ErrDuplicate) {
				panic(err)
			} else {
				log.Println("Duplicate autogempa data", "datetime:", autoGempaData.Infogempa.Gempa.DateTime, "coordinates:", autoGempaData.Infogempa.Gempa.Coordinates)
			}
		} else {
			log.Println("Successfully save autogempa data", "datetime:", autoGempaData.Infogempa.Gempa.DateTime, "coordinates:", autoGempaData.Infogempa.Gempa.Coordinates)
		}

		for _, resData := range data {
			for _, d := range resData.ToSliceOfMap() {
				if err := hyperbaseCollection.InsertOne(d); err != nil {
					if !errors.Is(err, hyperbaseclient.ErrDuplicate) {
						panic(err)
					} else {
						log.Println("Duplicate gempa data", "datetime:", d["datetime"], "coordinates:", d["coordinates"])
					}
				} else {
					log.Println("Successfully save gempa data", "datetime:", d["datetime"], "coordinates:", d["coordinates"])
				}
			}
		}
	}
}
