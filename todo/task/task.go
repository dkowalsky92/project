package task

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dkowalsky/todo/db"
	"github.com/dkowalsky/todo/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Task -
type Task struct {
	ID           int64      `json:"id"`
	Message      string     `json:"message"`
	Completed    bool       `json:"completed"`
	DateCreated  time.Time  `json:"date_created"`
	DateDeadline *time.Time `json:"date_deadline"`
}

// StatusChange -
type StatusChange struct {
	ID        int64 `json:"id"`
	Completed bool  `json:"completed"`
}

// Router -
type Router struct {
	Mux *chi.Mux
	Db  *db.DB
}

// NewRouter -
func NewRouter(db *db.DB) *Router {
	r := &Router{Db: db}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			//w.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(w, r)
		})
	})
	mux.Get("/{id}", r.Get)
	mux.Get("/", r.GetAll)
	mux.Post("/", r.Insert)
	mux.Put("/", r.ChangeStatus)
	mux.Options("/", r.Allow)
	mux.Delete("/{id}", r.Delete)
	mux.Options("/{id}", r.Allow)

	r.Mux = mux

	return r
}

// Allow -
func (r *Router) Allow(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
}

// DbGet -
func DbGet(db *db.DB, id int64) (*Task, error) {
	row := db.QueryRow(`SELECT id, message, completed, date_created, date_deadline FROM Task WHERE id = ?`, id)
	var task Task
	err := row.Scan(&task.ID, &task.Message, &task.Completed, &task.DateCreated, &task.DateDeadline)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Retreived 1 row.\n")

	return &task, nil
}

// Get -
func (r *Router) Get(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)
	if err != nil {
		panic(err)
	}
	tasks, err := DbGet(r.Db, id)
	if err != nil {
		panic(err)
	}
	err = util.ParseAndWrite(w, tasks)
	if err != nil {
		panic(err)
	}
}

// DbGetAll -
func DbGetAll(db *db.DB) ([]Task, error) {
	rows, err := db.Query(`SELECT id, message, completed, date_created, date_deadline FROM Task`)
	if err != nil {
		return nil, err
	}

	tasks := []Task{}

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Message, &task.Completed, &task.DateCreated, &task.DateDeadline)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	fmt.Printf("Retreived %d rows.\n", len(tasks))

	return tasks, nil
}

// GetAll -
func (r *Router) GetAll(w http.ResponseWriter, req *http.Request) {
	tasks, err := DbGetAll(r.Db)
	if err != nil {
		panic(err)
	}
	err = util.ParseAndWrite(w, tasks)
	if err != nil {
		panic(err)
	}
}

// DbInsert -
func DbInsert(db *db.DB, task Task) (*int64, error) {
	stmt, err := db.Prepare(`INSERT INTO Task(message, completed) VALUES (?, ?)`)
	if err != nil {
		return nil, err
	}

	inserted, err := stmt.Exec(task.Message, task.Completed)
	if err != nil {
		return nil, err
	}

	id, err := inserted.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// Insert -
func (r *Router) Insert(w http.ResponseWriter, req *http.Request) {
	task := &Task{}
	err := util.ReadBody(w, req, task)
	if err != nil {
		panic(err)
	}
	id, err := DbInsert(r.Db, *task)
	if err != nil {
		panic(err)
	}
	task, err = DbGet(r.Db, *id)
	if err != nil {
		panic(err)
	}
	err = util.ParseAndWrite(w, task)
	if err != nil {
		panic(err)
	}
}

// DbDelete -
func DbDelete(db *db.DB, id int64) (bool, error) {
	stmt, err := db.Prepare(`DELETE FROM task WHERE id = ?`)
	if err != nil {
		return false, err
	}

	deleted, err := stmt.Exec(id)
	if err != nil {
		return false, err
	}

	rowsAffected, _ := deleted.RowsAffected()
	fmt.Printf("Deleted %d rows.\n", rowsAffected)

	return true, nil
}

// Delete -
func (r *Router) Delete(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(req, "id"), 10, 64)

	if err != nil {
		panic(err)
	}
	success, err := DbDelete(r.Db, id)
	if err != nil {
		panic(err)
	}
	if success {
		w.Write([]byte(`OK`))
	}
}

// DbChangeStatus -
func DbChangeStatus(db *db.DB, id int64, completed bool) (bool, error) {
	stmt, err := db.Prepare(`UPDATE Task SET completed = ? WHERE id = ?`)
	if err != nil {
		return false, err
	}

	updated, err := stmt.Exec(completed, id)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := updated.RowsAffected()
	fmt.Printf("Updated %d rows.\n", rowsAffected)

	return true, nil
}

// ChangeStatus -
func (r *Router) ChangeStatus(w http.ResponseWriter, req *http.Request) {
	statusChange := &StatusChange{}
	err := util.ReadBody(w, req, statusChange)
	if err != nil {
		panic(err)
	}
	success, err := DbChangeStatus(r.Db, statusChange.ID, statusChange.Completed)
	if err != nil {
		panic(err)
	}
	if success {
		w.Write([]byte(`{"OK"}`))
	}
}
