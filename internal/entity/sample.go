package entity

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/database"
)

type Sample struct {
	model *Model

	id string
}

func NewSample(
	sourceModel *Model,
	init bool,
) (*Sample, error) {
	if !init {
		return retrieve(sourceModel), nil
	}
	sampleID := uuid.NewString()
	newModel, err := sourceModel.NewChild(sampleID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &Sample{
		id:    sampleID,
		model: newModel,
	}, nil
}

func retrieve(model *Model) *Sample {
	return &Sample{
		id:    model.sampleID,
		model: model,
	}
}

func (i *Sample) TaskContent() string {
	return fmt.Sprintf("sampling for model \"%s\"", i.model.name)
}

func (i *Sample) TaskType() string {
	return "sample"
}

func (i *Sample) Estimate() time.Duration {
	return time.Minute * 5
}

func (s *Sample) Insertion() database.InsertSampleParams {
	return database.InsertSampleParams{
		ID:      s.id,
		ModelID: s.model.id,
	}
}

func (s *Sample) Model() *Model {
	return s.model
}

func (s *Sample) SaveDir() string {
	return path.Join("./data", "assets/samples", s.id)
}

func (s *Sample) ViewImage() string {
	files, err := readFilesInFolder(fmt.Sprintf("./data/assets/samples/%s/images", s.id))
	if err != nil {
		fmt.Println("should not err here")
	}

	return path.Join("/data", "assets/samples", s.id, "images", files[0])
}

func (s *Sample) SampleImages() []string {
	files, err := readFilesInFolder(fmt.Sprintf("./data/assets/samples/%s/images", s.id))
	if err != nil {
		fmt.Println("should not err here", err.Error())
	}

	for _, image := range files {
		fmt.Println("image: ", image)
	}

	return files
}

func (s *Sample) ID() string {
	return s.id
}

func readFilesInFolder(folderPath string) ([]string, error) {
	var fileNames []string

	folder, err := os.Open(folderPath)
	if err != nil {
		return nil, err
	}
	defer folder.Close()

	files, err := folder.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}
