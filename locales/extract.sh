#!/bin/sh

go get -u github.com/nicksnyder/go-i18n/v2/goi18n
goi18n extract -outdir=l10n ..
goi18n merge -outdir=l10n l10n/active.*.toml
