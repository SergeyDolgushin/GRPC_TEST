package converter

import (
	"github.com/SergeyDolgushin/GRPC_TEST/internal/repository/note/model"
	desc "github.com/SergeyDolgushin/GRPC_TEST/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromRepo(note *model.Note) *desc.Note {
	var updated_at *timestamppb.Timestamp
	if note.Updated_at.Valid {
		updated_at = timestamppb.New(note.Updated_at.Time)
	}

	return &desc.Note{
		Id:        note.ID,
		Info:      ToNoteInfoFromRepo(note.Info),
		CreatedAt: timestamppb.New(note.Created_at),
		UpdatedAt: updated_at,
	}
}

func ToNoteInfoFromRepo(info *model.Info) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
