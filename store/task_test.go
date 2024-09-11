package store

import (
	"context"
	"testing"

	"example.com/sample/go_todo_app/clock"
	"example.com/sample/go_todo_app/entity"
	"example.com/sample/go_todo_app/testutil"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()

	// entity.Taskを作成するほかのテストケースと混ざるとテストがフェイルする
	// そのため、トランザクションを張ることでこのテストケースの中だけのテーブル状態にする
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// このテストケースが完了したら元に戻す
	t.Cleanup(func() {
		_ = tx.Rollback()
	})
	if err != nil {
		t.Fatal(err)
	}
	wants := prepareTasks(ctx, t, tx, tx)

	sut := &Repository{}
	tx.QueryRowContext(ctx, "SELECT NOW();").Scan(&sut.Clocker)
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want\n%s", d)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer, queryer Queryer) entity.Tasks {
	t.Helper()
	if _, err := con.ExecContext(ctx, "DELETE FROM tasks;"); err != nil {
		t.Logf("failed to delete tasks: %v", err)
	}
	if _, err := con.ExecContext(ctx, "TRUNCATE TABLE tasks RESTART IDENTITY;"); err != nil {
		t.Logf("failed to restart sequence: %v", err)
	}

	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title:     "want task 1",
			Status:    "todo",
			CreatedAt: c.Now(),
		},
		{
			Title:     "want task 2",
			Status:    "doing",
			CreatedAt: c.Now(),
		},
		{
			Title:     "want task 3",
			Status:    "done",
			CreatedAt: c.Now(),
		},
	}

	var id int
	err := queryer.QueryRowxContext(
		ctx,
		`INSERT INTO tasks (title, status, created_at)
		 VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9)
		 RETURNING id;
		`,
		wants[0].Title, wants[0].Status, wants[0].CreatedAt,
		wants[1].Title, wants[1].Status, wants[1].CreatedAt,
		wants[2].Title, wants[2].Status, wants[2].CreatedAt,
	).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}

	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:     "ok task",
		Status:    "todo",
		CreatedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { db.Close() })
	mock.ExpectQuery(
		`
		INSERT INTO tasks \(title, status, created_at\)
		VALUES \(\$1, \$2, \$3\)
		returning id
		`,
	).WithArgs(okTask.Title, okTask.Status, okTask.CreatedAt).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(wantID),
	)

	xdb := sqlx.NewDb(db, "postgres")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Fatal(err)
	}
}
