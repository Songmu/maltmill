package maltmill

import (
	"bytes"
	"os"
	"testing"
)

func TestTmplExecute(t *testing.T) {
	testCases := []struct {
		name       string
		nf         *formulaData
		expectFile string
	}{
		{
			"darwin amd64 only",
			&formulaData{
				Owner:           "Songmu",
				Repo:            "maltmill",
				Name:            "maltmill",
				Version:         "0.5.7",
				CapitalizedName: "Maltmill",
				Downloads: formulaDataDownloads{
					DarwinAmd64: &formulaDownload{
						URL:    "https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_darwin_amd64.zip",
						SHA256: "2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e2",
						OS:     "darwin",
						Arch:   "amd64",
					},
				},
			},
			"testdata/tmpl_darwin_amd_64_only.rb",
		},
		{
			"darwin only",
			&formulaData{
				Owner:           "Songmu",
				Repo:            "maltmill",
				Name:            "maltmill",
				Version:         "0.5.7",
				CapitalizedName: "Maltmill",
				Downloads: formulaDataDownloads{
					DarwinAmd64: &formulaDownload{
						URL:    "https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_darwin_amd64.zip",
						SHA256: "2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e2",
						OS:     "darwin",
						Arch:   "amd64",
					},
					DarwinArm64: &formulaDownload{
						URL:    "https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_darwin_arm64.zip",
						SHA256: "2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e3",
						OS:     "darwin",
						Arch:   "arm64",
					},
				},
			},
			"testdata/tmpl_darwin_only.rb",
		},
		{
			"amd64 only",
			&formulaData{
				Owner:           "Songmu",
				Repo:            "maltmill",
				Name:            "maltmill",
				Version:         "0.5.7",
				CapitalizedName: "Maltmill",
				Downloads: formulaDataDownloads{
					DarwinAmd64: &formulaDownload{
						URL:    "https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_darwin_amd64.zip",
						SHA256: "2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e2",
						OS:     "darwin",
						Arch:   "amd64",
					},
					LinuxAmd64: &formulaDownload{
						URL:    "https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_linux_amd64.tar.gz",
						SHA256: "2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e4",
						OS:     "linux",
						Arch:   "amd64",
					},
				},
			},
			"testdata/tmpl_amd64_only.rb",
		},
		{
			"linux amd64 only",
			&formulaData{
				Owner:           "Songmu",
				Repo:            "maltmill",
				Name:            "maltmill",
				Version:         "0.5.7",
				CapitalizedName: "Maltmill",
				Downloads: formulaDataDownloads{
					LinuxAmd64: &formulaDownload{
						URL:    "https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_linux_amd64.tar.gz",
						SHA256: "2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e4",
						OS:     "linux",
						Arch:   "amd64",
					},
				},
			},
			"testdata/tmpl_linux_amd64_only.rb",
		},
		{
			"with desc",
			&formulaData{
				Desc:            "create and update Homebrew thrid party Formulae",
				Owner:           "Songmu",
				Repo:            "maltmill",
				Name:            "maltmill",
				Version:         "0.5.7",
				CapitalizedName: "Maltmill",
				Downloads: formulaDataDownloads{
					DarwinAmd64: &formulaDownload{
						URL:    "https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_darwin_amd64.zip",
						SHA256: "2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e2",
						OS:     "darwin",
						Arch:   "amd64",
					},
				},
			},
			"testdata/tmpl_with_desc.rb",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			if err := formulaTmpl.Execute(out, tc.nf); err != nil {
				t.Fatal(err)
			}
			if os.Getenv("UPDATE_GOLDEN") != "" {
				if err := os.WriteFile(tc.expectFile, out.Bytes(), os.ModePerm); err != nil {
					t.Fatal(err)
				}
				return
			}
			b, err := os.ReadFile(tc.expectFile)
			if err != nil {
				t.Fatal(err)
			}

			expect := string(b)
			if out.String() != expect {
				t.Errorf("result not expected.\n  out: %s\nexpect: %s", out.String(), expect)
			}
		})
	}
}
