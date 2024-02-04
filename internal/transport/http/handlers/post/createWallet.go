package post

import (
	"EWallet/internal/domain/model"
	"EWallet/internal/transport/http/handlers"
	"log/slog"
	"net/http"
)

func CreateBalance(uc handlers.DatabaseUseCase, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := logger.With(
			slog.String("component", ""),
		)

		newWallet := model.NewWallet()

		walletOutput, err := uc.CreateWallet(newWallet)

		if err != nil {
			switch {
			// for handling future errors
			default:
				log.Error("error FindHandler.")
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Add("id", walletOutput.ID)
		w.Header().Add("balance", walletOutput.Balance)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("well done!"))

	}
}
