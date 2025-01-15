live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

live/server:
	go run github.com/air-verse/air@latest \
	--build.cmd "go build -o tmp/main main.go" --build.bin "tmp/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live/tailwind:
	npx --yes tailwindcss -i ./views/css/app.css -o ./public/styles.css --minify --watch=always

live/sync_assets:
	go run github.com/air-verse/air@latest \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "public" \
	--build.include_ext "js,css"

live:
	make -j4 live/templ live/server live/tailwind live/sync_assets
