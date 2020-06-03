package scenario

import "context"

// Verify Initialize後のアプリケーションサーバーに対して、副作用のない検証を実行する
// 早い段階でベンチマークをFailさせて早期リターンさせるのが目的
// ex) recommended API や Search API を叩いて初期状態を確認する
func Verify(ctx context.Context) {
	return
}
