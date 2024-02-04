package post

import (
	"EWallet/internal/domain/model"
	"EWallet/internal/repository/postgresql/outsideErrors"
	"EWallet/internal/transport/http/handlers"
	"errors"
	"log/slog"
	"net/http"
)

func CreateBalance(uc handlers.DatabaseUseCase, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "CreateBalance"),
		)

		newWallet := model.NewWallet()

		walletOutput, err := uc.CreateWallet(newWallet)
		if err != nil {
			switch {
			case errors.Is(err, outsideErrors.WalletAlreadyExist):
				log.Error("uuid already exist:", err)
				http.Error(w, outsideErrors.WalletAlreadyExist.Error(), http.StatusInternalServerError)
			default:
				log.Error("internal error:", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}
		w.Header().Add("id", walletOutput.ID)
		w.Header().Add("balance", walletOutput.Balance)
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("well done!")); err != nil {
			log.Error("responseWriter error:", err)
		}
	}
}
