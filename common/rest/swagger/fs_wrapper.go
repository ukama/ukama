package swagger

import (
	"bytes"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

//go:embed swagger-ui
var swaggerUiRes embed.FS

// AddOpenApiUIHandler adds handler that serves html for Swagger UI
func AddOpenApiUIHandler(e *gin.Engine, path string, openApiJsonPath string) {
	sub, err := fs.Sub(swaggerUiRes, "swagger-ui")
	if err != nil {
		panic(err)
	}

	e.StaticFS(path, http.FS(newFsWrapper(sub, openApiJsonPath)))
}

type fsWrapper struct {
	fs           fs.FS
	jsonFileName string
}

func newFsWrapper(fs fs.FS, jsonFilePath string) *fsWrapper {
	return &fsWrapper{
		fs:           fs,
		jsonFileName: jsonFilePath,
	}
}

func (f *fsWrapper) Open(name string) (fs.File, error) {
	if strings.HasSuffix(name, "index.html") {
		tpl, err := template.ParseFS(f.fs, "index.gohtml")
		if err != nil {
			return nil, err
		}
		data := map[string]interface{}{
			"openApiJson": strings.TrimPrefix(f.jsonFileName, "/"),
		}

		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, data); err != nil {
			return nil, err
		}

		return fsWrapperFile{
			r: io.NopCloser(bytes.NewReader(buf.Bytes())),
			fInfo: fileInfoMock{
				name: name,
				size: int64(buf.Len()),
			},
		}, nil
	}

	return f.fs.Open(name)
}

type fsWrapperFile struct {
	r     io.ReadCloser
	fInfo fs.FileInfo
}

func (f fsWrapperFile) Stat() (fs.FileInfo, error) {
	return f.fInfo, nil
}

func (f fsWrapperFile) Read(bytes []byte) (int, error) {
	return f.r.Read(bytes)
}

func (f fsWrapperFile) Close() error {
	return f.r.Close()
}

type fileInfoMock struct {
	name string
	size int64
}

func (f fileInfoMock) Name() string {
	return f.name
}

func (f fileInfoMock) Size() int64 {
	return f.size
}

func (f fileInfoMock) Mode() fs.FileMode {
	return fs.FileMode(0)
}

func (f fileInfoMock) ModTime() time.Time {
	return time.Now()
}

func (f fileInfoMock) IsDir() bool {
	return false
}

func (f fileInfoMock) Sys() interface{} {
	panic("implement me")
}
