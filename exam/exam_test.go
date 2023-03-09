package exam_test

import (
	"testing"

	"example.com/exam"
)

func TestParseFileName(t *testing.T) {
	type want struct {
		subjectName string
		times       int
		err         bool
	}

	type testcase struct {
		name     string
		fileName string
		want     want
	}

	cases := []testcase{
		{
			name:     "1",
			fileName: "hoge_1.csv",
			want: want{
				subjectName: "hoge",
				times:       1,
				err:         false,
			},
		},
		{
			name:     "2",
			fileName: "hoge_fuga_1_2_3.csv",
			want: want{
				subjectName: "",
				times:       0,
				err:         true,
			},
		},
		{
			name:     "3",
			fileName: "hoge_123fugapiyo.txt.csv",
			want: want{
				subjectName: "hoge",
				times:       123,
				err:         false,
			},
		},
		{
			name:     "4",
			fileName: "hoge1.csv",
			want: want{
				subjectName: "",
				times:       0,
				err:         true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			subjectName, times, err := exam.ParseFileName(tt.fileName)

			if tt.want.err {
				if err == nil {
					t.Errorf("ParseFileName is not error, want error")
				}
			} else if tt.want.subjectName != subjectName || tt.want.times != times {
				t.Errorf("got: %v, %v\nwant: %v, %v\n", subjectName, times, tt.want.subjectName, tt.want.times)
			}
		})
	}
}
