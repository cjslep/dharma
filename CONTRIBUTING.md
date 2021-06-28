# Contributing

Dharma is built with this technological world-salad:

* Golang, the overworked server workhorse
* HTML, because native apps are something someone else can build
* SCSS, Saner CSS
* Vuejs, to deliver a slightly less displeasurable user experience

The workflow philosophy is "as simple as possible". Namely, no bloated frontend
development tools are needed. The general process is:

1. `go generate ./...` to automatically generate the locale bindings, turn SCSS
   into CSS, and pack assets into the binary.
2. `go fmt ./...` to make sure the code adheres to the Biblical Golang style.
3. `go build` for production-grade goodness or `go build -tags dev` for garbage
   only someone as deranged as I can stomach.

The deployment philosophy is "as simple as possible". Namely, no mess of zipped
nor archive files. The goal is to have a single binary, and keep it that way.

This application relies on `github.com/go-fed/apcore` as a server framework for
an ActivityPub federating application. That library always is in need of further
improvements, based on the real-world needs of this application.

## File Layout

### Frontend Development

For any new webpage, the following are needed:

* `assets/src/templates/<foo.tmpl>` for the golang-style HTML templates. Note
  that the Vue delimiter is `${}`.
* `assets/src/scss/<foo.scss>` for composing other SCSS files. It keeps pages
  trimmed so it does not have to import the whole applications' CSS.
* `assets/src/js/<foo.js>` for composing other JS files. It also keeps pages
  trimmed so it does not have to import the whole applications' Javascript.

This arrangement in the age of ruthless caching seems overkill, but it is
desirable as it is expected dharma will have extremely specialized and very
heavy pages that are rarely visited for 90% of usecases. Think industry, mining,
logistics.

The web-service side of these pages is under `internal/api/<domain>/<foo.go>`,
which then is responsible for serving the page appropriately.

### Backend Development

At a high level:

* `esi/` contains all the openapi generated code for the application to be an
  ESI client.
* `gen/` contains specialized standalone programs for code-generation time.
* `internal/` is everything specific to the dharma application.
* `locales/` contains tools for managing translations, as well as the files
  containing translations for dharma.

## Locales & Translations

Messages that need to be localized are in `internal/render/i18n/messages.go`.
One can generate the files needed for translation by running the commands:

```
$ cd locales
$ ./extract.sh
```

This will generate a few files `locales/l10n/translate.*.toml`, which are ready
for a translator to modify. Once the files are ready to be used in the program,
they are published by:

```
$ cd locales
$ ./publish.sh
```

Which will turn the `translate.*.toml` files into `active.*.toml` files. These
files are then used in `go generate` commands.

## Assets

Assets are packed into the binary, for ease of deployment. If the `dev` build
tag is used, they are not packed so that local file changes take immediate
effect. Note for SCSS files, they will need to be re-processed manually by the
`sass` command or invoked via `go generate ./...`.
