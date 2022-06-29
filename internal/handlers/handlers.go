package handlers

import (
	"encoding/json"
	"last_lesson/internal/asyncCollector"
	"last_lesson/internal/mytypes"
	"net/http"
	"time"
)

var (
	timeSync        mytypes.Synchronizer
	PendingResponse mytypes.ResultT
	onInternal      bool
)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	if pendingAnswer() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(PendingResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func pendingAnswer() error {
	now := time.Now()

	if (now.Sub(timeSync.LastResponseTime) >= 30*time.Second) || (onInternal) {
		timeSync.NewLoop(now)
		onInternal = false

		data, err := asyncCollector.GetResult()
		if err != nil {
			onInternal = true
			return err
		}

		if !answerIsCorrect(data) {
			PendingResponse.Status = false
			PendingResponse.Error = "Error on collect data"
			PendingResponse.Data = mytypes.ResultSetT{}
		} else {
			PendingResponse.Status = true
			PendingResponse.Error = ""
			PendingResponse.Data = *data
		}
	}

	return nil
}

func answerIsCorrect(answer *mytypes.ResultSetT) bool {
	onEmptySMS := len((*answer).SMS) == 0
	onEmptyMMS := len((*answer).MMS) == 0
	onEmptyVC := len((*answer).VoiceCall) == 0
	onEmptyEmail := len((*answer).Email) == 0
	onEmptySup := len((*answer).Support) == 0
	onEmptyInc := len((*answer).Incidents) == 0

	return !onEmptySMS && !onEmptyMMS && !onEmptyVC && !onEmptyEmail && !onEmptySup && !onEmptyInc
}
