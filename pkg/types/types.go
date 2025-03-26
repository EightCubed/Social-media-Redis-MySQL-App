package v1alpha1

import (
	models "go-social-media/pkg/models"
)

type PostReturnType struct {
	Post          models.Post
	NumberOfLikes int64
}
