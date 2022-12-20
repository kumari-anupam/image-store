package dbhandler

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func handlerError(err error, txn *sqlx.Tx) error {
	if err != nil {
		err1 := txn.Rollback()
		if err1 != nil {
			log.Errorf("error while rollback the transaction, %v", err)

			return fmt.Errorf("error while rollback the transaction, %w", err1)
		}

		return err
	}

	err = txn.Commit()
	if err != nil {
		log.Errorf("error while transaction commit for db, %v", err)

		return fmt.Errorf("error while transaction commit,%w", err)
	}

	log.Debug("db transaction commit is done")

	return nil
}
