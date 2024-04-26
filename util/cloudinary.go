package util

import (
	"context"
	"medichat-be/config"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryProvider interface {
	SendFile(sendFile SendFileOpts) (*uploader.UploadResult, error)
}

type cloudinaryProviderImpl struct {
	cloud *cloudinary.Cloudinary
}

type SendFileOpts struct {
	Context  context.Context
	Filename string `json:"filename"`
	Roomid   string `json:"roomid"`
	File     multipart.File
}

func NewCloudinarylProvider() (*cloudinaryProviderImpl, error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	cld, _ := cloudinary.NewFromParams(conf.CloudinaryName, conf.CloudinaryAPIKey, conf.CloudinaryAPISecret)

	return &cloudinaryProviderImpl{
		cloud: cld,
	}, nil
}

func (p *cloudinaryProviderImpl) SendFile(sendFile SendFileOpts) (*uploader.UploadResult, error) {

	now := time.Now()

	params := uploader.UploadParams{
		Type:             api.Upload,
		ResourceType:     "auto",
		FilenameOverride: sendFile.Filename,
		PublicID:         sendFile.Roomid + now.Format("2006_01_02_T15_04_05"),
	}

	res, err := p.cloud.Upload.Upload(sendFile.Context, sendFile.File, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
