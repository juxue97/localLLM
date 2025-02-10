package chatbot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"chatbot/cmd/service/auth"
	"chatbot/types"
	"chatbot/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.ChatbotStore
	userStore types.UserStore
}

func NewHandler(store types.ChatbotStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/chat", auth.WithJWTAuth(h.handleQuery, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/chatroom/chatroomIDs", auth.WithJWTAuth(h.getAllChatroomIDs, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/chatroom/chatroomHistory/{roomID}", auth.WithJWTAuth(h.getChatroomHistory, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/chatroom/create", auth.WithJWTAuth(h.createChatroom, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/chatroom/delete/{roomID}", auth.WithJWTAuth(h.deleteChatroom, h.userStore)).Methods(http.MethodDelete)
}

// TODO 1 - If not found/selected session room ID, create a new one
// TODO 2 - Load the chat history (if any)
// TODO 3 - Store the conversation to the database with respect to the roomID
func (h *Handler) handleQuery(w http.ResponseWriter, r *http.Request) {
	var payload types.ChatbotPayload
	var event map[string]interface{}
	var accumulatedContent strings.Builder

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	// fmt.Println(userID)

	// Get JSON payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// Validate Request JSON Body
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// Check if SessionRoomID is provided in the payload
	if payload.SessionRoomID == "" {
		roomID, err := h.store.CreateSessionRoom(userID.Hex())
		if err != nil {
			fmt.Fprintf(w, "data: %s\n\n", err)
			return
		}
		err = h.userStore.UpdateUserSessionRooms(userID, roomID)
		if err != nil {
			fmt.Fprintf(w, "data: %s\n\n", err)
			return
		}
		payload.SessionRoomID = roomID
	}

	// Start implement sse streaming here
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Fprintf(w, "data: %s\n\n", fmt.Errorf("streaming unsupported"))
		return
	}

	// TODO - Load previous chat history into the payload.Query
	messages, err := h.store.LoadChatHistory(userID.Hex(), payload.SessionRoomID)
	if err != nil {
		fmt.Fprintf(w, "data: %s\n\n", err)
		return
	}
	appendMessages := append(messages, map[string]string{
		"role":    "user",
		"content": payload.Query,
	})

	// Generate the response stream from the chatbot service
	responseStreamBody, err := GenerateResponse(appendMessages)
	if err != nil {
		fmt.Fprintf(w, "data: %s\n\n", fmt.Errorf("error receiving response: %w", err))
		return
	}

	decoder := json.NewDecoder(responseStreamBody.Body)
	for {
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(w, "data: %s\n\n", fmt.Errorf("error decoding response: %w", err))
			return
		}
		// Check for stream completion (done field)
		if done, ok := event["done"].(bool); ok && done {
			// 1. extract information such as token usage, and duration for loading model and fetehing response.
			// 2. return Stream completed

			streamResponse := types.StreamResponse{
				Done:    true,
				Message: "Stream finished",
				Data: types.Output{
					Response:           accumulatedContent.String(),
					InputToken:         int(event["prompt_eval_count"].(float64)),
					OutputToken:        int(event["eval_count"].(float64)),
					LoadModelDuration:  event["load_duration"].(float64) / 1e9,
					PromptEvalDuration: event["prompt_eval_duration"].(float64) / 1e9,
					EvaluateDuration:   event["eval_duration"].(float64) / 1e9,
					TotalDuration:      event["total_duration"].(float64) / 1e9,
				},
			}

			// Save the convesation history onto the mongoDB
			err := h.store.StoreChatHistory(payload.Query, &streamResponse.Data, payload.SessionRoomID)
			if err != nil {
				fmt.Fprintf(w, "data: %s\n\n", err)
				return
			}

			response, err := json.Marshal(streamResponse)
			if err != nil {
				fmt.Fprintf(w, "data: %s\n\n", fmt.Errorf("error marshaling response: %w", err))
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", string(response))
			flusher.Flush()
			time.Sleep(10 * time.Millisecond) // Optional: control message frequency
			responseStreamBody.Body.Close()
			break
		}
		// Extract and accumulate the content
		if message, ok := event["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				accumulatedContent.WriteString(content)
			}
		}

		streamResponse := types.StreamResponse{
			Done:    false,
			Message: "Streaming in progress",
			Data: types.Output{
				Response: accumulatedContent.String(),
			},
		}

		response, err := json.Marshal(streamResponse)
		if err != nil {
			fmt.Fprintf(w, "data: %s\n\n", fmt.Errorf("error marshaling response: %w", err))
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", string(response))
		flusher.Flush()
		time.Sleep(10 * time.Millisecond) // Optional: control message frequency

	}
}

func (h *Handler) createChatroom(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	roomID, err := h.store.CreateSessionRoom(userID.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	er := h.userStore.UpdateUserSessionRooms(userID, roomID)
	if er != nil {
		utils.WriteError(w, http.StatusInternalServerError, er)
		return
	}

	responseFormat := types.APIResponse{
		Success: true,
		Message: "Session room creation successful.",
		Data:    roomID,
	}

	utils.WriteJSON(w, http.StatusCreated, responseFormat)
}

func (h *Handler) getAllChatroomIDs(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	roomIDs, err := h.store.GetAllSessionRoomID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	responseFormat := types.APIResponse{
		Success: true,
		Message: "Successfully retrieve roomIDs",
		Data:    roomIDs,
	}

	utils.WriteJSON(w, http.StatusOK, responseFormat)
}

func (h *Handler) getChatroomHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	vars := mux.Vars(r)
	roomID := vars["roomID"]
	if roomID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("roomID is required"))
		return
	}

	messages, err := h.store.LoadChatHistory(userID.Hex(), roomID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resposneFormat := types.APIResponse{
		Success: true,
		Message: "Successfully retrieve chat history",
		Data:    messages,
	}
	utils.WriteJSON(w, http.StatusOK, resposneFormat)
}

func (h *Handler) deleteChatroom(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	vars := mux.Vars(r)
	roomID := vars["roomID"]
	if roomID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("roomID is required"))
		return
	}

	er := h.store.DeleteSessionRoom(userID, roomID)
	if er != nil {
		utils.WriteError(w, http.StatusInternalServerError, er)
		return
	}
	responseFormat := types.APIResponse{
		Success: true,
		Message: "Session room successful delete.",
		Data:    nil,
	}

	utils.WriteJSON(w, http.StatusOK, responseFormat)
}
