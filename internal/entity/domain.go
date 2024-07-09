package entity

//go:generate go-enum --marshal --values

// ENUM(sessile, pendunculated, human, landscape)
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
	}

	panic("not found domain")
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
	}

	panic("not found domain")
}
