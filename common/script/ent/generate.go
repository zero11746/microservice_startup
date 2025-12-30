package main

//go:generate go run entgo.io/ent/cmd/ent generate --feature intercept,privacy,schema/snapshot,sql/lock,sql/upsert,sql/modifier,sql/execquery ./internal/ent/schema
