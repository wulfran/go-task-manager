package db

import "embed"

//go:embed queries/task/*.sql queries/user/*.sql queries/utils/*.sql migrations/*.sql
var SQLFiles embed.FS
