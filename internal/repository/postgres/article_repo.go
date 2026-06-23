package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type articleRepository struct{ pool *pgxpool.Pool }

func NewArticleRepository(pool *pgxpool.Pool) *articleRepository {
	return &articleRepository{pool: pool}
}

const articleSelectCols = `
    a.id, a.venue_id, a.external_id, a.location_id,
    a.placement, a.navigate, a.title, a.label, a.content,
    a.status, a.published_at, a.published_period_start, a.published_period_end,
    a.localization, a.created_at, a.updated_at`

func (r *articleRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Article, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM articles WHERE venue_id=$1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+articleSelectCols+` FROM articles a WHERE a.venue_id=$1 AND a.deleted_at IS NULL
         ORDER BY a.created_at DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.Article
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, a)
	}
	return out, total, rows.Err()
}

func (r *articleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+articleSelectCols+` FROM articles a WHERE a.id=$1 AND a.deleted_at IS NULL`, id)
	a, err := scanArticle(row)
	if err != nil {
		return nil, err
	}
	images, err := r.listImages(ctx, id)
	if err != nil {
		return nil, err
	}
	a.Images = images
	return a, nil
}

func (r *articleRepository) Create(ctx context.Context, a *domain.Article) error {
	a.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO articles
         (id,venue_id,external_id,location_id,placement,navigate,title,label,content,
          status,published_at,published_period_start,published_period_end,localization)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
         RETURNING created_at,updated_at`,
		a.ID, a.VenueID, a.ExternalID, a.LocationID, a.Placement, a.Navigate,
		a.Title, a.Label, a.Content, a.Status,
		a.PublishedAt, a.PublishedPeriodStart, a.PublishedPeriodEnd, a.Localization,
	).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *articleRepository) Update(ctx context.Context, a *domain.Article) error {
	return r.UpdateWithImages(ctx, a, domain.ArticleMediaChange{Replace: false})
}

func (r *articleRepository) UpdateWithImages(ctx context.Context, a *domain.Article, change domain.ArticleMediaChange) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	if err := updateArticleScalar(ctx, tx, a); err != nil {
		return err
	}

	var images []*domain.ArticleImage
	if change.Replace {
		if _, err := tx.Exec(ctx,
			`UPDATE article_images SET deleted_at=NOW(), updated_at=NOW()
             WHERE article_id=$1 AND deleted_at IS NULL`,
			a.ID,
		); err != nil {
			return err
		}
		images = make([]*domain.ArticleImage, 0, len(change.Keys))
		for i, key := range change.Keys {
			img := &domain.ArticleImage{
				ID:        newID(),
				ArticleID: a.ID,
				Image:     key,
				SortOrder: i,
			}
			if err := tx.QueryRow(ctx,
				`INSERT INTO article_images (id,article_id,image,sort_order)
                 VALUES ($1,$2,$3,$4) RETURNING created_at,updated_at`,
				img.ID, img.ArticleID, img.Image, img.SortOrder,
			).Scan(&img.CreatedAt, &img.UpdatedAt); err != nil {
				return err
			}
			images = append(images, img)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	committed = true
	if change.Replace {
		a.Images = images
	}
	return nil
}

func updateArticleScalar(ctx context.Context, q pgx.Tx, a *domain.Article) error {
	return q.QueryRow(ctx,
		`UPDATE articles SET
         external_id=$2,location_id=$3,placement=$4,navigate=$5,title=$6,label=$7,content=$8,
         status=$9,published_at=$10,published_period_start=$11,published_period_end=$12,localization=$13
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		a.ID, a.ExternalID, a.LocationID, a.Placement, a.Navigate,
		a.Title, a.Label, a.Content, a.Status,
		a.PublishedAt, a.PublishedPeriodStart, a.PublishedPeriodEnd, a.Localization,
	).Scan(&a.UpdatedAt)
}

func (r *articleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE articles SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *articleRepository) CreateImage(ctx context.Context, img *domain.ArticleImage) error {
	img.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO article_images (id,article_id,image,sort_order) VALUES ($1,$2,$3,$4)
         RETURNING created_at,updated_at`,
		img.ID, img.ArticleID, img.Image, img.SortOrder,
	).Scan(&img.CreatedAt, &img.UpdatedAt)
}

func (r *articleRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE article_images SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *articleRepository) listImages(ctx context.Context, articleID uuid.UUID) ([]*domain.ArticleImage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,article_id,image,sort_order,created_at,updated_at
         FROM article_images WHERE article_id=$1 AND deleted_at IS NULL ORDER BY sort_order`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.ArticleImage
	for rows.Next() {
		var img domain.ArticleImage
		if err := rows.Scan(&img.ID, &img.ArticleID, &img.Image, &img.SortOrder, &img.CreatedAt, &img.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &img)
	}
	return out, rows.Err()
}

func scanArticle(row scanner) (*domain.Article, error) {
	var a domain.Article
	if err := row.Scan(
		&a.ID, &a.VenueID, &a.ExternalID, &a.LocationID,
		&a.Placement, &a.Navigate, &a.Title, &a.Label, &a.Content,
		&a.Status, &a.PublishedAt, &a.PublishedPeriodStart, &a.PublishedPeriodEnd,
		&a.Localization, &a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &a, nil
}
