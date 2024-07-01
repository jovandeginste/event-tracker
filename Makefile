NAME=event-tracker
GIT_REF ?= $(shell git symbolic-ref HEAD)
GIT_REF_NAME ?= $(shell git branch --show-current)
GIT_REF_TYPE ?= branch
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
OUTPUT_FILE ?= tmp/main

I18N_LANGUAGES ?= nl de fr
LANG_TO_GENERATE = $(patsubst generate-translation-%,%,$@)
THEME_SCREENSHOT_WIDTH ?= 1200
THEME_SCREENSHOT_HEIGHT ?= 900

.PHONY: all clean test build screenshots

all: clean install-deps test build

install-deps:
	npm install

clean:
	rm -fv ./assets/output.css ./$(NAME)
	rm -rf ./tmp/ ./node_modules/ ./assets/dist/

dev:
	air

build: build-swagger build-dist build-tw build-server build-docker screenshots

build-server:
	go build \
		-ldflags "-X 'main.buildTime=$(BUILD_TIME)' -X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.gitRef=$(GIT_REF)' -X 'main.gitRefName=$(GIT_REF_NAME)' -X 'main.gitRefType=$(GIT_REF_TYPE)'" \
		-o $(OUTPUT_FILE) ./

build-docker:
	docker build -t $(NAME)
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg GIT_COMMIT="$(GIT_COMMIT)" \
		--build-arg GIT_REF="$(GIT_REF)" \
		--build-arg GIT_REF_NAME="$(GIT_REF_NAME)" \
		--build-arg GIT_REF_TYPE="$(GIT_REF_TYPE)" \
		.

build-swagger:
	swag init \
		--parseDependency \
		--dir ./pkg/app/,./,./vendor/gorm.io/gorm/,./vendor/github.com/codingsince1985/geo-golang/ \
		--generalInfo api_handlers.go

build-tw:
	npx tailwindcss -i ./main.css -o ./assets/output.css

clean-dist:
	rm -rf ./assets/dist/

build-dist: clean-dist
	mkdir -p ./assets/dist/images
	cp -v ./node_modules/fullcalendar/index.global.min.js ./assets/dist/fullcalendar.min.js
	cp -R ./node_modules/@fortawesome/fontawesome-free/ ./assets/dist/fontawesome/
	cp -v ./node_modules/htmx.org/dist/htmx.min.js ./assets/dist/


watch-tw:
	npx tailwindcss -i ./main.css -o ./assets/output.css --watch

generate-messages:
	xspreak -p ./translations/ -f json --template-keyword "i18n" -t "views/**/*.html"

generate-translations: $(patsubst %,generate-translation-%, $(I18N_LANGUAGES))

$(patsubst %,generate-translation-%, $(I18N_LANGUAGES)):
	xspreak merge -i translations/messages.json \
		-o translations/${LANG_TO_GENERATE}.json -l ${LANG_TO_GENERATE}
	prettier --write translations/${LANG_TO_GENERATE}.json

serve:
	$(OUTPUT_FILE)

test: test-go test-assets

test-assets:
	prettier --check .

test-go:
	go test -short -count 1 -mod vendor -covermode=atomic ./...
	golangci-lint run --allow-parallel-runners

screenshots: generate-screenshots screenshots-theme screenshots-responsive

generate-screenshots:
	K6_BROWSER_ARGS="force-dark-mode" k6 run screenshots.js

go-cover:
	go test -short -count 1 -mod vendor -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm -vf coverage.out

update-deps:
	npm update
	go get -d -t ./...
