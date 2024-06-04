package entity

type TaskSample struct {
	ResumeFrom      *string
	ImageFn         string
	PromptFnInpaint string
}

type TaskTrain struct {
	ResumeFrom      *string
	TrainSamplePath string
	TrainJsonPath   string
}
