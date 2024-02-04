package post

import (
	"EWallet/internal/repository/postgresql/outsideErrors"
	"EWallet/internal/transport/dto"
	"EWallet/internal/transport/http/handlers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func SendMoney(uc handlers.DatabaseUseCase, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		from, ok := vars["walletID"]

		if !ok {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		var dtoInput dto.WalletInput

		dtoInput.WalletID = from

		err := json.NewDecoder(r.Body).Decode(&dtoInput)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = uc.SendMoney(dtoInput)

		if err != nil {
			switch {
			case errors.Is(err, outsideErrors.InvalidWalletFormat):
				http.Error(w, outsideErrors.InvalidWalletFormat.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, outsideErrors.InvalidMoneyFormat):
				http.Error(w, outsideErrors.InvalidMoneyFormat.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, outsideErrors.NotEnoughMoney):
				http.Error(w, outsideErrors.NotEnoughMoney.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, outsideErrors.NoSuchOutgoingWallet):
				http.Error(w, outsideErrors.NoSuchOutgoingWallet.Error(), http.StatusNotFound)
				return
			case errors.Is(err, outsideErrors.NoSuchDestinationWallet):
				fmt.Println("should be there")
				http.Error(w, outsideErrors.NoSuchDestinationWallet.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
				return

			}
		}
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("well done!"))
	}
}
