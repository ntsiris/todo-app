package api

import (
	"fmt"
	"github.com/ntsiris/todo-app/internal/service"
	"github.com/ntsiris/todo-app/internal/store"
	"github.com/ntsiris/todo-app/internal/types"
	"github.com/ntsiris/todo-app/internal/utils"
	"net/http"
	"strconv"
)

type TodoHandler struct {
	*service.TodoService
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			_ = utils.WriteJSON(w, err.(*types.APIError).Code, err)
		}
	}
}

func NewTodoHandler(store store.Store) *TodoHandler {
	return &TodoHandler{
		service.NewTodoService(store),
	}
}

func (handler *TodoHandler) RegisterRoutes(router *http.ServeMux) {
	// Register Add
	router.HandleFunc("POST /todo", makeHTTPHandleFunc(handler.handleAdd))

	// Register Get
	// TODO: get all should be paginated
	router.HandleFunc("GET /todo", makeHTTPHandleFunc(handler.handleGetAll))
	router.HandleFunc("GET /todo/{id}", makeHTTPHandleFunc(handler.handleGet))

	// Register Update
	router.HandleFunc("PUT /todo/{id}", makeHTTPHandleFunc(handler.handleUpdate))

	// Register Delete
	router.HandleFunc("DELETE /todo/{id}", makeHTTPHandleFunc(handler.handleDelete))

	// Register Search
	router.HandleFunc("GET /todo/search", makeHTTPHandleFunc(handler.handleSearch))
}

func (handler *TodoHandler) handleAdd(w http.ResponseWriter, r *http.Request) error {
	// Parse the request
	var todoItem service.Item
	if err := utils.ParseJSON(r, &todoItem); err != nil {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "Failed to add a new item",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	// Add the item to the store
	err := handler.Add(&todoItem)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "failed to add a new item",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	return utils.WriteJSON(w, http.StatusCreated, todoItem)
}

func (handler *TodoHandler) handleGetAll(w http.ResponseWriter, r *http.Request) error {
	todoItems, err := handler.GetAll()

	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "failed to retrieve all items",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	return utils.WriteJSON(w, http.StatusOK, todoItems)
}

func (handler *TodoHandler) handleGet(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIntPathValue(r, "id")
	if err != nil {
		return err
	}

	item, err := handler.retrieveItem(r, id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, item)
}

func (handler *TodoHandler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	// Get the id from the request
	id, err := parseIntPathValue(r, "id")
	if err != nil {
		return err
	}

	// Get the updated item from the request
	updatedItem := service.Item{}
	if err = utils.ParseJSON(r, &updatedItem); err != nil {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "Invalid request format",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	// Update the item to the store
	if err = handler.Update(id, &updatedItem); err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "Failed to update item",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	return utils.WriteJSON(w, http.StatusOK, updatedItem)
}

func (handler *TodoHandler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIntPathValue(r, "id")
	if err != nil {
		return err
	}

	item, err := handler.retrieveItem(r, id)
	if err != nil {
		return err
	}

	err = handler.Delete(item)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       fmt.Sprintf("failed to delete item with id %v", id),
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	return utils.WriteJSON(w, http.StatusOK, item)
}

func (handler *TodoHandler) handleSearch(w http.ResponseWriter, r *http.Request) error {
	// Get query parameter from URL
	query := r.URL.Query().Get("q")
	if query == "" {
		return &types.APIError{
			Code:      http.StatusBadRequest,
			Message:   "query parameter 'q' is required",
			Operation: types.FormatOperation(r.Method, r.URL.Path),
		}
	}

	results, err := handler.Search(query)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       fmt.Sprintf("failed to search item with query %v", query),
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	return utils.WriteJSON(w, http.StatusOK, results)
}

func (handler *TodoHandler) retrieveItem(r *http.Request, id int) (service.Item, error) {

	item, err := handler.Get(id)
	if err != nil {
		return service.Item{}, &types.APIError{
			Code:          http.StatusNotFound,
			Message:       fmt.Sprintf("failed to retrieve item with id %v", id),
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	return *item, nil
}

func parseIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)

	intValue, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "invalid format of ID value",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err,
		}
	}

	return intValue, nil
}
