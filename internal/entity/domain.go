package entity

//go:generate go-enum --marshal --values

// ENUM(sessile, pendunculated, human, landscape, simple_animal)
type Domain int

func (d Domain) ImageFn() string {
	switch d {
	case DomainSessile:
		return "sessile_imgs"
	case DomainPendunculated:
		return "penduculated_imgs"
	case DomainHuman:
		return "human_imgs"
	case DomainLandscape:
		return "landscape_imgs"
	case DomainSimpleAnimal:
		return ""
	}

	panic("not found domain")
}

func (d Domain) SampleScriptPath() string {
	switch d {
	case DomainSessile, DomainPendunculated:
		return "./d3po/scripts/sample_inpaint.py"
	case DomainHuman, DomainLandscape:
		return "./d3po/scripts/sample_outpaint.py"
	case DomainSimpleAnimal:
		return "./d3po/scripts/sample_txt2img.py"
	}

	panic("not found domain")
}

func (d Domain) TrainScriptPath() string {
	switch d {
	case DomainSimpleAnimal:
        return "./d3po/scripts/train_txt2img.py"
	default:
        return "./d3po/scripts/train_inpaint.py"
	}
}

func (d Domain) PromptFn() string {
	switch d {
	case DomainSessile:
		return "sessile_prompt"
	case DomainPendunculated:
		return "penduculated_prompt"
	case DomainHuman:
		return "human_prompt"
	case DomainLandscape:
		return "landscape_prompt"
	case DomainSimpleAnimal:
		return "simple_animal"
	}

	panic("not found domain")
}
