package postgres

func int64ScanDest(i *int64) []any {
	return []any{i}
}

func boolScanDest(b *bool) []any {
	return []any{b}
}

func stringScanDest(s *string) []any {
	return []any{s}
}
