package note

import (
	"context"
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyDolgushin/GRPC_TEST/internal/repository"
	"github.com/SergeyDolgushin/GRPC_TEST/internal/repository/note/converter"
	"github.com/SergeyDolgushin/GRPC_TEST/internal/repository/note/model"
	desc "github.com/SergeyDolgushin/GRPC_TEST/pkg/note_v1"
)

const (
	tableName = "note"

	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository.NoteRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *desc.NoteInfo) (int64, error) {

	builder := sq.Insert(tableName).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	log.Printf("CREATE REQUEST, %v, %v, %v", query, info.Title, info.Content)
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	log.Printf("CREATE REQUEST, %v", err)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.Note, error) {
	builder := sq.Select(idColumn, titleColumn, contentColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(10)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	log.Printf("GET REQUEST, %v, %v, %v", query, id, args)

	var note model.Note
	note.Info = &model.Info{}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.Created_at, &note.Updated_at)

	log.Printf("GET REQUEST, %v", err)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
