package cmd

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/platform/database"
)

var seedReset bool

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with test data",
	Long:  `Inserts a full set of test data covering all tables. Safe to run multiple times (skips if seed venue already exists). Use --reset to wipe and re-seed.`,
	RunE:  runSeed,
}

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().BoolVar(&seedReset, "reset", false, "Delete existing seed data before inserting")
}

func runSeed(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pool, err := database.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	ctx := context.Background()

	if seedReset {
		if err := resetSeed(ctx, pool); err != nil {
			return err
		}
	}

	var exists bool
	if err := pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM venues WHERE public_key='seed-pub-key-000000000000000000000000000000' AND deleted_at IS NULL)`,
	).Scan(&exists); err != nil {
		return err
	}
	if exists {
		fmt.Println("seed data already present — skipping (use --reset to re-seed)")
		return nil
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// ── Customer ──────────────────────────────────────────────────────────────
	var customerID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO customers (name, email, phone, address, image, url, description, is_active)
		VALUES ('Digitran Asia', 'admin@digitran.asia', '+81-3-5555-0100',
		        '1-2-3 Ariake, Koto City, Tokyo 135-0063, Japan',
		        'https://example.com/logo.png',
		        'https://digitran.asia', 'Indoor mapping solutions provider.', true)
		RETURNING id`).Scan(&customerID); err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}
	fmt.Printf("customer:        %s\n", customerID)

	// ── Venue ─────────────────────────────────────────────────────────────────
	var venueID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO venues (
			customer_id, name, external_id, type,
			public_key, private_key,
			address, city, state, country, postal,
			lat, lng, timezone, telephone, description,
			theme, app_configs, app_domains,
			seo_title, seo_description, seo_keywords,
			head_tag, start_at, end_at
		) VALUES (
			$1, 'Foodex Mar 2026', 'foodex_mar_2026', 5,
			'seed-pub-key-000000000000000000000000000000',
			'seed-priv-key-00000000000000000000000000000000000000000000000000',
			'3-11-1 Ariake, Koto City', 'Tokyo', 'Tokyo', 'Japan', '135-0063',
			35.6269, 139.7956, 'Asia/Tokyo', '+81-3-5530-1111',
			'Annual international food and beverage trade show.',
			'{"primaryColor":"#E63946","secondaryColor":"#457B9D"}'::jsonb,
			'{"mapStyle":"default","searchEnabled":true}'::jsonb,
			'["foodex2026.example.com"]'::jsonb,
			'Foodex Mar 2026 — Tokyo Big Sight',
			'International food and beverage exhibition in Tokyo.',
			'foodex,exhibition,tokyo,food',
			'<!-- GA: G-SEED2026 -->',
			NOW(), NOW() + INTERVAL '5 days'
		) RETURNING id`, customerID).Scan(&venueID); err != nil {
		return fmt.Errorf("insert venue: %w", err)
	}
	fmt.Printf("venue:           %s\n", venueID)

	// ── Users ─────────────────────────────────────────────────────────────────
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	var editorID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, is_active)
		VALUES ('editor@digitran.asia', $1, 'Hung', 'Nguyen', '+81-90-1234-5678', true)
		RETURNING id`, string(hash)).Scan(&editorID); err != nil {
		return fmt.Errorf("insert editor user: %w", err)
	}
	var viewerID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, first_name, last_name, is_active)
		VALUES ('viewer@digitran.asia', $1, 'Alice', 'Viewer', true)
		RETURNING id`, string(hash)).Scan(&viewerID); err != nil {
		return fmt.Errorf("insert viewer user: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO venue_user_roles (venue_id, user_id, role)
		VALUES ($1, $2, 'editor'), ($1, $3, 'viewer')`,
		venueID, editorID, viewerID); err != nil {
		return fmt.Errorf("insert venue roles: %w", err)
	}
	fmt.Printf("users:           editor=%s  viewer=%s  (password123)\n", editorID, viewerID)

	// ── Map group + perspective + 2 levels ───────────────────────────────────
	var mapGroupID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO map_groups (venue_id, name, short_name, type, sort_index)
		VALUES ($1, 'Tokyo Big Sight East', 'East Hall', 'hall', 0)
		RETURNING id`, venueID).Scan(&mapGroupID); err != nil {
		return fmt.Errorf("insert map_group: %w", err)
	}

	var perspID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO perspectives (
			name, camera_zoom, camera_type,
			camera_max_zoom, camera_min_zoom,
			camera_target_center_lng, camera_target_center_lat,
			camera_target_zoom, camera_target_bearing, camera_target_pitch
		) VALUES ('East Hall Default', 17.5, 0, 20.0, 14.0, 139.7956, 35.6269, 17.5, 0, 0)
		RETURNING id`).Scan(&perspID); err != nil {
		return fmt.Errorf("insert perspective: %w", err)
	}

	var level1ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO levels (
			venue_id, map_group_id, perspective_id, name, short_name, external_id,
			type, latitude, longitude, bearing, width, height, scale,
			level_width, level_height, elevation, is_published
		) VALUES ($1,$2,$3, 'Hall 1 (Ground)', 'H1-GF', 'EAST-H1-GF',
			1, 35.6269, 139.7956, 0.0, 1024, 768, 0.05,
			200.0, 150.0, 0, true)
		RETURNING id`, venueID, mapGroupID, perspID).Scan(&level1ID); err != nil {
		return fmt.Errorf("insert level 1: %w", err)
	}
	var level2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO levels (
			venue_id, map_group_id, name, short_name, external_id,
			type, latitude, longitude, bearing, width, height, scale,
			level_width, level_height, elevation, is_published
		) VALUES ($1,$2, 'Hall 2 (Ground)', 'H2-GF', 'EAST-H2-GF',
			1, 35.6270, 139.7958, 0.0, 1024, 768, 0.05,
			200.0, 150.0, 0, true)
		RETURNING id`, venueID, mapGroupID).Scan(&level2ID); err != nil {
		return fmt.Errorf("insert level 2: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO geo_references (level_id, control_x, control_y, target_x, target_y)
		VALUES ($1, 0, 0, 139.7950, 35.6265), ($1, 1024, 768, 139.7962, 35.6273)`,
		level1ID); err != nil {
		return fmt.Errorf("insert geo_references: %w", err)
	}
	fmt.Printf("levels:          %s  %s\n", level1ID, level2ID)

	// ── Location categories ───────────────────────────────────────────────────
	var catFoodID, catElecID, catFeaturedID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO location_categories (venue_id, name, short_name, color, icon, icon_default, sort_index, visible, source, description)
		VALUES ($1, 'Food & Beverage', 'F&B', '#E63946', 'restaurant', 'food', 0, true, 'internal', 'Food, drinks and catering booths')
		RETURNING id`, venueID).Scan(&catFoodID); err != nil {
		return fmt.Errorf("insert category food: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO location_categories (venue_id, name, short_name, color, icon, icon_default, sort_index, visible, source, description)
		VALUES ($1, 'Electronics', 'Elec', '#457B9D', 'bolt', 'electronics', 1, true, 'internal', 'Electronics and tech exhibitors')
		RETURNING id`, venueID).Scan(&catElecID); err != nil {
		return fmt.Errorf("insert category electronics: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO location_categories (venue_id, name, short_name, color, icon, sort_index, visible, source, description)
		VALUES ($1, 'Premium Zone', 'PZ', '#2A9D8F', 'star', 2, true, 'external', 'VIP and premium exhibitor zone')
		RETURNING id`, venueID).Scan(&catFeaturedID); err != nil {
		return fmt.Errorf("insert category featured: %w", err)
	}
	fmt.Printf("location_cats:   food=%s  elec=%s  featured=%s\n", catFoodID, catElecID, catFeaturedID)

	// ── Locations ─────────────────────────────────────────────────────────────
	var loc1ID, loc2ID, loc3ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO locations (
			venue_id, level_id, main_category_id, external_id,
			common_name, common_short_name, common_description,
			common_color, common_location_type, common_location_sub_type,
			common_latitude, common_longitude,
			common_contact_email, common_contact_phone,
			common_social_website, booth_number, is_top_location
		) VALUES (
			$1, $2, $3, 'LOC-F001',
			'Tokyo Ramen Co.', 'Ramen', 'Authentic Tokyo-style ramen booth.',
			'#E63946', 1, 0,
			35.62695, 139.79565,
			'ramen@tokyoramen.jp', '+81-3-1111-2222',
			'https://tokyoramen.jp', 'A-101', true
		) RETURNING id`, venueID, level1ID, catFoodID).Scan(&loc1ID); err != nil {
		return fmt.Errorf("insert location 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO locations (
			venue_id, level_id, main_category_id, external_id,
			common_name, common_short_name, common_description,
			common_color, common_location_type, common_location_sub_type,
			common_latitude, common_longitude,
			common_contact_email, booth_number, is_top_location
		) VALUES (
			$1, $2, $3, 'LOC-E001',
			'FutureTech Solutions', 'FutureTech', 'Next-gen food processing equipment.',
			'#457B9D', 1, 0,
			35.62700, 139.79570,
			'info@futuretech.jp', 'B-201', false
		) RETURNING id`, venueID, level1ID, catElecID).Scan(&loc2ID); err != nil {
		return fmt.Errorf("insert location 2: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO locations (
			venue_id, level_id, main_category_id, external_id,
			common_name, common_description,
			common_location_type, common_latitude, common_longitude,
			booth_number
		) VALUES (
			$1, $2, $3, 'LOC-P001',
			'Premium Lounge', 'VIP lounge for premium exhibitors.',
			2, 35.62710, 139.79580, 'VIP-01'
		) RETURNING id`, venueID, level2ID, catFeaturedID).Scan(&loc3ID); err != nil {
		return fmt.Errorf("insert location 3: %w", err)
	}
	fmt.Printf("locations:       %s  %s  %s\n", loc1ID, loc2ID, loc3ID)

	// ── Product categories + products ─────────────────────────────────────────
	var prodCatID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO product_categories (venue_id, name, source)
		VALUES ($1, 'Beverages', 'internal')
		RETURNING id`, venueID).Scan(&prodCatID); err != nil {
		return fmt.Errorf("insert product_category: %w", err)
	}
	var prod1ID, prod2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO products (venue_id, location_id, main_category_id, name, price, size, description, source)
		VALUES ($1,$2,$3, 'Craft Beer Sampler', '¥1500', '330ml x3',
			'A selection of three artisan craft beers from local breweries.', 'internal')
		RETURNING id`, venueID, loc1ID, prodCatID).Scan(&prod1ID); err != nil {
		return fmt.Errorf("insert product 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO products (venue_id, location_id, main_category_id, name, price, size, description, source)
		VALUES ($1,$2,$3, 'Premium Sake Set', '¥3000', '180ml x2',
			'Premium aged sake from the Niigata region.', 'internal')
		RETURNING id`, venueID, loc1ID, prodCatID).Scan(&prod2ID); err != nil {
		return fmt.Errorf("insert product 2: %w", err)
	}
	for _, pid := range []string{prod1ID, prod2ID} {
		if _, err = tx.Exec(ctx,
			`INSERT INTO product_category_links (product_id, category_id) VALUES ($1,$2)`,
			pid, prodCatID); err != nil {
			return fmt.Errorf("insert product_category_link: %w", err)
		}
	}
	fmt.Printf("products:        %s  %s\n", prod1ID, prod2ID)

	// ── Event types + events ──────────────────────────────────────────────────
	var eventTypeID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO event_types (venue_id, name)
		VALUES ($1, 'Seminar')
		RETURNING id`, venueID).Scan(&eventTypeID); err != nil {
		return fmt.Errorf("insert event_type: %w", err)
	}
	var event1ID, event2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO events (venue_id, type_id, title, description, banner_image,
			start_time, end_time, show_start_time, show_end_time, content_detail, content_url)
		VALUES ($1,$2,
			'Foodex Opening Ceremony', 'Grand opening ceremony with keynote speeches.',
			'https://example.com/banner/opening.jpg',
			NOW() + INTERVAL '1 hour', NOW() + INTERVAL '3 hours',
			NOW(), NOW() + INTERVAL '4 hours',
			'Join us for the grand opening of Foodex 2026.',
			'https://foodex2026.example.com/opening')
		RETURNING id`, venueID, eventTypeID).Scan(&event1ID); err != nil {
		return fmt.Errorf("insert event 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO events (venue_id, type_id, title, description,
			start_time, end_time, content_detail)
		VALUES ($1,$2,
			'Future of Food Tech Seminar', 'Panel discussion on food technology innovation.',
			NOW() + INTERVAL '5 hours', NOW() + INTERVAL '7 hours',
			'Industry leaders discuss the next decade of food tech.')
		RETURNING id`, venueID, eventTypeID).Scan(&event2ID); err != nil {
		return fmt.Errorf("insert event 2: %w", err)
	}
	if _, err = tx.Exec(ctx,
		`INSERT INTO event_location_links (event_id, location_id) VALUES ($1,$2)`,
		event1ID, loc1ID); err != nil {
		return fmt.Errorf("insert event_location_link: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO event_images (event_id, image)
		VALUES ($1, 'https://example.com/events/opening1.jpg'),
		       ($1, 'https://example.com/events/opening2.jpg')`, event1ID); err != nil {
		return fmt.Errorf("insert event_images: %w", err)
	}
	fmt.Printf("events:          %s  %s\n", event1ID, event2ID)

	// ── Notifications ─────────────────────────────────────────────────────────
	var notif1ID, notif2ID string
	// status=1(sent) type=1(normal) send_status=1(success) send_type=3(immediate)
	if err = tx.QueryRow(ctx, `
		INSERT INTO notifications (venue_id, title, content, topic, status, type, send_status, send_type, published_at, created_by)
		VALUES ($1, 'Welcome to Foodex 2026',
			'Explore 2,500+ exhibitors across 6 halls. Get your map now!',
			'general', 1, 1, 1, 3, NOW(), $2)
		RETURNING id`, venueID, editorID).Scan(&notif1ID); err != nil {
		return fmt.Errorf("insert notification 1: %w", err)
	}
	// status=2(unsent) send_type=2(scheduled)
	if err = tx.QueryRow(ctx, `
		INSERT INTO notifications (venue_id, title, content, status, type, send_status, send_type, scheduled_at, created_by)
		VALUES ($1, 'Opening Ceremony Starting Soon',
			'The opening ceremony begins in 30 minutes. Head to Hall 1!',
			2, 1, 0, 2, NOW() + INTERVAL '30 minutes', $2)
		RETURNING id`, venueID, editorID).Scan(&notif2ID); err != nil {
		return fmt.Errorf("insert notification 2: %w", err)
	}
	fmt.Printf("notifications:   %s  %s\n", notif1ID, notif2ID)

	// ── Surveys + questions + options ─────────────────────────────────────────
	var survey1ID, survey2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO surveys (venue_id, title, content, status, publish_type, source, app, is_forced, created_by)
		VALUES ($1, 'Visitor Satisfaction Survey',
			'Help us improve your experience at Foodex 2026.',
			2, 3, 1, 'all', false, $2)
		RETURNING id`, venueID, editorID).Scan(&survey1ID); err != nil {
		return fmt.Errorf("insert survey 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO surveys (venue_id, title, content, status, publish_type, source, app, is_forced, created_by)
		VALUES ($1, 'Exhibitor Quality Check',
			'Rate the exhibitors you visited today.',
			2, 2, 1, 'all', false, $2)
		RETURNING id`, venueID, editorID).Scan(&survey2ID); err != nil {
		return fmt.Errorf("insert survey 2: %w", err)
	}
	type question struct {
		surveyID string
		num      int
		qtype    string
		text     string
		required bool
		opts     []string
	}
	questions := []question{
		{survey1ID, 1, "single_choice", "How would you rate your overall experience?", true,
			[]string{"Excellent", "Good", "Fair", "Poor"}},
		{survey1ID, 2, "paragraph", "What did you enjoy most about this event?", false, nil},
		{survey1ID, 3, "multiple_choice", "Which areas did you visit?", false,
			[]string{"Hall 1", "Hall 2", "Conference Area", "Food Court"}},
		{survey2ID, 1, "single_choice", "How would you rate the exhibitor quality?", true,
			[]string{"5 stars", "4 stars", "3 stars", "2 stars", "1 star"}},
	}
	for _, q := range questions {
		var qID string
		if err = tx.QueryRow(ctx, `
			INSERT INTO questions (survey_id, question_number, question_type, question_text, is_required)
			VALUES ($1,$2,$3,$4,$5) RETURNING id`,
			q.surveyID, q.num, q.qtype, q.text, q.required).Scan(&qID); err != nil {
			return fmt.Errorf("insert question: %w", err)
		}
		for i, opt := range q.opts {
			if _, err = tx.Exec(ctx,
				`INSERT INTO options (question_id, option_number, option_text) VALUES ($1,$2,$3)`,
				qID, i+1, opt); err != nil {
				return fmt.Errorf("insert option: %w", err)
			}
		}
	}
	fmt.Printf("surveys:         %s  %s\n", survey1ID, survey2ID)

	// ── Coupons ───────────────────────────────────────────────────────────────
	var coupon1ID, coupon2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO coupons (venue_id, external_id, coupon_name, coupon_code, status, issued_at, expired_at)
		VALUES ($1, 'CPN-2026-001', '10% Off Premium Sake', 'SAKE10', 'active',
			NOW(), NOW() + INTERVAL '5 days')
		RETURNING id`, venueID).Scan(&coupon1ID); err != nil {
		return fmt.Errorf("insert coupon 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO coupons (venue_id, external_id, coupon_name, coupon_code, status, issued_at, expired_at)
		VALUES ($1, 'CPN-2026-002', 'Free Beer Sample', 'BEERFREE', 'active',
			NOW(), NOW() + INTERVAL '3 days')
		RETURNING id`, venueID).Scan(&coupon2ID); err != nil {
		return fmt.Errorf("insert coupon 2: %w", err)
	}
	fmt.Printf("coupons:         %s  %s\n", coupon1ID, coupon2ID)

	// ── Advertisements ────────────────────────────────────────────────────────
	var ad1ID, ad2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO advertisements (
			venue_id, location_id, type, status, navigate,
			content_image, content_cta_url, placement,
			size_width, size_height, display_duration,
			published_at, start_at, end_at
		) VALUES (
			$1,$2, 'dialog', 'published', 'location',
			NULL, 'https://foodex2026.example.com/sake',
			'home_screen', 800, 600, 10,
			NOW(), NOW(), NOW() + INTERVAL '5 days'
		) RETURNING id`, venueID, loc1ID).Scan(&ad1ID); err != nil {
		return fmt.Errorf("insert ad 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO advertisements (
			venue_id, type, status, placement,
			content_image, size_width, size_height,
			start_at, end_at
		) VALUES (
			$1, 'banner', 'published', 'map_screen',
			NULL, 1200, 200,
			NOW(), NOW() + INTERVAL '5 days'
		) RETURNING id`, venueID).Scan(&ad2ID); err != nil {
		return fmt.Errorf("insert ad 2: %w", err)
	}
	fmt.Printf("ads:             %s  %s\n", ad1ID, ad2ID)

	// ── Articles ──────────────────────────────────────────────────────────────
	var art1ID, art2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO articles (
			venue_id, location_id, title, label, content,
			status, placement, published_at,
			published_period_start, published_period_end
		) VALUES (
			$1,$2,
			'Top 5 Must-Visit Booths at Foodex 2026',
			'Editor''s Pick',
			'<p>From award-winning sake breweries to innovative food tech startups...</p>',
			'published', 'article', NOW(),
			CURRENT_DATE, CURRENT_DATE + 5
		) RETURNING id`, venueID, loc1ID).Scan(&art1ID); err != nil {
		return fmt.Errorf("insert article 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO articles (venue_id, title, content, status, placement, published_at)
		VALUES ($1,
			'Venue Map & Navigation Guide',
			'<p>Use the interactive map to find your way around Tokyo Big Sight.</p>',
			'published', 'article', NOW())
		RETURNING id`, venueID).Scan(&art2ID); err != nil {
		return fmt.Errorf("insert article 2: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO article_images (article_id, image, sort_order)
		VALUES ($1,'https://example.com/articles/foodex-booth1.jpg',0),
		       ($1,'https://example.com/articles/foodex-booth2.jpg',1)`, art1ID); err != nil {
		return fmt.Errorf("insert article_images: %w", err)
	}
	fmt.Printf("articles:        %s  %s\n", art1ID, art2ID)

	// ── Beacons ───────────────────────────────────────────────────────────────
	var beacon1ID, beacon2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO beacons (venue_id, level_id, name, hw_id, uuid_val, mac, major, minor, radius, battery, position_x, position_y, is_enable)
		VALUES ($1,$2, 'Hall 1 Entrance', 'HW-E001', 'E2C56DB5-DFFB-48D2-B060-D0F5A71096E0',
			'AA:BB:CC:DD:EE:01', 1, 1, 5, 85, 100.0, 200.0, true)
		RETURNING id`, venueID, level1ID).Scan(&beacon1ID); err != nil {
		return fmt.Errorf("insert beacon 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO beacons (venue_id, level_id, name, hw_id, uuid_val, mac, major, minor, radius, battery, position_x, position_y, is_enable)
		VALUES ($1,$2, 'Hall 2 Entrance', 'HW-E002', 'E2C56DB5-DFFB-48D2-B060-D0F5A71096E1',
			'AA:BB:CC:DD:EE:02', 1, 2, 5, 72, 500.0, 200.0, true)
		RETURNING id`, venueID, level2ID).Scan(&beacon2ID); err != nil {
		return fmt.Errorf("insert beacon 2: %w", err)
	}
	fmt.Printf("beacons:         %s  %s\n", beacon1ID, beacon2ID)

	// ── Connections ───────────────────────────────────────────────────────────
	var connID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO connections (venue_id, external_id, name, type, x, y, state, status, accessible, active)
		VALUES ($1, 'CONN-001', 'Main Escalator', 1, 512.0, 384.0, 1, 1, true, true)
		RETURNING id`, venueID).Scan(&connID); err != nil {
		return fmt.Errorf("insert connection: %w", err)
	}
	for _, lid := range []string{level1ID, level2ID} {
		if _, err = tx.Exec(ctx,
			`INSERT INTO connection_levels (connection_id, level_id, active) VALUES ($1,$2,true)`,
			connID, lid); err != nil {
			return fmt.Errorf("insert connection_level: %w", err)
		}
	}
	fmt.Printf("connection:      %s\n", connID)

	// ── Tags + entity_tags ────────────────────────────────────────────────────
	var tag1ID, tag2ID string
	if err = tx.QueryRow(ctx, `INSERT INTO tags (name) VALUES ('sake') RETURNING id`).Scan(&tag1ID); err != nil {
		return fmt.Errorf("insert tag 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `INSERT INTO tags (name) VALUES ('food-tech') RETURNING id`).Scan(&tag2ID); err != nil {
		return fmt.Errorf("insert tag 2: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO entity_tags (tag_id, entity_type, entity_id)
		VALUES ($1,'location',$2::uuid), ($3,'location',$2::uuid)`,
		tag1ID, loc1ID, tag2ID); err != nil {
		return fmt.Errorf("insert entity_tags: %w", err)
	}
	fmt.Printf("tags:            %s  %s\n", tag1ID, tag2ID)

	// ── Videos ────────────────────────────────────────────────────────────────
	var videoID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO videos (venue_id, title, description, url, thumbnail, duration, status, published_at)
		VALUES ($1, 'Foodex 2026 Highlights',
			'Recap of the best moments from Foodex 2026.',
			'https://example.com/videos/foodex2026-highlights.mp4',
			'https://example.com/videos/foodex2026-thumb.jpg',
			180, 'published', NOW())
		RETURNING id`, venueID).Scan(&videoID); err != nil {
		return fmt.Errorf("insert video: %w", err)
	}
	fmt.Printf("video:           %s\n", videoID)

	// ── App users ─────────────────────────────────────────────────────────────
	var appUser1ID, appUser2ID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO app_users (
			venue_id, external_id, source, type,
			first_name, last_name, first_name_en, last_name_en,
			email, phone, company_name, department, position,
			visitor_type, is_consented, app
		) VALUES (
			$1, 'EXT-APP-001', 'app', 0,
			'Hanako', 'Yamada', 'Hanako', 'Yamada',
			'hanako@example.jp', '+81-90-5555-0001', 'Japan Foods Inc.', 'Marketing', 'Manager',
			1, true, 'foodex'
		) RETURNING id`, venueID).Scan(&appUser1ID); err != nil {
		return fmt.Errorf("insert app_user 1: %w", err)
	}
	if err = tx.QueryRow(ctx, `
		INSERT INTO app_users (
			venue_id, external_id, source, type,
			first_name, last_name,
			email, company_name, visitor_type, is_consented, app
		) VALUES (
			$1, 'EXT-APP-002', 'import', 1,
			'John', 'Smith',
			'john@globalfoods.com', 'Global Foods Ltd.', 2, true, 'foodex'
		) RETURNING id`, venueID).Scan(&appUser2ID); err != nil {
		return fmt.Errorf("insert app_user 2: %w", err)
	}
	fmt.Printf("app_users:       %s  %s\n", appUser1ID, appUser2ID)

	// ── Assets ────────────────────────────────────────────────────────────────
	var assetID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO assets (venue_id, name, key, content_type, size_bytes, url, created_by)
		VALUES ($1, 'Venue Floor Plan', 'venues/seed/floor-plan.pdf',
			'application/pdf', 2048000,
			'https://example.com/assets/floor-plan.pdf', $2)
		RETURNING id`, venueID, editorID).Scan(&assetID); err != nil {
		return fmt.Errorf("insert asset: %w", err)
	}
	fmt.Printf("asset:           %s\n", assetID)

	// ── Level type ────────────────────────────────────────────────────────────
	var levelTypeID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO level_types (venue_id, name, icon)
		VALUES ($1, 'Exhibition Hall', 'warehouse')
		RETURNING id`, venueID).Scan(&levelTypeID); err != nil {
		return fmt.Errorf("insert level_type: %w", err)
	}
	fmt.Printf("level_type:      %s\n", levelTypeID)

	// ── Theme ─────────────────────────────────────────────────────────────────
	var themeID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO themes (venue_id, scope, name, data, storage_path)
		VALUES (
			$1,
			'custom',
			'Foodex Brand',
			'{"primary_color":"#E63946","secondary_color":"#457B9D"}'::jsonb,
			''
		)
		RETURNING id`, venueID).Scan(&themeID); err != nil {
		return fmt.Errorf("insert theme: %w", err)
	}
	fmt.Printf("theme:           %s\n", themeID)

	// ── Snapshot + level bundle ───────────────────────────────────────────────
	var snapshotID string
	if err = tx.QueryRow(ctx, `
		INSERT INTO snapshots (venue_id, state, method, created_by, publish_at)
		VALUES ($1, 2, 1, $2, NOW())
		RETURNING id`, venueID, editorID).Scan(&snapshotID); err != nil {
		return fmt.Errorf("insert snapshot: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO level_bundles (snapshot_id, venue_id, level_id, state)
		VALUES ($1,$2,$3,2)`, snapshotID, venueID, level1ID); err != nil {
		return fmt.Errorf("insert level_bundle: %w", err)
	}
	fmt.Printf("snapshot:        %s\n", snapshotID)

	// ── Analytics ─────────────────────────────────────────────────────────────
	if _, err = tx.Exec(ctx, `
		INSERT INTO event_logs (venue_id, name, params, device_id, user_agent, ip_address)
		VALUES
			($1, 'view', '{"screen":"map"}', 'device-001', 'Mozilla/5.0 (iPhone)', '203.0.113.1'),
			($1, 'search', '{"term":"ramen"}', 'device-001', 'Mozilla/5.0 (iPhone)', '203.0.113.1'),
			($1, 'view', '{"screen":"location","id":"`+loc1ID+`"}', 'device-002', 'Mozilla/5.0 (Android)', '203.0.113.2')`,
		venueID); err != nil {
		return fmt.Errorf("insert event_logs: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO search_queries (venue_id, search_term, search_count, last_searched, is_promoted, status)
		VALUES
			($1, 'ramen', 42, NOW(), false, 'published'),
			($1, 'sake',  28, NOW(), true,  'published'),
			($1, 'beer',  15, NOW(), false, 'published')`,
		venueID); err != nil {
		return fmt.Errorf("insert search_queries: %w", err)
	}
	fmt.Printf("analytics:       3 event_logs, 3 search_queries\n")

	// ── Languages ─────────────────────────────────────────────────────────────
	var langTableExists bool
	_ = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='languages')`,
	).Scan(&langTableExists)
	if langTableExists {
		if _, err = tx.Exec(ctx, `
			INSERT INTO languages (venue_id, code, name, is_default)
			VALUES ($1,'en','English',true), ($1,'ja','Japanese',false)`, venueID); err != nil {
			return fmt.Errorf("insert languages: %w", err)
		}
		fmt.Printf("languages:       en (default), ja\n")
	} else {
		fmt.Printf("languages:       skipped (run make migrate-up first)\n")
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	fmt.Println("\nSeed complete.")
	fmt.Printf("  public_key:  seed-pub-key-000000000000000000000000000000\n")
	fmt.Printf("  private_key: seed-priv-key-00000000000000000000000000000000000000000000000000\n")
	fmt.Printf("  editor:      editor@digitran.asia / password123\n")
	fmt.Printf("  viewer:      viewer@digitran.asia / password123\n")
	return nil
}

func resetSeed(ctx context.Context, pool *pgxpool.Pool) error {
	// Find the seed venue ID first.
	var venueID string
	err := pool.QueryRow(ctx,
		`SELECT id FROM venues WHERE public_key='seed-pub-key-000000000000000000000000000000'`,
	).Scan(&venueID)
	if err != nil {
		// Nothing to reset.
		return nil
	}

	// Delete tables that reference venues without ON DELETE CASCADE.
	for _, stmt := range []string{
		`DELETE FROM venue_user_roles WHERE venue_id=$1`,
		`DELETE FROM venue_invitations WHERE venue_id=$1`,
		`DELETE FROM location_categories WHERE venue_id=$1`,
		`DELETE FROM venue_amenities WHERE venue_id=$1`,
		`DELETE FROM map_groups WHERE venue_id=$1`,
	} {
		if _, err := pool.Exec(ctx, stmt, venueID); err != nil {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	if _, err := pool.Exec(ctx,
		`DELETE FROM users WHERE email IN ('editor@digitran.asia','viewer@digitran.asia')`,
	); err != nil {
		return fmt.Errorf("delete users: %w", err)
	}

	// Delete the venue — cascades to all other venue-scoped tables.
	if _, err := pool.Exec(ctx,
		`DELETE FROM venues WHERE id=$1`, venueID,
	); err != nil {
		return fmt.Errorf("delete venue: %w", err)
	}

	// Clean up the seed customer.
	if _, err := pool.Exec(ctx,
		`DELETE FROM customers WHERE email='admin@digitran.asia'`,
	); err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}

	fmt.Println("reset: existing seed data removed")
	return nil
}
