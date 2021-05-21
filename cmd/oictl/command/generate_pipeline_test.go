package command

import (
	"os"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

func Test_splitServiceIntoParts(t *testing.T) {
	tests := []struct {
		svc     string
		name    string
		version string
		alias   string
		errMsg  string
	}{
		{
			svc:     "mysvc",
			name:    "mysvc",
			version: "0.0.1",
			alias:   "mysvc",
		},
		{
			svc:     "mysvc:0.0.2",
			name:    "mysvc",
			version: "0.0.2",
			alias:   "mysvc",
		},
		{
			svc:     "mysvc@svc",
			name:    "mysvc",
			version: "0.0.1",
			alias:   "svc",
		},
		{
			svc:     "mysvc:0.0.2@svc",
			name:    "mysvc",
			version: "0.0.2",
			alias:   "svc",
		},
		{
			svc:    "",
			errMsg: "name must be part of service",
		},
		{
			svc:    ":0.0.1",
			errMsg: "name must be part of service",
		},
		{
			svc:    "mysvc:v1",
			errMsg: "version must be semantic version: No Major.Minor.Patch elements found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, version, alias, err := splitServiceIntoParts(tt.svc)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			}
			assert.Equal(t, tt.name, name)
			assert.Equal(t, tt.version, version)
			assert.Equal(t, tt.alias, alias)
		})
	}
}

func Test_resolveProjectFinalLocation(t *testing.T) {
	tests := []struct {
		name   string
		dir    string
		want   string
		errMsg string
		statFn func(name string) (os.FileInfo, error)
	}{
		{
			name: ". should use current working dir",
			dir:  ".",
			want: "/some/path",
		},
		{
			name: "./inner should prepand to working dir",
			dir:  "./inner",
			want: "/some/path/inner",
		},
		{
			name: "/path root shoud be absolute",
			dir:  "/path",
			want: "/path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getwd = func() (string, error) {
				return "/some/path", nil
			}
			stat = tt.statFn
			if stat == nil {
				stat = func(name string) (os.FileInfo, error) {
					memfs := afero.NewMemMapFs()
					if err := memfs.MkdirAll(name, os.ModePerm); err != nil {
						return nil, err
					}
					return memfs.Stat(name)
				}
			}

			got, err := resolveProjectFinalLocation(tt.dir)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
