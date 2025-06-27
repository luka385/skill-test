package report_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/luke385/skill-test/internal/report"
	"github.com/luke385/skill-test/internal/report/mocks"
	"github.com/luke385/skill-test/internal/report/usecase/domain"
)

func TestReportUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repo *mocks.MockStudentRepository
		gen  *mocks.MockFileGenerator
	}

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		wantBuf *bytes.Buffer
		wantCT  string
		wantExt string
		wantErr bool
	}{
		{
			name: "successful generate report",
			setup: func(f *fields) {
				student := &domain.Student{ID: 1, Name: "Lucas", Email: "lucas@example.com", Class: "Math"}
				buf := bytes.NewBufferString("PDF data")
				f.repo.EXPECT().GetByID("1").Return(student, nil)
				f.gen.EXPECT().Generate(student).Return(buf, nil)
				f.gen.EXPECT().GetContentType().Return("application/pdf")
				f.gen.EXPECT().GetFileExtension().Return(".pdf")
			},
			args:    args{id: "1"},
			wantBuf: bytes.NewBufferString("PDF data"),
			wantCT:  "application/pdf",
			wantExt: ".pdf",
			wantErr: false,
		},
		{
			name: "error getting student",
			setup: func(f *fields) {
				f.repo.EXPECT().GetByID("99").Return(nil, errors.New("not found"))
			},
			args:    args{id: "99"},
			wantBuf: nil,
			wantCT:  "",
			wantExt: "",
			wantErr: true,
		},
		{
			name: "error generating file",
			setup: func(f *fields) {
				student := &domain.Student{ID: 2, Name: "Ana"}
				f.repo.EXPECT().GetByID("2").Return(student, nil)
				f.gen.EXPECT().Generate(student).Return(nil, errors.New("generate error"))
			},
			args:    args{id: "2"},
			wantBuf: nil,
			wantCT:  "",
			wantExt: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewMockStudentRepository(ctrl)
			genMock := mocks.NewMockFileGenerator(ctrl)

			f := &fields{
				repo: repoMock,
				gen:  genMock,
			}

			tt.setup(f)

			gotBuf, gotCT, gotExt, err := report.NewReportUseCase(f.repo, f.gen).Execute(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, gotBuf)
				assert.Equal(t, "", gotCT)
				assert.Equal(t, "", gotExt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBuf.String(), gotBuf.String())
				assert.Equal(t, tt.wantCT, gotCT)
				assert.Equal(t, tt.wantExt, gotExt)
			}
		})
	}

}
