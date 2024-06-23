package entity

//go:generate go-enum --marshal --values

// ENUM(outpaint, sessile, pendunculated)
type Domain int

func (d Domain) ImageFn() string {
	switch d {
	case DomainSessile:
		return "sessile_imgs"
	case DomainPendunculated:
		return "penduculated_imgs"
	}

	panic("not found domain")
}

func (d Domain) PromptFn() string {
	switch d {
	case DomainSessile:
		return "sessile_prompt"
	case DomainPendunculated:
		return "penduculated_imgs"
	}

	panic("not found domain")
}
