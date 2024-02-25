package model

type File struct {
    Id, Filename, ContentType, Sha1Sum, UploadKey string
    Size, UploadTime int64
}
