package dto

type CreateLinkDTO struct {
	OriginalUrl *string
	ShortName   *string
	ShortUrl    *string
}

type UpdateLinkDTO struct {
	ShortName   *string
	OriginalUrl *string
}
