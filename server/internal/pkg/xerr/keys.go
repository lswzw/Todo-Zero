package xerr

// localeKey 用于在 context 中存储语言信息的 key 类型
type localeKey struct{}

// LocaleKey context 中存储语言的 key
var LocaleKey interface{} = localeKey{}