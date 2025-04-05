package sqlite

import (
	"context"
	"fmt"

	"github.com/artifact-space/ArtiSpace/consts/queries"
	"github.com/artifact-space/ArtiSpace/log"
)

func (s *sqliteDb) CreateDockerNamespaceAndRepositoryIfMissing(ctx context.Context, namespace string, repository string) error {
	tx, err := s.Begin()
	if err != nil {
		log.Logger().Error().Err(err).Msg("unable to create db transaction")
		return err
	}
	defer func() {
		if err != nil {
			log.Logger().Warn().Msgf("transaction is going to be rolled-back due to errors")
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var namespaceId string

	err = tx.QueryRowContext(ctx, queries.SqliteGetDockerNamespaceID, namespace).Scan(&namespaceId)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to retrieve docker namespace")
		return err
	}

	if namespaceId == "" {
		log.Logger().Info().Msgf("docker namespace %s does not exists", namespace)

		// create docker namespace
		err = tx.QueryRowContext(ctx, queries.SqliteCreateDockerNamespace, namespace).Scan(&namespaceId)
		if err != nil {
			log.Logger().Error().Err(err).Msgf("unable to create docker namespace: %s", namespace)
			return err
		}
		if namespaceId == "" {
			log.Logger().Error().Msgf("docker namespace %s was not created", namespace)
			return fmt.Errorf("unable to create docker namespace: %s", namespace)
		}
	}

	var dockerRepositoryId string

	err = tx.QueryRowContext(ctx, queries.SqliteGetDockerRepositoryID, namespaceId, repository).Scan(&dockerRepositoryId)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to retrive docker repository: %s:%s", namespace, repository)
		return err
	}

	if dockerRepositoryId == "" {
		log.Logger().Info().Msgf("docker respository %s:%s doesn not exists.", namespace, repository)

		res, err := tx.ExecContext(ctx, queries.SqliteCreateDockerRepository, repository, namespaceId)
		if err != nil {
			log.Logger().Error().Err(err).Msgf("unable to create docker repository: %s:%s", namespace, repository)
			return err
		}

		if insertCount, err := res.RowsAffected(); err != nil || insertCount != 1 {
			log.Logger().Error().Err(err).Msgf("docker repository creation was not successful. %s:%s", namespace, repository)
			if err != nil {
				return err
			}
			return fmt.Errorf("no new docker repository was added")
		}
	}

	return nil
}
