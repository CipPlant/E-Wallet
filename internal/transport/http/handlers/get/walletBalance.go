package get

import (
	"EWallet/internal/repository/postgresql/outsideErrors"
	"EWallet/internal/transport/dto"
	"EWallet/internal/transport/http/handlers"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func WalletBalance(uc handlers.DatabaseUseCase, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "WalletBalance"),
		)

		vars := mux.Vars(r)
		from, ok := vars["walletID"]
		if !ok {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		inputData := dto.WalletInput{
			WalletID: from,
		}

		dtoOutput, err := uc.GetWalletBalance(inputData)
		if err != nil {
			switch {
			case errors.Is(err, outsideErrors.InvalidWalletFormat):
				http.Error(w, outsideErrors.InvalidWalletFormat.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, outsideErrors.NoSuchWallet):
				http.Error(w, outsideErrors.NoSuchWallet.Error(), http.StatusNotFound)
				return
			default:
				log.Error("internal server error:", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

		}
		jsonData, err := json.Marshal(dtoOutput)
		if err != nil {
			log.Error("internal server error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(jsonData); err != nil {
			log.Error("responseWriter error:", err)
		}
	}
}
