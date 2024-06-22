package entity

const (
	data_dir = "./data"

	inf_dir         = "inferences"
	inf_mask_path   = "image/mask.png"
	inf_image_path  = "images/input.png"
	inf_output_path = "images/output.png"
)

type Task interface {
	Sample | Train | Inference
}
