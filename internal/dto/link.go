package dto

type GetLinksDTO struct {
	Limit  *int64
	Offset *int64
}

type CreateLinkDTO struct {
	OriginalUrl *string
	ShortName   *string
	ShortUrl    *string
}

type UpdateLinkDTO struct {
	ShortName   *string
	OriginalUrl *string
}
