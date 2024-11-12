DOWN:
	goose -dir db/migrations postgres "user=youruser password=yourpassword dbname=yourdb sslmode=disable" down