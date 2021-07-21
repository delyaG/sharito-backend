package postgres

import (
	"backend/internal/domain"
	"backend/internal/infra/postgres/models"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type adapter struct {
	logger logrus.FieldLogger
	config *Config
	db     *sqlx.DB
}

func NewAdapter(logger logrus.FieldLogger, config *Config) (domain.Database, error) {
	a := &adapter{
		logger: logger,
		config: config,
	}

	db, err := sqlx.Open("pgx", config.ConnectionString())
	if err != nil {
		logger.Errorf("cannot open sql connection: %w", err)
		return nil, err
	}
	a.db = db

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifeTime)

	// Migrations block
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(config.MigrationsSourceURL, config.Name, driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return a, nil
}

func (a *adapter) SaveUser(user *domain.User) (int, error) {
	var id int
	if err := a.db.Get(
		&id,
		`INSERT INTO users (login, first_name, last_name, email, password_hash, salt)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id`,
		user.Login,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Salt,
	); err != nil {
		a.logger.WithError(err).Error("Error while saving user info!")
		return 0, domain.ErrInternalDatabase
	}

	return id, nil
}

func (a *adapter) GetUserByLogin(login string) (*domain.User, error) {
	var user models.User

	if err := a.db.Get(
		&user,
		`SELECT id, login, first_name, last_name, email, password_hash, salt
				FROM users
				WHERE login = $1`,
		login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.WithError(err).Error("There is no such user!")
			return nil, domain.ErrNoSuchUser
		}
		a.logger.WithError(err).Error("Error while getting user info by id!")
		return nil, domain.ErrInternalDatabase
	}

	return user.Domain(), nil
}

func (a *adapter) GetUserByID(id int) (*domain.User, error) {
	var user models.User

	if err := a.db.Get(
		&user,
		`SELECT id, login, first_name, last_name, email, password_hash, salt
				FROM users
				WHERE id = $1`,
		id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.WithError(err).Error("There is no such user!")
			return nil, domain.ErrNoSuchUser
		}
		a.logger.WithError(err).Error("Error while getting user info by id!")
		return nil, domain.ErrInternalDatabase
	}

	return user.Domain(), nil
}

func (a *adapter) SaveProduct(product *domain.Product) (int, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		a.logger.WithError(err).Error("Error while trying to begin a database transaction!")
		return 0, err
	}

	defer func(err *error) {
		if *err != nil {
			if err := tx.Rollback(); err != nil {
				a.logger.WithError(err).Error("Error while trying to rollback a database transaction!")
			}
		}
	}(&err)

	var id int

	if err := tx.Get(
		&id,
		`INSERT INTO products (owner_id, name, per_hour, description)
				VALUES ($1, $2, $3, $4)
				RETURNING id`,
		product.OwnerID,
		product.Name,
		product.PerHour,
		product.Description,
	); err != nil {
		a.logger.WithError(err).Error("Error while saving product info!")
		return 0, domain.ErrInternalDatabase
	}

	for _, v := range product.Photos {
		if _, err := tx.Exec(
			`INSERT INTO product_photos (product_id, photo)
				VALUES ($1, $2)`,
			id,
			v,
		); err != nil {
			a.logger.WithError(err).Error("Error while saving product photos!")
			return 0, domain.ErrInternalDatabase
		}
	}

	if err := tx.Commit(); err != nil {
		a.logger.WithError(err).Error("Error while trying to commit a database transaction!")
		return 0, domain.ErrInternalDatabase
	}

	return id, nil
}

func (a *adapter) GetProductByID(id int) (*domain.Product, error) {
	var product models.Product

	if err := a.db.Get(
		&product,
		`SELECT id, owner_id, name, per_hour, description
				FROM products
				WHERE id = $1`,
		id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		a.logger.WithError(err).Error("Error while getting product by id!")
		return nil, domain.ErrInternalDatabase
	}

	if err := a.db.Select(
		&product.Photos,
		`SELECT photo FROM product_photos WHERE product_id = $1`,
		id); err != nil {
		a.logger.WithError(err).Error("Error while getting photos by product_id!")
		return nil, domain.ErrInternalDatabase
	}

	return product.Domain(), nil
}

func (a *adapter) GetProductsWithPagination(limit, offset int, search string) ([]*domain.Product, int, error) {
	template := "%" + search + "%"
	var p models.Products
	if err := a.db.Select(&p,
		`SELECT id, owner_id, name, per_hour, description
				FROM products
				WHERE name ILIKE $1
				LIMIT $2 OFFSET $3`,
		template,
		limit,
		offset); err != nil {
		a.logger.WithError(err).Error("Error while getting products with pagination!")
		return nil, 0, domain.ErrInternalDatabase
	}

	for _, v := range p {
		if err := a.db.Select(&v.Photos,
			`SELECT photo FROM product_photos WHERE product_id = $1 ORDER BY created_at LIMIT 1`,
			v.ID,
		); err != nil {
			a.logger.WithError(err).Error("Error while getting main product photo!")
			return nil, 0, domain.ErrInternalDatabase
		}
	}

	var count int
	if err := a.db.Get(
		&count,
		`SELECT count(id) FROM products`,
	); err != nil {
		a.logger.WithError(err).Error("Error while getting product count!")
		return nil, 0, domain.ErrInternalDatabase
	}

	return p.Domain(), count, nil
}

func (a *adapter) RentProduct(productID, userID int, from, to time.Time) error {
	if _, err := a.db.Exec(
		`INSERT INTO orders (user_id, product_id, order_start, order_end)
				VALUES ($1, $2, $3, $4)`,
		userID,
		productID,
		from,
		to,
	); err != nil {
		a.logger.WithError(err).Error("Error while saving order!")
		return domain.ErrInternalDatabase
	}

	return nil
}

func (a *adapter) GetOrders(userID int, isMine bool) ([]*domain.Order, error) {
	var orders models.Orders

	if isMine {
		if err := a.db.Select(&orders,
			`SELECT id, user_id, product_id, order_start, order_end, 
       				extract(EPOCH FROM order_end - order_start) / 3600 AS price 
					FROM orders 
					WHERE user_id = $1`,
			userID,
		); err != nil {
			a.logger.WithError(err).Error("Error while getting orders!")
			return nil, domain.ErrInternalDatabase
		}
	} else {
		if err := a.db.Select(&orders,
			`SELECT p.id, user_id, product_id, order_start, order_end,
       				extract(EPOCH FROM order_end - order_start) / 3600 AS price
					FROM orders
					LEFT JOIN products p ON p.id = orders.product_id 
					WHERE p.owner_id = $1`,
			userID,
		); err != nil {
			a.logger.WithError(err).Error("Error while getting orders!")
			return nil, domain.ErrInternalDatabase
		}
	}

	return orders.Domain(), nil
}
