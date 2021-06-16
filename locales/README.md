# Locales

This is the directory responsible for user-facing translations. Translators are
welcome to directly contribute to any existing `translate.*.toml` files.

## Translators

Volunteers for translations are greatly appreciated. Please:

1) Modify the TRANSLATORS file with your first name, last name, and contact
information. You deserve credit!
2) Modify the `translate.<locale>.toml` file for the `locale` you are willing
and able to translate. All entries that are not the `description`, `hash`, or
`[id]` need translating.
3) Open a pull request with the modified translations. If you are unfamiliar
with `git` and do not know how, please reach out.

## Developers

### Requesting New Translations

Once the code has been updated with new or modified localized strings, they need
to be translated into all targeted languages.

To update internationalization into, first `cd` to this directory, then run
`extract.sh` to generate updated `translate.*.toml` files for supported locales.

The `translate.*.toml` files are then ready for translators.

### Publishing Translations

Translations can only be published in batch form for all locales. First `cd` to
this directory, then run `publish.sh`.

### Adding a new Locale

For a new `${LOCALE}`

```bash
touch translate.${LOCALE}.toml
goi18n merge active.en.toml translate.${LOCALE}.toml
cp translate.${LOCALE}.toml active.${LOCALE}.toml
```

Then, modify `internal/render/renderer.go` to include

```go
bundle.LoadMessageFile("active.${LOCALE}.toml")
```

Since we `cp` the untranslated file, it will begin in all-English, and translate
over time.
