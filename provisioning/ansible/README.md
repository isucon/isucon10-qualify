# isucon10-provisioning

ansible 2.9.13で動作確認しています

## playbooks
- benchmarker.yaml
  - ベンチマーカーがセットアップされます
- competitor.yaml
  - 競技者に提供された各種言語実装がセットアップされます
- allinone.yaml
  - 各種言語実装に加えてベンチマーカーのセットアップもされています

## Vagrantを利用して，環境をセットアップする

本Vagrantファイルは1台構成で，allinone.yamlを実行した結果を提供します
別の，設定を実行したい場合は

### 初回構築
下記コマンドで，VMの作り直しから始まります
- make vagrant/init

### ファイル初期化/再構築
ansibleの実行中に，通信環境エラーなどが起きた場合，下記のコマンドで再実行できます
- vagrant provision

